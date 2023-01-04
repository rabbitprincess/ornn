package config

import (
	"fmt"

	"ariga.io/atlas/sql/schema"
	"github.com/gokch/ornn/atlas"
)

type Schema struct {
	DbType atlas.DbType `json:"-"`

	*schema.Schema `json:"-"`
}

func (t *Schema) Init(dbType atlas.DbType, sch *schema.Schema) {
	t.DbType = dbType
	t.Schema = sch
}

func (t *Schema) GetTableFieldMatched(fieldName string, tablesName []string) (matched []string, err error) {
	matched = make([]string, 0, 10)

	for _, tableName := range tablesName {
		table, exist := t.Table(tableName)
		if exist != true {
			return nil, fmt.Errorf("wrong table name in sql query, table name is not exist in schema | table_name : %s", tableName)
		}
		_, exist = table.Column(fieldName)
		if exist == true {
			matched = append(matched, tableName)
		}
	}
	return matched, nil
}

func (t *Schema) GetFieldType(tableName, fieldName string) (fieldType string, err error) {
	if tableName == "" {
		for _, tbl := range t.Tables {
			fld, exist := tbl.Column(fieldName)
			if exist == true {
				return fld.Type.Raw, nil
			}
		}
	} else {
		tbl, exist := t.Schema.Table(tableName)
		if exist != true {
			return "", fmt.Errorf("not exist table | table name : %s", tableName)
		}
		fld, exist := tbl.Column(fieldName)
		if exist != true {
			return "", fmt.Errorf("not exist field | table name : %s | field name : %s", tableName, fieldName)
		}
		return fld.Type.Raw, nil
	}
	return "", fmt.Errorf("not exist field | field name : %s", fieldName)
}
