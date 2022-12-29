package go_orm_gen

import (
	"fmt"
	"module/db"

	"strings"

	"github.com/gokch/go-orm-gen/sql"
)

type GenData struct {
	arrpt_group []*T_BP__gen__data__group

	pc_db             *db.C_DB_conn
	td_n1_db_rds_type TD_N1_db_rds_type
}

func (t *GenData) Init(_pc_db *db.C_DB_conn, _td_n1_db_rds_type TD_N1_db_rds_type) {
	t.pc_db = _pc_db
	t.td_n1_db_rds_type = _td_n1_db_rds_type
}

const (
	DEF_s_sql__prepare_statement__delimiter = "%"
	DEF_s_sql__prepare_statement__after     = "?"

	DEF_s_sql__tpl__delimiter = "#"
	DEF_s_sql__tpl__after     = "%s"
	DEF_s_sql__tpl__split     = "/"
)

func (t *GenData) prepare_db_schema(_pt_bp_config *T_BP__config) (err error) {
	pt_schema := &_pt_bp_config.T_db.T_schema
	t.pc_db.Exec("Insert from test_table ")

	// group
	for _, pt_group := range _pt_bp_config.T_db.T_caller.Arrpt_group {
		var pt_gen__group *T_BP__gen__data__group
		{
			pt_gen__group = &T_BP__gen__data__group{}
			pt_gen__group.Set(pt_group.S_group_name)
			t.add_group(pt_gen__group)
		}

		// func
		for _, pt_query := range pt_group.Arrpt_query {
			pt_gen__query := &T_BP__gen__data__query{}

			// set args
			{
				// tpl args ( # name # )를 배열로 추출
				{
					arrs_tpl, err := Util__export_str_between_delimiter(pt_query.S_sql, DEF_s_sql__tpl__delimiter)
					if err != nil {
						return err
					}

					for _, s_tpl := range arrs_tpl {
						arrs_tmp := strings.Split(s_tpl, DEF_s_sql__tpl__split)
						var arg_name string
						var arg_example_data string
						if len(arrs_tmp) == 1 {
							arg_name = arrs_tmp[0]
							arg_example_data = ""
						} else if len(arrs_tmp) == 2 {
							arg_name = arrs_tmp[0]
							arg_example_data = arrs_tmp[1]
						} else {
							err = fmt.Errorf("tpl format is wrong - %s", s_tpl)
							return err
						}

						pt_gen__query.t_tpl.set_key_value(arg_name, arg_example_data)
					}
				}

				// args ( % name % )를 배열로 추출
				{
					arrs_arg, err := Util__export_str_between_delimiter(pt_query.S_sql, DEF_s_sql__prepare_statement__delimiter)
					if err != nil {
						return err
					}
					pt_gen__query.t_arg.set_key(arrs_arg)
				}

			}

			// %arg% -> ? # # +  /
			s_sql__after_arg := Util__replace_str__between_delimiter(pt_query.S_sql, DEF_s_sql__prepare_statement__delimiter, DEF_s_sql__prepare_statement__after)

			// 쿼리 분석 후 struct 화
			var i_sql sql.SQL
			{
				// #tpl# -> tpl
				s_sql__after_arg_clear_tpl := Util__Replace_str__in_delimiter_value(s_sql__after_arg, DEF_s_sql__tpl__delimiter, DEF_s_sql__tpl__split)

				i_sql, err = sql.New(s_sql__after_arg_clear_tpl)
				if err != nil {
					_pt_bp_config.T_db.T_caller.is_exist_error__caller_sql = true
					pt_query.S_tmp_err__caller__parser = fmt.Sprintf("%v", err)
					continue
				}
			}

			switch data_parser := i_sql.(type) {
			case *sql.Select:
				err = t.proc__select(_pt_bp_config, pt_group, pt_query, pt_gen__query, data_parser)
			case *sql.Insert:
				err = t.proc__insert(_pt_bp_config, pt_group, pt_query, pt_gen__query, data_parser, pt_schema)
			case *sql.Update:
				err = t.proc__update(_pt_bp_config, pt_group, pt_query, pt_gen__query, data_parser, pt_schema)
			case *sql.Delete:
				err = t.proc__delete(_pt_bp_config, pt_group, pt_query, pt_gen__query, data_parser)
			}

			// 에러처리 - 사용자 입력 sql 문을 디버깅 하기 위해 query 에 대한 err msg 를 caller json 안에 넣어 제공
			if err != nil {
				_pt_bp_config.T_db.T_caller.is_exist_error__caller_sql = true
				pt_query.S_tmp_err__caller__query = fmt.Sprintf("%v", err)
				continue
			}

			// pt_gen__query 데이터 구성 후처리
			{
				// 그룹 이름 복사
				pt_gen__query.s_group_name = pt_group.S_group_name

				// 쿼리 이름 복사
				pt_gen__query.s_query_name = pt_query.S_query_name

				// sql 문 복사 ( #이름# -> %s 로 변경 )
				s_sql__after_arg_tpl := Util__replace_str__between_delimiter(s_sql__after_arg, DEF_s_sql__tpl__delimiter, DEF_s_sql__tpl__after)
				pt_gen__query.s_query = s_sql__after_arg_tpl

				// group list 에 func 추가
				pt_gen__group.add_query(pt_gen__query)
			}
		}
	}

	return nil
}

