package bp

import (
	"fmt"
	"io/ioutil"
	"module/db"
)

type I_BP__gen interface {
	gen_source_code(
		_pt_bp_config *T_BP__config,
		_maps_config map[string]string,
		_pt_gen_data *T_BP__gen__data,
	) (
		s_code string,
		err error,
	)
}

type T_BP__gen struct {
}

func (t *T_BP__gen) Get_source_code(
	_pc_db *db.C_DB_conn,
	_td_n1_db_rds_type TD_N1_db_rds_type,
	_pt_bp_config *T_BP__config,
	_td_lang_name TD_S_lang_name,
	_maps_gen_config map[string]string,
	_s_filepath string,
) (
	s_code string,
	err error,
) {

	// json -> 코드 생성을 위한 gen 데이터 준비
	pt_gen_data := &T_BP__gen__data{}
	{
		pt_gen_data.Init(_pc_db, _td_n1_db_rds_type)

		err = pt_gen_data.prepare_db_schema(_pt_bp_config)
		if err != nil {
			return "", err
		}

		if _pt_bp_config.T_db.T_caller.is_exist_error__caller_sql == true {
			_pt_bp_config.T_db.T_caller.Print__sql_errors()
		}

	}

	// 언어별 코드 생성
	{
		switch _td_lang_name {
		case TD_S_lang_name__go:
			pt := &T_BP__gen__go{}
			s_code, err = pt.gen_source_code(_pt_bp_config, _maps_gen_config, pt_gen_data)
		case TD_S_lang_name__java:
			pt := &T_BP__gen__java{}
			s_code, err = pt.gen_source_code(_pt_bp_config, _maps_gen_config, pt_gen_data)
		default:
			return "", fmt.Errorf("json is emtpy")
		}
		if err != nil {
			return "", err
		}
	}
	// write file
	{
		err = ioutil.WriteFile(_s_filepath, []byte(s_code), 0700)
		if err != nil {
			return "", err
		}
	}
	return s_code, nil
}
