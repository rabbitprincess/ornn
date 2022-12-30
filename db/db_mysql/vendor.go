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
	var unsigned bool

	opts := strings.Split(string(dbType), " ")
	for _, opt := range opts {
		opt = strings.ToLower(opt)
		if opt == "unsigned" {
			unsigned = true
		}
	}

	if len(opts) == 0 {
		return ""
	}

	s_field_type_with_len := opts[0]

	n_pos_start_len := strings.Index(s_field_type_with_len, "(")
	if n_pos_start_len != -1 {
		s_field_type_with_len = s_field_type_with_len[0:n_pos_start_len]
	}

	return t.convType(s_field_type_with_len, unsigned)
}

func (t *Vendor) convType(dbType string, unsigned bool) string {
	switch strings.ToLower(dbType) {
	case "char", "varchar", "tinytext", "text", "mediumtext", "longtext", "json":
		return "string"
	case "binary", "varbinary", "tinyblob", "blob", "mediumblob", "longblob":
		return "[]byte"
	case "tinyint":
		if unsigned == true {
			return "uint8"
		}
		return "int8"
	case "smallint":
		if unsigned == true {
			return "uint16"
		}
		return "int16"
	case "int":
		if unsigned == true {
			return "uint32"
		}
		return "int32"
	case "bigint":
		if unsigned == true {
			return "uint64"
		}
		return "int64"
	case "float":
		return "float32"
	case "double", "real":
		return "float64"
	default:
		return "error"
	}
}

func (t *Vendor) CreateTable() (sql []string, err error) {
	// db 안 모든 테이블 이름을 가져옴
	tableNames, err := t.allTable()
	if err != nil {
		return nil, err
	}

	sql = make([]string, 0, len(tableNames))
	for _, tableName := range tableNames {
		s_sql := "show create table `" + tableName + "`"

		var table string
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
			err = rows.Scan(&table, &sqlCreateTable)
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
		var tableName string = "table_name"
		err = rows.Scan(&tableName)
		if err != nil {
			return nil, err
		}
		tables = append(tables, tableName)
	}
	return tables, nil
}