func (t *GenData) proc__select(
	_pt_bp_config *T_BP__config,
	_pt_group *T_BP__config__db__caller__group,
	_pt_query *T_BP__config__db__caller__query,
	_pt_gen__query *T_BP__gen__data__query,
	_pt_select *sql.Select,
) error {

	pt_rds := &T_DB_RDS{}
	pt_rds.Init(t.pc_db, TD_N1_db_rds_type__mysql)

	_pt_gen__query.TD_n1_query_type = TD_N1_query_type__select

	// 필드 정보를 얻어온다.
	{
		s_sql, _ := Util__split_by_delimiter(_pt_query.S_sql, "where")
		s_sql__after_arg := Util__replace_str__between_delimiter(s_sql, DEF_s_sql__prepare_statement__delimiter, DEF_s_sql__prepare_statement__after)
		s_sql__after_arg_clear_tpl := Util__Replace_str__in_delimiter_value(s_sql__after_arg, DEF_s_sql__tpl__delimiter, DEF_s_sql__tpl__split)

		pt_cols_arr, err := t.pc_db.Query(s_sql__after_arg_clear_tpl).Cols__info__arr()
		if err != nil {
			return err
		}

		for _, pt := range pt_cols_arr {
			s_name := pt.S_name

			var s_field_type string
			s_field_type__custom := _pt_query.Get_select_field_type(s_name)
			if s_field_type__custom != "" {
				s_field_type = s_field_type__custom
			} else {
				s_type := pt.PT_type.DatabaseTypeName()
				s_field_type__bp := string(pt_rds.i_dbms.conv_field_type__to_bp(s_type))
				s_field_type = s_field_type__bp
			}
			_pt_gen__query.t_ret.set_key_value(s_name, s_field_type)
		}
	}
	// single select 처리
	// 코드 생성 시 단일 구조체 반환 목적
	if _pt_select.Limit != nil && *(_pt_select.Limit) == 1 {
		_pt_gen__query.is_select__single = true
	}
	return nil
}

func (t *GenData) proc__insert(
	_pt_bp_config *T_BP__config,
	_pt_group *T_BP__config__db__caller__group,
	_pt_query *T_BP__config__db__caller__query,
	_pt_gen__query *T_BP__gen__data__query,
	_pt_insert *sql.Insert,
	_pt_schema *T_BP__config__db__schema,
) error {

	pt_rds := &T_DB_RDS{}
	pt_rds.Init(t.pc_db, TD_N1_db_rds_type__mysql)

	_pt_gen__query.TD_n1_query_type = TD_N1_query_type__insert

	// 필드 정보를 얻어온다.
	{
		pt_schema__table := _pt_schema.get_table(_pt_insert.TblName)
		if pt_schema__table == nil {
			return fmt.Errorf("table name is not exist | table name - %s", _pt_insert.TblName)
		}

		// 스키마와 파서의 전체 필드 숫자가 다르면 -> 파서에서 모든 필드 이름이 제공되어야 함 -> 하나라도 없으면 에러
		if len(_pt_insert.Fields) != len(pt_schema__table.Arrpt_field) {
			for _, pt_field_value := range _pt_insert.Fields {
				if pt_field_value.FldName == "" {
					return fmt.Errorf("field name is empty")
				}
			}
		} else {
			// 스키마와 파서의 전체 필드수가 같으면 -> 파서에서 모든 필드 이름이 없어도 가능 -> 스키마에서 추출하여 모든 필드명을 채움
			for i, field := range _pt_insert.Fields {
				field.FldName = pt_schema__table.Arrpt_field[i].S_name
			}
		}

		// 필드 이름을 모두 채운 상태에서 처리 시작
		for _, field := range _pt_insert.Fields {
			// 입력값이 ? (arg) 형식이 아니면 func arg 를 만들 필요가 없음으로 continue
			if Util__is_parser_val__arg(field.Val) == false {
				continue
			}

			// 입력값이 ? (arg) 일 때만 필드이름 조사 = func arg 의 name 으로 활용
			pt_schema__field := pt_schema__table.get_field(field.FldName)
			if pt_schema__field == nil {
				return fmt.Errorf("not exist field in schema | field name : %s", field.FldName)
			}

			s_field_type__bp := string(pt_schema__field.S_type__BP)
			_pt_gen__query.t_arg.set_key_value(field.FldName, s_field_type__bp)
		}
	}
	// multi insert 처리
	_pt_gen__query.is_insert__multi = _pt_query.IS_insert__multi

	return nil
}

