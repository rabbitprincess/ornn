package go_orm_gen

// BP 버젼
type TD_N4_bp_version int32

const (
	TD_N4_bp_version__now_build TD_N4_bp_version = 1
)

type LangType string

const (
	LangTypeNone LangType = "none"
	LangTypeGo   LangType = "go"
	LangTypeJava LangType = "java"
)

type DBType string

const (
	DBTypeNone  LangType = "none"
	DBTypeMysql LangType = "mysql"
)
