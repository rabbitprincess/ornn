package config

import (
	"encoding/json"
	"io/ioutil"

	"ariga.io/atlas/sql/schema"
)

type Config struct {
	Global  Global  `json:"global"`
	Schema  Schema  `json:"-"`
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

	return nil
}

func (t *Config) InitSchema(schema *schema.Schema) error {
	// set schema
	t.Schema.Init(schema)

	// set queries by schema
	t.Queries.init(&t.Schema)
	t.Queries.InitQueryTables()

	return nil
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
