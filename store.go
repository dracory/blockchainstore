package blockchainstore

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"log/slog"
	"strings"
	"time"

	"github.com/dracory/neat"
	ormcontracts "github.com/dracory/neat/contracts/database/orm"
	contractsschema "github.com/dracory/neat/contracts/database/schema"
	"github.com/dracory/neat/database/schema/constants"
	"github.com/dracory/uid"
	"github.com/dromara/carbon/v2"
)

var _ StoreInterface = (*Store)(nil) // verify it extends the interface

type Store struct {
	blockTableName     string
	db                 *neat.Database
	timeoutSeconds     int64
	automigrateEnabled bool
	debugEnabled       bool
	sqlLogger          *slog.Logger
}

type blockRow struct {
	ID        string    `db:"id"`
	ParentID  string    `db:"parent_id"`
	Data      string    `db:"data"`
	Hash      string    `db:"hash"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	DeletedAt time.Time `db:"deleted_at"`
}

// AutoMigrate auto migrate
func (store *Store) AutoMigrate() error {
	return store.MigrateUp(context.Background())
}

// MigrateUp creates the block store tables
func (store *Store) MigrateUp(ctx context.Context, tx ...*sql.Tx) error {
	if !store.db.Schema().HasTable(store.blockTableName) {
		err := store.db.Schema().Create(store.blockTableName, func(table contractsschema.Blueprint) {
			table.String(COLUMN_ID, 40)
			table.Primary(COLUMN_ID)
			table.String(COLUMN_PARENT_ID, 40).Default("")
			table.LongText(COLUMN_DATA)
			table.String(COLUMN_HASH, 200).Default("")
			table.DateTime(COLUMN_CREATED_AT).GetUseCurrent()
			table.DateTime(COLUMN_UPDATED_AT).GetUseCurrent()
			table.DateTime(COLUMN_DELETED_AT).Default(constants.MaxSoftDeletedAtDefault)
		})
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

// MigrateDown drops the block store tables
func (store *Store) MigrateDown(ctx context.Context, tx ...*sql.Tx) error {
	if store.db.Schema().HasTable(store.blockTableName) {
		err := store.db.Schema().Drop(store.blockTableName)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

// EnableDebug - enables the debug option
func (st *Store) EnableDebug(debug bool) StoreInterface {
	st.debugEnabled = debug
	return st
}

func (store *Store) BlockCreate(ctx context.Context, block BlockInterface) error {
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if block == nil {
		return errors.New("block is nil")
	}

	if block.GetID() == "" {
		block.SetID(uid.HumanUid())
	}
	if block.Timestamp() == "" {
		block.SetTimestamp(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}

	deletedAt := block.Get(COLUMN_DELETED_AT)
	var deletedAtTime time.Time
	if deletedAt != "" && deletedAt != constants.MaxSoftDeletedAtDefault {
		deletedAtTime = carbon.Parse(deletedAt, carbon.UTC).StdTime()
	} else {
		deletedAtTime = carbon.Parse(constants.MaxSoftDeletedAtDefault, carbon.UTC).StdTime()
	}

	createdAtTime := carbon.Parse(block.Timestamp(), carbon.UTC).StdTime()

	row := blockRow{
		ID:        block.GetID(),
		ParentID:  block.PreviousHash(),
		Data:      block.Data(),
		Hash:      block.ThisHash(),
		CreatedAt: createdAtTime,
		UpdatedAt: createdAtTime,
		DeletedAt: deletedAtTime,
	}

	err := store.db.Query().Table(store.blockTableName).Create(&row)
	if err != nil {
		return err
	}

	block.MarkAsNotDirty()

	return nil
}

func (store *Store) BlockDelete(ctx context.Context, block BlockInterface) error {
	if block == nil {
		return errors.New("block is nil")
	}

	return store.BlockDeleteByID(ctx, block.GetID())
}

func (store *Store) BlockDeleteByID(ctx context.Context, id string) error {
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if id == "" {
		return errors.New("block id is empty")
	}

	_, err := store.db.Query().
		Table(store.blockTableName).
		Where(COLUMN_ID+" = ?", id).
		Delete()

	return err
}

func (store *Store) BlockFindByID(ctx context.Context, id string) (BlockInterface, error) {
	if id == "" {
		return nil, errors.New("block id is empty")
	}

	list, err := store.BlockList(ctx, BlockQueryOptions{
		ID:    id,
		Limit: 1,
	})

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

func (store *Store) BlockList(ctx context.Context, options BlockQueryOptions) ([]BlockInterface, error) {
	if ctx == nil {
		return nil, errors.New("ctx is nil")
	}

	q := store.blockQuery(options)

	var rows []blockRow
	if err := q.Table(store.blockTableName).Get(&rows); err != nil {
		return []BlockInterface{}, err
	}

	list := make([]BlockInterface, 0, len(rows))

	for _, r := range rows {
		model := mapRowToBlock(r)
		list = append(list, model)
	}

	return list, nil
}

func (store *Store) BlockUpdate(ctx context.Context, block BlockInterface) error {
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if block == nil {
		return errors.New("block is nil")
	}

	rawChanged := block.GetDataChanged()
	dataChanged := map[string]interface{}{}
	if _, ok := rawChanged[COLUMN_PREVIOUS_HASH]; ok {
		dataChanged[COLUMN_PARENT_ID] = block.PreviousHash()
	}
	if _, ok := rawChanged[COLUMN_THIS_HASH]; ok {
		dataChanged[COLUMN_HASH] = block.ThisHash()
	}
	if _, ok := rawChanged[COLUMN_DATA]; ok {
		dataChanged[COLUMN_DATA] = block.Data()
	}
	if _, ok := rawChanged[COLUMN_DELETED_AT]; ok {
		deletedAt := block.Get(COLUMN_DELETED_AT)
		if deletedAt != "" && deletedAt != constants.MaxSoftDeletedAtDefault {
			dataChanged[COLUMN_DELETED_AT] = carbon.Parse(deletedAt, carbon.UTC).StdTime()
		} else {
			dataChanged[COLUMN_DELETED_AT] = carbon.Parse(constants.MaxSoftDeletedAtDefault, carbon.UTC).StdTime()
		}
	}

	if len(dataChanged) < 1 {
		return nil
	}

	dataChanged[COLUMN_UPDATED_AT] = carbon.Now(carbon.UTC).StdTime()

	_, err := store.db.Query().
		Table(store.blockTableName).
		Where(COLUMN_ID+" = ?", block.GetID()).
		Update(dataChanged)

	if err != nil {
		return err
	}

	block.MarkAsNotDirty()

	return nil
}

func (store *Store) blockQuery(options BlockQueryOptions) ormcontracts.Query {
	q := store.db.Query().Table(store.blockTableName)

	if options.ID != "" {
		q = q.Where(COLUMN_ID+" = ?", options.ID)
	}

	if len(options.IDIn) > 0 {
		inClause := COLUMN_ID + " IN ("
		placeholders := make([]interface{}, 0, len(options.IDIn))
		for i, id := range options.IDIn {
			if i > 0 {
				inClause += ", "
			}
			inClause += "?"
			placeholders = append(placeholders, id)
		}
		inClause += ")"
		q = q.Where(inClause, placeholders...)
	}

	if !options.CountOnly {
		if options.Limit > 0 {
			q = q.Limit(options.Limit)
		}

		if options.Offset > 0 {
			q = q.Offset(options.Offset)
		}
	}

	if options.OrderBy != "" {
		sortOrder := "desc"
		if options.SortOrder != "" {
			sortOrder = strings.ToLower(options.SortOrder)
		}
		q = q.OrderBy(options.OrderBy, sortOrder)
	}

	if options.WithDeleted {
		q = q.WithSoftDeleted()
	} else {
		q = q.Where(COLUMN_DELETED_AT+" > ?", carbon.Now(carbon.UTC).StdTime())
	}

	return q
}

type BlockQueryOptions struct {
	ID          string
	IDIn        []string
	Offset      int
	Limit       int
	SortOrder   string
	OrderBy     string
	CountOnly   bool
	WithDeleted bool
}

// logSql logs sql to the sql logger, if debug mode is enabled
func (store *Store) logSql(sqlOperationType string, sql string, params ...interface{}) {
	if !store.debugEnabled {
		return
	}

	if store.sqlLogger != nil {
		store.sqlLogger.Debug("sql: "+sqlOperationType, slog.String("sql", sql), slog.Any("params", params))
	}
}

func mapRowToBlock(r blockRow) BlockInterface {
	block := NewBlock()
	block.SetID(r.ID)
	block.SetPreviousHash(r.ParentID)
	block.SetData(r.Data)
	block.SetThisHash(r.Hash)
	block.SetTimestamp(carbon.CreateFromStdTime(r.CreatedAt).ToDateTimeString(carbon.UTC))
	if !r.DeletedAt.IsZero() && carbon.CreateFromStdTime(r.DeletedAt).ToDateTimeString(carbon.UTC) != constants.MaxSoftDeletedAtDefault {
		block.Set(COLUMN_DELETED_AT, carbon.CreateFromStdTime(r.DeletedAt).ToDateTimeString(carbon.UTC))
	}
	block.MarkAsNotDirty()
	return block
}
