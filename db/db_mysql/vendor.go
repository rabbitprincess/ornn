package db_mysql

import (
	"strings"

	"github.com/gokch/go-orm-gen/db"
)

// ---------------------------------------------------------------------------------//
// mysql
type T_DB_RDS__vendor__mysql struct {
}

func (t *T_DB_RDS__vendor__mysql) conv_field_type__to_bp(_s_field_type__db string) (_s_field_type__bp string) {
	var is_unsigned bool

	arrs_option := strings.Split(string(_s_field_type__db), " ")
	for _, s_option := range arrs_option {
		s_option = strings.ToLower(s_option)
		if s_option == "unsigned" {
			is_unsigned = true

		}
	}

	if len(arrs_option) == 0 {
		return ""
	}

	s_field_type_with_len := arrs_option[0]

	n_pos_start_len := strings.Index(s_field_type_with_len, "(")
	if n_pos_start_len != -1 {
		s_field_type_with_len = s_field_type_with_len[0:n_pos_start_len]
	}

	_s_field_type__bp = t.conv_field_type__sub(s_field_type_with_len, is_unsigned)
	return _s_field_type__bp
}

func (t *T_DB_RDS__vendor__mysql) conv_field_type__sub(_s_field_type_db string, _is_unsigned bool) string {
	switch strings.ToLower(_s_field_type_db) {
	case "char", "varchar", "tinytext", "text", "mediumtext", "longtext", "json":
		return "s"
	case "binary", "varbinary", "tinyblob", "blob", "mediumblob", "longblob":
		return "bt"
	case "tinyint":
		if _is_unsigned == true {
			return "u1"
		}
		return "n1"
	case "smallint":
		if _is_unsigned == true {
			return "u2"
		}
		return "n2"
	case "int":
		if _is_unsigned == true {
			return "u4"
		}
		return "n4"
	case "bigint":
		if _is_unsigned == true {
			return "u8"
		}
		return "n8"
	case "float":
		return "f"
	case "double", "real":
		return "d"
	default:
		return "error"
	}
}

// 임시 - 작업 필요
// func (t *T_DB_RDS__vendor__mysql) sql(_pc_db *db.C_DB_conn) (arrs_sql__create_table []string, err error) {
// }

//--------------------------------------------------------------------------------------------------------------------------------------//

func (t *T_DB_RDS__vendor__mysql) get_sql__create_table(_pc_db *db.DB) (arrs_sql__create_table []string, err error) {
	// db 안 모든 테이블 이름을 가져옴
	arrs_table_name, err := t.get_table_name__all(_pc_db)
	if err != nil {
		return nil, err
	}

	arrs_sql__create_table = make([]string, 0, len(arrs_table_name))
	for _, s_table_name := range arrs_table_name {
		s_sql := "show create table `" + s_table_name + "`"

		var s_table_name string
		var s_sql__create_table string

		is_end, err := _pc_db.Query(s_sql).Row_next(&s_table_name, &s_sql__create_table)
		if err != nil {
			return nil, err
		}
		if is_end == true {
			return nil, nil
		}

		arrs_sql__create_table = append(arrs_sql__create_table, s_sql__create_table)
	}
	return arrs_sql__create_table, nil
}

func (t *T_DB_RDS__vendor__mysql) get_table_name__all(_pc_db *db.DB) (arrs_table_name []string, err error) {
	s_sql := "select `table_name` from `information_schema`.`tables` where `table_schema` = '" + _pc_db.DB__name_get() + "'"

	_, arrpt_row_map, err := _pc_db.Query(s_sql).Row_all__map()
	if err != nil {
		return nil, err
	}
	arrs_table_name = make([]string, 0, len(arrpt_row_map))
	for _, pt_row_map := range arrpt_row_map {
		s_table_name := pt_row_map["table_name"].(*db.T_Col__bt).String()
		arrs_table_name = append(arrs_table_name, s_table_name)
	}
	return arrs_table_name, nil
}
