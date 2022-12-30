package db_mysql

import (
	"strings"

	"github.com/gokch/ornn/db"
)

func NewVendor(db *db.Conn) *Vendor {
	return &Vendor{
		db: db,
	}
}

type Vendor struct {
	db *db.Conn
}

func (t *Vendor) ConvType(dbType string) (genType string) {
	var is_unsigned bool

	arrs_option := strings.Split(string(dbType), " ")
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

	return t.convType(s_field_type_with_len, is_unsigned)
}

func (t *Vendor) convType(_s_field_type_db string, _is_unsigned bool) string {
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

func (t *Vendor) CreateTable() (sql []string, err error) {
	// db 안 모든 테이블 이름을 가져옴
	tblNames, err := t.allTable()
	if err != nil {
		return nil, err
	}

	sql = make([]string, 0, len(tblNames))
	for _, tblName := range tblNames {
		s_sql := "show create table `" + tblName + "`"

		var tbl string
		var sqlCreateTable string

		job, err := t.db.Begin()
		if err != nil {
			return nil, err
		}

		rows, err := job.Query(s_sql)
		if err != nil {
			return nil, err
		}
		if rows.Next() == true {
			err = rows.Scan(&tbl, &sqlCreateTable)
			if err != nil {
				return nil, err
			}
		}

		sql = append(sql, sqlCreateTable)
	}
	return sql, nil
}

func (t *Vendor) allTable() (tables []string, err error) {
	s_sql := "select `table_name` from `information_schema`.`tables` where `table_schema` = '" + t.db.DbName + "'"

	job, err := t.db.Begin()
	if err != nil {
		return nil, err
	}
	rows, err := job.Query(s_sql)
	if err != nil {
		return nil, err
	}
	tables = make([]string, 0, 10)
	for rows.Next() {
		var tblName string = "table_name"
		err = rows.Scan(&tblName)
		if err != nil {
			return nil, err
		}
		tables = append(tables, tblName)
	}
	return tables, nil
}