func (t *GenData) proc__update(
	_pt_bp_config *T_BP__config,
	_pt_group *T_BP__config__db__caller__group,
	_pt_query *T_BP__config__db__caller__query,
	_pt_gen__query *T_BP__gen__data__query,
	_pt_update *sql.Update,
	_pt_schema *T_BP__config__db__schema,
) error {

	pt_rds := &T_DB_RDS{}
	pt_rds.Init(t.pc_db, TD_N1_db_rds_type__mysql)

	_pt_gen__query.TD_n1_query_type = TD_N1_query_type__update

	// set
	{
		for _, pt_field_value := range _pt_update.Field {
			// 입력값이 ? (arg) 형식이 아니면 func arg 를 만들 필요가 없음으로 continue
			if Util__is_parser_val__arg(pt_field_value.Val) == false {
				continue
			}

			s_field_name := pt_field_value.FldName
			s_table_name := pt_field_value.TblName

			// 정의된 table name 이 없으면 update 대상 테이블 중 매칭되는 테이블을 찾는다
			if s_table_name == "" {
				tables := _pt_update.GetTableNames()
				tablesMatch, err := _pt_schema.get_table_list__have__field_name(s_field_name, tables)
				if err != nil {
					return err
				}

				// parse 에러 처리
				{
					// 두개 이상의 테이블이 매칭됨
					if len(tablesMatch) > 1 {
						var s_table_name_dup string
						for _, s_table_name__match := range tablesMatch {
							s_table_name_dup += fmt.Sprintf("%s, ", s_table_name__match)
						}
						s_table_name_dup = s_table_name_dup[:len(s_table_name_dup)-2]
						return fmt.Errorf("duplicated field name in multiple table | field name - %s | tables name - %s", s_field_name, s_table_name_dup)
					}
					// 매칭되는 테이블이 한개도 없음
					if len(tablesMatch) == 0 {
						return fmt.Errorf("no tables match the field | field name - %s", s_field_name)
					}
				}

				// 테이블 이름 설정
				s_table_name = tablesMatch[0]
			}

			// 테이블과 필드 이름을 이용해 필드 타입을 찾아낸다
			var s_field_type__bp string
			{
				pt_schema__table := _pt_schema.get_table(s_table_name)
				if pt_schema__table == nil {
					return fmt.Errorf("not exist table | table name - %s", s_table_name)
				}
				pt_schema__field := pt_schema__table.get_field(s_field_name)
				if pt_schema__field == nil {
					return fmt.Errorf("not exist field | field name - %s", pt_field_value.FldName)
				}
				s_field_type__bp = string(pt_schema__field.S_type__BP)
			}

			_pt_gen__query.t_arg.set_key_value(pt_field_value.FldName, s_field_type__bp)
		}
	}
	// update 시 null 값 ignore 처리
	_pt_gen__query.is_update__null_ignore = _pt_query.IS_update__null_ignore

	return nil
}

func (t *GenData) proc__delete(
	_pt_bp_config *T_BP__config,
	_pt_group *T_BP__config__db__caller__group,
	_pt_query *T_BP__config__db__caller__query,
	_pt_gen__query *T_BP__gen__data__query,
	_pt_delete *sql.Delete,
) error {

	_pt_gen__query.TD_n1_query_type = TD_N1_query_type__delete

	// 임시 - 할게 없음
	return nil
}

//---------------------------------------------------------------------------------------------------//
