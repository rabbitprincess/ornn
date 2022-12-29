package go_orm_gen

import (
	"fmt"
	"strings"

	"github.com/gokch/go-orm-gen/db"
	"github.com/gokch/go-orm-gen/sql"
)

type GenData struct {
	groups []*GenDataGroup

	db     *db.DB
	dbType TD_N1_db_rds_type
}

func (t *GenData) Init(db *db.DB, dbType TD_N1_db_rds_type) {
	t.db = db
	t.dbType = dbType
}

const (
	DEF_s_sql__prepare_statement__delimiter = "%"
	DEF_s_sql__prepare_statement__after     = "?"

	DEF_s_sql__tpl__delimiter = "#"
	DEF_s_sql__tpl__after     = "%s"
	DEF_s_sql__tpl__split     = "/"
)

func (t *GenData) prepare_db_schema(config *T_BP__config) (err error) {
	pt_schema := &config.T_db.T_schema
	t.db.Exec("Insert from test_table ")

	// group
	for _, group := range config.T_db.T_caller.Arrpt_group {
		genGroup := &GenDataGroup{}
		genGroup.Set(group.S_group_name)
		t.add_group(genGroup)

		// func
		for _, pt_query := range group.Arrpt_query {
			query := &GenDataQuery{}

			// set args
			{
				// tpl args ( # name # )를 배열로 추출
				{
					tpls, err := Util__export_str_between_delimiter(pt_query.S_sql, DEF_s_sql__tpl__delimiter)
					if err != nil {
						return err
					}

					for _, s_tpl := range tpls {
						tmps := strings.Split(s_tpl, DEF_s_sql__tpl__split)
						var arg_name string
						var arg_example_data string
						if len(tmps) == 1 {
							arg_name = tmps[0]
							arg_example_data = ""
						} else if len(tmps) == 2 {
							arg_name = tmps[0]
							arg_example_data = tmps[1]
						} else {
							err = fmt.Errorf("tpl format is wrong - %s", s_tpl)
							return err
						}

						query.tpl.set_key_value(arg_name, arg_example_data)
					}
				}

				// args ( % name % )를 배열로 추출
				{
					args, err := Util__export_str_between_delimiter(pt_query.S_sql, DEF_s_sql__prepare_statement__delimiter)
					if err != nil {
						return err
					}
					query.arg.set_key(args)
				}

			}

			// %arg% -> ? # # +  /
			sqlAfterArg := Util__replace_str__between_delimiter(pt_query.S_sql, DEF_s_sql__prepare_statement__delimiter, DEF_s_sql__prepare_statement__after)

			// 쿼리 분석 후 struct 화
			// #tpl# -> tpl
			sqlAfterArgClearTpl := Util__Replace_str__in_delimiter_value(sqlAfterArg, DEF_s_sql__tpl__delimiter, DEF_s_sql__tpl__split)

			i_sql, err := sql.New(sqlAfterArgClearTpl)
			if err != nil {
				config.T_db.T_caller.is_exist_error__caller_sql = true
				pt_query.S_tmp_err__caller__parser = fmt.Sprintf("%v", err)
				continue
			}

			switch data_parser := i_sql.(type) {
			case *sql.Select:
				err = t.proc__select(config, group, pt_query, query, data_parser)
			case *sql.Insert:
				err = t.proc__insert(config, group, pt_query, query, data_parser, pt_schema)
			case *sql.Update:
				err = t.proc__update(config, group, pt_query, query, data_parser, pt_schema)
			case *sql.Delete:
				err = t.proc__delete(config, group, pt_query, query, data_parser)
			}

			// 에러처리 - 사용자 입력 sql 문을 디버깅 하기 위해 query 에 대한 err msg 를 caller json 안에 넣어 제공
			if err != nil {
				config.T_db.T_caller.is_exist_error__caller_sql = true
				pt_query.S_tmp_err__caller__query = fmt.Sprintf("%v", err)
				continue
			}

			// pt_gen__query 데이터 구성 후처리
			{
				// 그룹 이름 복사
				query.groupName = group.S_group_name

				// 쿼리 이름 복사
				query.queryName = pt_query.S_query_name

				// sql 문 복사 ( #이름# -> %s 로 변경 )
				s_sql__after_arg_tpl := Util__replace_str__between_delimiter(sqlAfterArg, DEF_s_sql__tpl__delimiter, DEF_s_sql__tpl__after)
				query.query = s_sql__after_arg_tpl

				// group list 에 func 추가
				genGroup.add_query(query)
			}
		}
	}

	return nil
}

