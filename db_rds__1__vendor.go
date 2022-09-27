package bp

import "module/db"

type I_DB_RDS__vendor interface {
	conv_field_type__to_bp(string) string // DBMS  에 있는 field type string 을 받아 BP 용 data type 으로 변환
	get_sql__create_table(_pc_db *db.C_DB_conn) (arrs_sql__create_table []string, err error)
}
