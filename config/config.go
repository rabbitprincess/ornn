package config

import (
	"encoding/json"
	"io/ioutil"
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