func (t *GenData) proc__select(
	_pt_bp_config *T_BP__config,
	_pt_group *T_BP__config__db__caller__group,
	_pt_query *T_BP__config__db__caller__query,
	_pt_gen__query *GenDataQuery,
	_pt_select *sql.Select,
) error {

	pt_rds := &T_DB_RDS{}
	pt_rds.Init(t.db, TD_N1_db_rds_type__mysql)

	_pt_gen__query.queryType = QueryTypeSelect

	// 필드 정보를 얻어온다.
	{
		s_sql, _ := Util__split_by_delimiter(_pt_query.S_sql, "where")
		s_sql__after_arg := Util__replace_str__between_delimiter(s_sql, DEF_s_sql__prepare_statement__delimiter, DEF_s_sql__prepare_statement__after)
		s_sql__after_arg_clear_tpl := Util__Replace_str__in_delimiter_value(s_sql__after_arg, DEF_s_sql__tpl__delimiter, DEF_s_sql__tpl__split)

		pt_cols_arr, err := t.db.Query(s_sql__after_arg_clear_tpl).Cols__info__arr()
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
			_pt_gen__query.ret.set_key_value(s_name, s_field_type)
		}
	}
	// single select 처리
	// 코드 생성 시 단일 구조체 반환 목적
	if _pt_select.Limit != nil && *(_pt_select.Limit) == 1 {
		_pt_gen__query.isSelectSingle = true
	}
	return nil
}

func (t *GenData) proc__insert(
	_pt_bp_config *T_BP__config,
	_pt_group *T_BP__config__db__caller__group,
	_pt_query *T_BP__config__db__caller__query,
	_pt_gen__query *GenDataQuery,
	_pt_insert *sql.Insert,
	_pt_schema *T_BP__config__db__schema,
) error {

	pt_rds := &T_DB_RDS{}
	pt_rds.Init(t.db, TD_N1_db_rds_type__mysql)

	_pt_gen__query.queryType = QueryTypeInsert

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
			_pt_gen__query.arg.set_key_value(field.FldName, s_field_type__bp)
		}
	}
	// multi insert 처리
	_pt_gen__query.isInsertMulti = _pt_query.IS_insert__multi

	return nil
}

func (t *GenData) proc__update(
	_pt_bp_config *T_BP__config,
	_pt_group *T_BP__config__db__caller__group,
	_pt_query *T_BP__config__db__caller__query,
	_pt_gen__query *GenDataQuery,
	_pt_update *sql.Update,
	_pt_schema *T_BP__config__db__schema,
) error {

	pt_rds := &T_DB_RDS{}
	pt_rds.Init(t.db, TD_N1_db_rds_type__mysql)

	_pt_gen__query.queryType = QueryTypeUpdate

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

			_pt_gen__query.arg.set_key_value(pt_field_value.FldName, s_field_type__bp)
		}
	}
	// update 시 null 값 ignore 처리
	_pt_gen__query.isUpdateNullIgnore = _pt_query.IS_update__null_ignore

	return nil
}

func (t *GenData) proc__delete(
	_pt_bp_config *T_BP__config,
	_pt_group *T_BP__config__db__caller__group,
	_pt_query *T_BP__config__db__caller__query,
	_pt_gen__query *GenDataQuery,
	_pt_delete *sql.Delete,
) error {

	_pt_gen__query.queryType = QueryTypeDelete

	// 임시 - 할게 없음
	return nil
}

//---------------------------------------------------------------------------------------------------//

func (t *GenData) add_group(_pt *GenDataGroup) {
	if t.groups == nil {
		t.groups = make([]*GenDataGroup, 0, 10)
	}
	t.groups = append(t.groups, _pt)
}

// ------------------------------------------------------------------------------------------------------------//
type GenDataGroup struct {
	Name   string
	Querys []*GenDataQuery
}

func (t *GenDataGroup) Set(_s_group_name string) {
	if t.Querys == nil {
		t.Querys = make([]*GenDataQuery, 0, 10)
	}
	t.Name = _s_group_name
}

func (t *GenDataGroup) add_query(_pt *GenDataQuery) {
	t.Querys = append(t.Querys, _pt)
}

//------------------------------------------------------------------------------------------------------------//

type QueryType int8

const (
	QueryTypeSelect QueryType = iota + 1
	QueryTypeInsert
	QueryTypeUpdate
	QueryTypeDelete
)

type GenDataQuery struct {
	queryType QueryType
	groupName string
	queryName string
	query     string

	tpl genDataStruct
	arg genDataStruct
	ret genDataStruct

	isSelectSingle     bool
	isInsertMulti      bool
	isUpdateNullIgnore bool
}

// ------------------------------------------------------------------------------------------------------------//

type Pair struct {
	Key   string
	Value string
}

type genDataStruct struct {
	arrpt_pair []*Pair
}

func (t *genDataStruct) set_key(Keys []string) {
	for _, s_field_name := range Keys {
		t.set_key_value(s_field_name, "")
	}
}

func (t *genDataStruct) set_key_value(key string, valueNew string) {
	if t.arrpt_pair == nil {
		t.arrpt_pair = make([]*Pair, 0, 10)
	}

	for _, pt_field_type := range t.arrpt_pair {
		if pt_field_type.Key == key {
			pt_field_type.Value = valueNew
			return
		}
	}

	pt_field_type := &Pair{}
	pt_field_type.Key = key
	pt_field_type.Value = valueNew

	t.arrpt_pair = append(t.arrpt_pair, pt_field_type)
}
