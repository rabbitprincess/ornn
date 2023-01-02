package atlas

type DbType int8

const (
	DbTypeEmpty DbType = iota
	DbTypeMySQL
	DbTypePostgre
	DbTypeSQLite

	// TODO
	// DbTypeMaria
	// DbTypeTiDB
	// DbTypeCockroachDB
)
