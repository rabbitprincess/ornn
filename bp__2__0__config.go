package bp

import (
	"encoding/json"
	"io/ioutil"
)

type T_BP__config struct {
	T_base_info T_BP__config__base_info `json:"base_info"`
	T_gen       T_BP__config__gen       `json:"gen"`
	T_db        T_BP__config__db        `json:"db"`
}

func (t *T_BP__config) File__load(_s_path__json_config_file string) error {
	var err error
	bt_data, err := ioutil.ReadFile(_s_path__json_config_file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bt_data, &t)
	if err != nil {
		return err
	}
	return nil
}

func (t *T_BP__config) File__save(_s_path__json_config_file string) error {
	var err error

	// 현재 버젼 저장
	t.T_base_info.Set_version(TD_N4_bp_version__now_build)

	bt_data, err := json.MarshalIndent(&t, "", "\t")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(_s_path__json_config_file, bt_data, 0700)
	if err != nil {
		return err
	}
	return nil
}
