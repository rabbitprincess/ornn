package config

import (
	"fmt"

	"ariga.io/atlas/sql/schema"
	"github.com/gokch/ornn/atlas"
	"github.com/gokch/ornn/db"
)

type Schema struct {
	*schema.Schema `json:"-"`
}

func (t *Schema) Init(dbType atlas.DbType, db *db.Conn) error {
	var err error
	atlas := &atlas.Atlas{}
	err = atlas.Init(dbType, db)
	if err != nil {
		return err
	}

	t.Schema, err = atlas.InspectSchema()
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
