package go_orm_gen

import (
	"fmt"
	"module/db"
)

type T_BP struct {
	pc_db             *db.C_DB_conn
	td_n1_db_rds_type TD_N1_db_rds_type

	pt_json *T_BP__config
}

func (t *T_BP) Init(_pc_db *db.C_DB_conn, _td_n1_db_rds_type TD_N1_db_rds_type) {
	// config
	t.pc_db = _pc_db
	t.td_n1_db_rds_type = _td_n1_db_rds_type
}

func (t *T_BP) Json__load(_s_filepath string) error {
	var err error

	t.pt_json = &T_BP__config{}
	err = t.pt_json.File__load(_s_filepath)
	if err != nil {
		return err
	}

	return nil
}

func (t *T_BP) Json__save(_s_filepath string) error {
	var err error
	if t.pt_json == nil {
		return fmt.Errorf("json is emtpy")
	}

	err = t.pt_json.File__save(_s_filepath)
	if err != nil {
		return err
	}

	return nil
}

func (t *T_BP) DBMS__load__schema(_s_table_name___prefix string) error {
	pt_db__schema := &T_DB_RDS{}
	pt_db__schema.Init(t.pc_db, t.td_n1_db_rds_type)

	pt_json_db_schema__from_db, err := pt_db__schema.get_schema()
	if err != nil {
		return err
	}
	// 1. arrs_table_name__prefix 가 존재할 시 해당 prefix 를 가지고 있는 테이블 스키마만 생성
	// 2. custom type 은 업데이트 되지 않음
	err = t.pt_json.T_db.T_schema.update_schema_table(pt_json_db_schema__from_db, _s_table_name___prefix)
	if err != nil {
		return err
	}
	return nil
}

func (t *T_BP) Code__save(_s_filepath string,
	_td_lang_name LangType,
	_maps_config map[string]string,
) (
	err error,
) {
	// 전처리
	{
		if _maps_config == nil {
			return fmt.Errorf("config is emtpy")
		}
		if t.pt_json == nil {
			return fmt.Errorf("json is emtpy")
		}
	}

	t_bp_gen := &Gen{}
	t.pt_json.T_db.T_caller.is_exist_error__caller_sql = false // 초기화
	_, err = t_bp_gen.Get_source_code(
		t.pc_db,
		t.td_n1_db_rds_type,
		t.pt_json,
		_td_lang_name,
		_maps_config,
		_s_filepath)
	if err != nil {
		return err
	}

	if t.pt_json.T_db.T_caller.is_exist_error__caller_sql == true {
		// 에러 출력
		t.pt_json.T_db.T_caller.Print__sql_errors()
	}

	return nil
}
