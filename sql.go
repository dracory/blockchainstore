package blockchainstore

import "github.com/dracory/sb"

// SQLCreateTable returns a SQL string for creating the cache table
func (st *Store) sqlCreateTable() (string, error) {
	sql, err := sb.NewBuilder(st.dbDriverName).
		Table(st.blockTableName).
		Column(sb.Column{
			Name:       "id",
			Type:       sb.COLUMN_TYPE_STRING,
			Length:     40,
			PrimaryKey: true,
		}).
		Column(sb.Column{
			Name:   "parent_id",
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 40,
		}).
		Column(sb.Column{
			Name: "data",
			Type: sb.COLUMN_TYPE_LONGTEXT,
		}).
		Column(sb.Column{
			Name:   "hash",
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 200,
		}).
		Column(sb.Column{
			Name: "created_at",
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name: "updated_at",
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name: "deleted_at",
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		CreateIfNotExists()

	return sql, err
}
