package ornn

import (
	"fmt"
	"strings"

	"github.com/blastrain/vitess-sqlparser/sqlparser"
	"github.com/gokch/ornn/config"
	"github.com/gokch/ornn/db"
)

func NewVendor(vendor db.Vendor) *Vendor {
	return &Vendor{dbVendor: vendor}
}

// vendor config by schema
type Vendor struct {
	dbVendor db.Vendor
}

func (t *Vendor) VendorBySchema(conf *config.Config) error {
	// sql(create table) 추출
	sqlCreateTables, err := t.dbVendor.CreateTable()
	if err != nil {
		return err
	}

	// Schema 및 Default query 초기화
	conf.Schema = config.Schema{}
	conf.Queries.Default = make(map[string][]*config.Query)

	// 뽑아낸 쿼리를 이용해서 table 을 제작
	for _, sqlCreateTable := range sqlCreateTables {
		table, queries, err := t.vendorByTable(sqlCreateTable)
		if err != nil {
			return err
		}
		conf.Schema.AddTable(table)

		for _, query := range queries {
			conf.Queries.AddQueryTables(table.Name, query)
		}
	}

	return nil
}

func (t *Vendor) vendorByTable(sqlCreateTable string) (*config.Table, []*config.Query, error) {
	stmt, err := sqlparser.Parse(sqlCreateTable)
	if err != nil {
		panic(err)
	}
	parser := stmt.(*sqlparser.CreateTable)
	tableName := parser.NewName.Name.String()

	// parse table
	table := &config.Table{}
	table.Init(tableName)

	// parse index
	for _, idx := range parser.Constraints {
		index := &config.Index{}
		keys := make([]string, 0, len(idx.Keys))
		for _, key := range idx.Keys {
			keys = append(keys, key.String())
		}
		index.Set(idx.Name, idx.Type.String(), keys)
		table.AddIndex(index)
	}

	// parse field
	for _, column := range parser.Columns {
		field := &config.Field{}
		field.Set(column.Name, column.Type, t.dbVendor.ConvType(column.Type))
		table.AddField(field)
	}

	// make query - TODO : Query Generator 공통화
	queries := make([]*config.Query, 0, 100)

	// insert all
	questionare := strings.Repeat("?, ", len(table.Fields))
	questionare = questionare[:len(questionare)-2]
	queries = append(queries, &config.Query{
		Name:    "insert",
		Comment: "default query - insert all",
		Sql:     fmt.Sprintf("INSERT INTO %s VALUES (%s)", tableName, questionare),
	})

	// select all
	queries = append(queries, &config.Query{
		Name:    "select",
		Comment: "default query - select all",
		Sql:     fmt.Sprintf("SELECT * FROM %s", tableName),
	})

	// TODO: select where by index

	// TODO: update

	// delete
	queries = append(queries, &config.Query{
		Name:    "delete",
		Comment: "default query - delete all",
		Sql:     fmt.Sprintf("DELETE FROM %s", tableName),
	})

	return table, queries, nil
}
