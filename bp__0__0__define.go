package bp

// BP 버젼
type TD_N4_bp_version int32

const (
	TD_N4_bp_version__now_build TD_N4_bp_version = 1
)

type TD_S_lang_name string

const (
	TD_S_lang_name__error TD_S_lang_name = "error"
	TD_S_lang_name__go    TD_S_lang_name = "go"
	TD_S_lang_name__java  TD_S_lang_name = "java"
)

type TD_S_dbms_name string

const (
	TD_S_dbms_name__error TD_S_lang_name = "error"
	TD_S_dbms_name__mysql TD_S_lang_name = "mysql"
)
