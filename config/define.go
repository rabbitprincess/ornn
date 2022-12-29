package config

// BP 버젼
const (
	version = 1
)

type DBType string

const (
	DBTypeNone  DBType = "none"
	DBTypeMysql DBType = "mysql"
)
