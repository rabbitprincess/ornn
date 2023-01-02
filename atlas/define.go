package atlas

type DbType int8

const (
	DbTypeEmpty DbType = iota
	DbTypeMySQL
	DbTypeMaria
	DbTypePostgre
	DbTypeSQLite
	DbTypeTiDB
	DbTypeCockroachDB
)
