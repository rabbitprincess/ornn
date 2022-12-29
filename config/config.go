package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Schema  Schema  `json:"db"`
	GenCode GenCode `json:"gen"`
}

func (t *Config) Load(config string) error {
	bt_data, err := ioutil.ReadFile(config)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bt_data, &t)
	if err != nil {
		return err
	}
	return nil
}

func (t *Config) Save(config string) error {
	bt_data, err := json.MarshalIndent(&t, "", "\t")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(config, bt_data, 0700)
	if err != nil {
		return err
	}
	return nil
}
