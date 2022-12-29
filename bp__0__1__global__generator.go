package go_orm_gen

/*
func Generate__mysql__golang(
	_s_db__addr,
	_s_db__port,
	_s_db__id,
	_s_db__pw,
	_s_db__name string,
	_s_table_name__prefix string,
	_s_filepath__json,
	_s_filepath__go,
	_s_package_name string,
	_s_class_name string,
) (
	err error,
) {
	err = Generate__config_json____mysql(_s_db__addr, _s_db__port, _s_db__id, _s_db__pw, _s_db__name, _s_table_name__prefix, _s_filepath__json)
	if err != nil {
		return err
	}
	err = Generate__source_code____golang(_s_db__addr, _s_db__port, _s_db__id, _s_db__pw, _s_db__name, _s_filepath__json, _s_filepath__go, _s_package_name, _s_class_name)
	if err != nil {
		return err
	}
	return nil
}

func Generate__config_json____mysql(
	_s_db__addr,
	_s_db__port,
	_s_db__id, _s_db__pw,
	_s_db__name string,
	_s_table_name__prefix string,
	_s_filepath__json string,
) (
	err error,
) {
	var pt_bp *T_BP
	// bp init
	{
		pc_db := db.New_C_DB_conn()
		pc_db.Init()
		pc_db.DB__connect("mysql", db_mysql.Make_str_dsn(_s_db__id, _s_db__pw, _s_db__addr, _s_db__port, _s_db__name), _s_db__name)

		pt_bp = &T_BP{}
		pt_bp.Init(pc_db, TD_N1_db_rds_type__mysql)
	}

	// json schema - load & save
	{
		err = pt_bp.Json__load(_s_filepath__json)
		if err != nil {
			return err
		}
		err = pt_bp.DBMS__load__schema(_s_table_name__prefix)
		if err != nil {
			return err
		}
		err = pt_bp.Json__save(_s_filepath__json)
		if err != nil {
			return err
		}
	}
	return nil
}

func Generate__source_code____golang(
	_s_db__addr,
	_s_db__port,
	_s_db__id,
	_s_db__pw,
	_s_db__name string,
	_s_filepath__json,
	_s_filepath__go,
	_s_package_name string,
	_s_class_name string,
) error {
	var err error
	var pt_bp *T_BP
	// bp - init
	{
		pc_db := db.New_C_DB_conn()
		pc_db.Init()
		pc_db.DB__connect("mysql", db_mysql.Make_str_dsn(_s_db__id, _s_db__pw, _s_db__addr, _s_db__port, _s_db__name), _s_db__name)

		pt_bp = &T_BP{}
		pt_bp.Init(pc_db, TD_N1_db_rds_type__mysql)
	}

	// json schema - load
	{
		err = pt_bp.Json__load(_s_filepath__json)
		if err != nil {
			return err
		}
	}
	// code - generate
	{
		err = pt_bp.Code__save(_s_filepath__go, TD_S_lang_name__go, map[string]string{
			DEF_s_gen_config__go__db__package_name: _s_package_name,
			DEF_s_gen_config__go__db__class__name:  _s_class_name,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
*/
