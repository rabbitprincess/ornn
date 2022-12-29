package go_orm_gen

import (
	"io/ioutil"
	"module/db"
)

type Gen struct {
	data *GenData
	code *GenCode
}

func (t *Gen) Gen(
	_pc_db *db.C_DB_conn,
	_td_n1_db_rds_type TD_N1_db_rds_type,
	_pt_bp_config *T_BP__config,
	_td_lang_name LangType,
	_maps_gen_config map[string]string,
	_s_filepath string,
) (code string, err error) {
	// json -> 코드 생성을 위한 gen 데이터 준비
	t.data = &GenData{}
	t.data.Init(_pc_db, _td_n1_db_rds_type)
	err = t.data.prepare_db_schema(_pt_bp_config)
	if err != nil {
		return "", err
	}

	if _pt_bp_config.T_db.T_caller.is_exist_error__caller_sql == true {
		_pt_bp_config.T_db.T_caller.Print__sql_errors()
	}

	// gen code
	t.code = &GenCode{}
	code, err = t.code.gen_source_code(_pt_bp_config, _maps_gen_config, t.data)
	if err != nil {
		return "", err
	}

	// write file
	{
		err = ioutil.WriteFile(_s_filepath, []byte(code), 0700)
		if err != nil {
			return "", err
		}
	}
	return code, nil
}
