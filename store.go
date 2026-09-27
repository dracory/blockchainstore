package blockchainstore

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"log/slog"
	"strings"

	"github.com/dracory/database"
	"github.com/dracory/neat"
	ormcontracts "github.com/dracory/neat/contracts/database/orm"
	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
	"github.com/samber/lo"
)

const BLOCK_TABLE_NAME = "blocks_block"

var _ StoreInterface = (*Store)(nil) // verify it extends the interface

type Store struct {
	blockTableName     string
	db                 *sql.DB
	neatDB             *neat.Database
	dbDriverName       string
	timeoutSeconds     int64
	automigrateEnabled bool
	debugEnabled       bool
	sqlLogger          *slog.Logger
}

// AutoMigrate auto migrate
func (store *Store) AutoMigrate() error {
	sql, err := store.sqlCreateTable()
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = store.db.Exec(sql)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// EnableDebug - enables the debug option
func (st *Store) EnableDebug(debug bool) {
	st.debugEnabled = debug
}

func (store *Store) BlockCreate(ctx context.Context, block *Block) error {
	block.SetTimestamp(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	data := blockToMap(block)

	query := store.neatDB.Query().Table(store.blockTableName)
	sqlStr := query.ToRawSql().Create(data)

	store.logSql("insert", sqlStr)

	_, err := database.Execute(store.toQuerableContext(ctx), sqlStr)
	if err != nil {
		return err
	}

	block.MarkAsNotDirty()

	return nil
}

func (store *Store) BlockDelete(ctx context.Context, block *Block) error {
	if block == nil {
		return errors.New("block is nil")
	}

	return store.BlockDeleteByID(ctx, block.ID())
}

func (store *Store) BlockDeleteByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("block id is empty")
	}

	query := store.neatDB.Query().Table(store.blockTableName).Where("id = ?", id)
	sqlStr := query.ToRawSql().Delete()

	store.logSql("delete", sqlStr)

	_, err := database.Execute(store.toQuerableContext(ctx), sqlStr)

	return err
}

func (store *Store) BlockFindByID(ctx context.Context, id string) (*Block, error) {
	if id == "" {
		return nil, errors.New("exam id is empty")
	}

	list, err := store.BlockList(ctx, BlockQueryOptions{
		ID:    id,
		Limit: 1,
	})

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return &list[0], nil
	}

	return nil, nil
}

func (store *Store) BlockList(ctx context.Context, options BlockQueryOptions) ([]Block, error) {
	q := store.blockQuery(options)

	sqlStr := q.ToRawSql().Get(nil)

	store.logSql("select", sqlStr)

	modelMaps, err := database.SelectToMapString(store.toQuerableContext(ctx), sqlStr)
	if err != nil {
		return []Block{}, err
	}

	list := []Block{}

	lo.ForEach(modelMaps, func(modelMap map[string]string, index int) {
		model := mapToBlock(modelMap)
		list = append(list, *model)
	})

	return list, nil
}

func (store *Store) BlockUpdate(ctx context.Context, block *Block) error {
	if block == nil {
		return errors.New("order is nil")
	}

	rawChanged := block.DataChanged()
	dataChanged := map[string]interface{}{}
	if _, ok := rawChanged["previous_hash"]; ok {
		dataChanged["parent_id"] = block.PreviousHash()
	}
	if _, ok := rawChanged["deleted_at"]; ok {
		dataChanged["deleted_at"] = block.Get("deleted_at")
	}

	if len(dataChanged) < 1 {
		return nil
	}

	query := store.neatDB.Query().Table(store.blockTableName).Where("id = ?", block.ID())
	sqlStr := query.ToRawSql().Update(dataChanged)

	store.logSql("update", sqlStr)

	_, err := database.Execute(store.toQuerableContext(ctx), sqlStr)

	block.MarkAsNotDirty()

	return err
}

func (store *Store) blockQuery(options BlockQueryOptions) ormcontracts.Query {
	q := store.neatDB.Query().Table(store.blockTableName)

	if options.ID != "" {
		q = q.Where("id = ?", options.ID)
	}

	if !options.CountOnly {
		if options.Limit > 0 {
			q = q.Limit(options.Limit)
		}

		if options.Offset > 0 {
			q = q.Offset(options.Offset)
		}
	}

	sortOrder := "desc"
	if options.SortOrder != "" {
		sortOrder = options.SortOrder
	}

	if options.OrderBy != "" {
		if strings.EqualFold(sortOrder, sb.ASC) {
			q = q.OrderBy(options.OrderBy, "asc")
		} else {
			q = q.OrderBy(options.OrderBy, "desc")
		}
	}

	if !options.WithDeleted {
		q = q.Where("deleted_at = ?", sb.NULL_DATETIME)
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

func (store *Store) toQuerableContext(context context.Context) database.QueryableContext {
	if database.IsQueryableContext(context) {
		return context.(database.QueryableContext)
	}

	return database.Context(context, store.db)
}

func blockToMap(block *Block) map[string]interface{} {
	m := map[string]interface{}{}
	m["id"] = block.ID()
	m["parent_id"] = block.PreviousHash()
	m["data"] = block.Data()
	m["hash"] = block.ThisHash()
	m["created_at"] = block.Timestamp()
	m["updated_at"] = block.Timestamp()
	if block.Get("deleted_at") != "" {
		m["deleted_at"] = block.Get("deleted_at")
	} else {
		m["deleted_at"] = sb.NULL_DATETIME
	}
	return m
}

func mapToBlock(m map[string]string) *Block {
	block := NewBlock()
	block.SetID(m["id"])
	block.SetPreviousHash(m["parent_id"])
	block.SetData(m["data"])
	block.SetThisHash(m["hash"])
	block.SetTimestamp(m["created_at"])
	if deletedAt, ok := m["deleted_at"]; ok {
		block.Set("deleted_at", deletedAt)
	}
	block.MarkAsNotDirty()
	return block
}
