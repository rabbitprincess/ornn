package go_orm_gen

import (
	"github.com/blastrain/vitess-sqlparser/sqlparser"
	"github.com/gokch/go-orm-gen/config"
	"github.com/gokch/go-orm-gen/db"
)

type DbVendor struct {
	vendor db.Vendor
}

//------------------------------------------------------------------------------------------------//
// Schema

func (t *DbVendor) Init(vendor db.Vendor) {
	t.vendor = vendor
}

func (t *DbVendor) SchemaGet() (*config.Schema, error) {
	// sql(create table) 추출
	sqlCreateTables, err := t.vendor.CreateTable()
	if err != nil {
		return nil, err
	}

	schema := &config.Schema{}
	schema.Init()

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

func (t *DbVendor) schemaGet(_s_sql__create_table string) (*config.Table, error) {
	stmt, err := sqlparser.Parse(_s_sql__create_table)
	if err != nil {
		panic(err)
	}
	parser := stmt.(*sqlparser.CreateTable)
	tblName := parser.NewName.Name.String()

	table := &config.Table{}
	table.Init(tblName)

	for _, pt_index__parser := range parser.Constraints {
		index := &config.Index{}
		arrs_keys := make([]string, 0, len(pt_index__parser.Keys))
		for _, pt_key := range pt_index__parser.Keys {
			arrs_keys = append(arrs_keys, pt_key.String())
		}
		index.Set(pt_index__parser.Name, pt_index__parser.Type.String(), arrs_keys)
		table.AddIndex(index)
	}

	for _, pt_field__parser := range parser.Columns {
		field := &config.Field{}
		field.Set(pt_field__parser.Name, pt_field__parser.Type, t.vendor.ConvType(pt_field__parser.Type))
		table.AddField(field)
	}
	return table, nil
}
