package ornn

import (
	"github.com/blastrain/vitess-sqlparser/sqlparser"
	"github.com/gokch/ornn/config"
	"github.com/gokch/ornn/db"
)

func NewVendor(vendor db.Vendor) *DbVendor {
	return &DbVendor{vendor: vendor}
}

type DbVendor struct {
	vendor db.Vendor
}

func (t *DbVendor) SchemaGet() (*config.Schema, error) {
	// sql(create table) 추출
	sqlCreateTables, err := t.vendor.CreateTable()
	if err != nil {
		return nil, err
	}

	schema := &config.Schema{}

	// 뽑아낸 쿼리를 이용해서 table 을 제작
	for _, sqlCreateTable := range sqlCreateTables {
		table, err := t.schemaGet(sqlCreateTable)
		if err != nil {
			return nil, err
		}
		schema.AddTable(table)
	}
	return schema, nil
}

func (t *DbVendor) schemaGet(sqlCreateTable string) (*config.Table, error) {
	stmt, err := sqlparser.Parse(sqlCreateTable)
	if err != nil {
		panic(err)
	}
	parser := stmt.(*sqlparser.CreateTable)
	tableName := parser.NewName.Name.String()

	table := &config.Table{}
	table.Init(tableName)

	for _, idx := range parser.Constraints {
		index := &config.Index{}
		keys := make([]string, 0, len(idx.Keys))
		for _, key := range idx.Keys {
			keys = append(keys, key.String())
		}
		index.Set(idx.Name, idx.Type.String(), keys)
		table.AddIndex(index)
	}

	for _, column := range parser.Columns {
		field := &config.Field{}
		field.Set(column.Name, column.Type, t.vendor.ConvType(column.Type))
		table.AddField(field)
	}
	return table, nil
}
