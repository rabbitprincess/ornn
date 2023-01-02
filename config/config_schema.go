package config

import (
	"fmt"

	"ariga.io/atlas/sql/schema"
	"github.com/gokch/ornn/atlas"
	"github.com/gokch/ornn/db/db_mysql"
	"github.com/gokch/ornn/db/db_postgres"
	"github.com/gokch/ornn/db/db_sqlite"
)

type Schema struct {
	dbType   atlas.DbType        `json:"-"`
	ConvType func(string) string `json:"-"`

	*schema.Schema `json:"-"`
}

func (t *Schema) Init(dbType atlas.DbType, sch *schema.Schema) {
	t.dbType = dbType
	t.Schema = sch

	switch dbType {
	case atlas.DbTypeMySQL, atlas.DbTypeMaria, atlas.DbTypeTiDB:
		t.ConvType = db_mysql.ConvType
	case atlas.DbTypePostgre, atlas.DbTypeCockroachDB:
		t.ConvType = db_postgres.ConvType
	case atlas.DbTypeSQLite:
		t.ConvType = db_sqlite.ConvType
	}

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
