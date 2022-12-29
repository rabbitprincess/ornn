package db_mysql

import "fmt"

func NewDsn(_s_id, _s_pw, _s_addr, _s_port, _s_db_name string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true", _s_id, _s_pw, _s_addr, _s_port, _s_db_name)
}
