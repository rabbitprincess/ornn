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
	tblName := parser.NewName.Name.String()

	table := &config.Table{}
	table.Init(tblName)

	for _, idx := range parser.Constraints {
		index := &config.Index{}
		arrs_keys := make([]string, 0, len(idx.Keys))
		for _, pt_key := range idx.Keys {
			arrs_keys = append(arrs_keys, pt_key.String())
		}
		index.Set(idx.Name, idx.Type.String(), arrs_keys)
		table.AddIndex(index)
	}

	for _, fld := range parser.Columns {
		field := &config.Field{}
		field.Set(fld.Name, fld.Type, t.vendor.ConvType(fld.Type))
		table.AddField(field)
	}
	return table, nil
}
