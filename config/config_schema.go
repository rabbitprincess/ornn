package config

import (
	"context"
	"fmt"

	"ariga.io/atlas/sql/mysql"
	"ariga.io/atlas/sql/schema"
	"github.com/gokch/ornn/db"
)

type Schema struct {
	*schema.Schema `json:"-"`
}

func (t *Schema) InitDefault(db *db.Conn) error {
	var err error

	migrate, err := mysql.Open(db.Db)
	if err != nil {
		return err
	}
	t.Schema, err = migrate.InspectSchema(context.Background(), "", nil)
	if err != nil {
		return err
	}
	return nil
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
