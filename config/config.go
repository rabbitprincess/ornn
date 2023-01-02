package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"ariga.io/atlas/sql/schema"
)

type Config struct {
	Global  Global  `json:"global"`
	Schema  Schema  `json:"schema"`
	Queries Queries `json:"queries"`
}

// TODO - 추후 config 형식 변경 예정
func (t *Config) Load(config string) error {
	data, err := ioutil.ReadFile(config)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &t)
	if err != nil {
		return err
	}

	// queries 초기화
	t.Queries.init(&t.Schema)

	return nil
}

func (t *Config) InitSchema(schema *schema.Schema) {
	t.Schema.Init(schema)
}

func (t *Config) Save(config string) error {
	data, err := json.MarshalIndent(&t, "", "\t")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(config, data, 0700)
	if err != nil {
		return err
	}
	return nil
}

func (t *Config) VendorBySchema() error {
	t.Queries.Default = make(map[string][]*Query)

	// 뽑아낸 쿼리를 이용해서 table 을 제작
	for _, table := range t.Schema.Tables {
		queries, err := t.vendorByTable(table)
		if err != nil {
			return err
		}

		for _, query := range queries {
			t.Queries.AddQueryTables(table.Name, query)
		}
	}

	return nil
}

func (t *Config) vendorByTable(table *schema.Table) ([]*Query, error) {
	// make query - TODO : Query Generator 공통화
	queries := make([]*Query, 0, 100)

	// insert all
	questionare := strings.Repeat("?, ", len(table.Columns))
	questionare = questionare[:len(questionare)-2]
	queries = append(queries, &Query{
		Name:    "insert",
		Comment: "default query - insert all",
		Sql:     fmt.Sprintf("INSERT INTO %s VALUES (%s)", table.Name, questionare),
	})

	// select all
	queries = append(queries, &Query{
		Name:    "select",
		Comment: "default query - select all",
		Sql:     fmt.Sprintf("SELECT * FROM %s", table.Name),
	})

	// TODO: select where by index

	// TODO: update

	// delete
	queries = append(queries, &Query{
		Name:    "delete",
		Comment: "default query - delete all",
		Sql:     fmt.Sprintf("DELETE FROM %s", table.Name),
	})

	return queries, nil
}
