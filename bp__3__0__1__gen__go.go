package bp

import (
	"fmt"
	"module/debug/logf"
	"strings"
)

const (
	DEF_s_gen_config__go__comment__do_not_edit string = "// Code generated - DO NOT EDIT.\n// This file is a generated and any changes will be lost.\n"

	DEF_s_gen_config__go__db__package_name   string = "package_name"
	DEF_s_gen_config__go__db__class__name    string = "class_name"
	DEF_s_gen_config__go__db__instance__type string = "db_instance_type"

	DEF_s_gen_config__default__go__db__package_name  string = "bp_result"
	DEF_s_gen_config__default__go__db__class_name    string = "C_DB"
	DEF_s_gen_config__default__go__db__instance_type string = "*C_DB_job"
	DEF_s_gen_config__go__db__instance_name          string = "pc_db_job"
	DEF_s_gen_config__go__db__class__prefix          string = "C_"
	DEF_s_gen_config__go__db__struct__prefix         string = "T_"
	DEF_s_gen_config__go__db__func__struct_var_name  string = "t"
	DEF_s_gen_config__go__db__func__arg__prefix      string = "_"
	DEF_s_gen_config__go__db__func__arg__prefix__tpl string = "tpl__"
	DEF_s_gen_config__go__db__func__arg__prefix__arg string = "arg__"

	DEF_s_gen_config__go__db__func__in_arg__name__prefix string = "i_"
	DEF_s_gen_config__go__db__func__in_arg__type         string = "interface{}"

	DEF_s_gen_config__go__db__func__body__func_name__SQL_remove__update_set_field__null string = "SQL_remove__update_set_field__null"

	DEF_s_gen_config__go__db__func__ret__error__type  string = "error"
	DEF_s_gen_config__go__db__func__ret__error__name  string = "err"
	DEF_s_gen_config__go__db__struct__func_name__init string = "Init"
)

type T_BP__gen__go struct {
	pt_gen_go *T_Gen__go

	s_cfg__comment__do_not_edit string
	s_cfg__db__package__name    string
	s_cfg__db__class__name      string
	s_cfg__db__instance__type   string

	pt_bp_config *T_BP__config
}

func (t *T_BP__gen__go) init(
	_pt_bp_config *T_BP__config,
	_maps_config map[string]string,
) {
	t.pt_bp_config = _pt_bp_config
	t.pt_gen_go = &T_Gen__go{}

	// config
	if _maps_config == nil {
		_maps_config = make(map[string]string)
	}

	if _maps_config[DEF_s_gen_config__go__comment__do_not_edit] == "" {
		_maps_config[DEF_s_gen_config__go__comment__do_not_edit] = DEF_s_gen_config__go__comment__do_not_edit
	}
	if _maps_config[DEF_s_gen_config__go__db__package_name] == "" {
		_maps_config[DEF_s_gen_config__go__db__package_name] = DEF_s_gen_config__default__go__db__package_name
	}
	if _maps_config[DEF_s_gen_config__go__db__class__name] == "" {
		_maps_config[DEF_s_gen_config__go__db__class__name] = DEF_s_gen_config__default__go__db__class_name
	}
	if _maps_config[DEF_s_gen_config__go__db__instance__type] == "" {
		_maps_config[DEF_s_gen_config__go__db__instance__type] = DEF_s_gen_config__default__go__db__instance_type
	}

	t.s_cfg__comment__do_not_edit = _maps_config[DEF_s_gen_config__go__comment__do_not_edit]
	t.s_cfg__db__package__name = _maps_config[DEF_s_gen_config__go__db__package_name]
	t.s_cfg__db__class__name = _maps_config[DEF_s_gen_config__go__db__class__name]
	t.s_cfg__db__instance__type = _maps_config[DEF_s_gen_config__go__db__instance__type]
	return
}

func (t *T_BP__gen__go) gen_source_code(
	_pt_bp_config *T_BP__config,
	_maps_config map[string]string,
	_pt_gen_data *T_BP__gen__data,
) (
	s_code string,
	err error,
) {
	// config 설정
	t.init(_pt_bp_config, _maps_config)

	// gen_go 에 소스 생성을 위한 데이터 넣기
	{
		t.pt_gen_go.T_global.S_comment__do_not_edit = t.s_cfg__comment__do_not_edit

		t.pt_gen_go.T_global.S_package_name = t.s_cfg__db__package__name

		// import 경로 추가
		{
			var pt_import__item *T_Go__import__item
			for _, pt_lang := range t.pt_bp_config.T_gen.Arrpt_lang_type {
				if pt_lang.TD_s_lang_name != "go" {
					continue
				}

				for _, pt_import := range pt_lang.Arrpt_import {
					pt_import__item = &T_Go__import__item{}
					pt_import__item.S_path = pt_import.S_path
					pt_import__item.S_use_name = pt_import.S_name

					t.pt_gen_go.T_global.t_import.Add(pt_import__item)
				}
			}
		}

		// group 을 담을 root 구조체 생성
		pt_root := &T_Go__struct{}
		t.pt_gen_go.T_global.Add(pt_root)
		pt_root__func__init := &T_Go__func{}
		t.pt_gen_go.T_global.Add(pt_root__func__init)

		// root
		{
			// 구조체 작성
			{
				pt_root.t_var_field.td_n1_code_format = TD_N1_code_format__var__in_struct_field
				pt_root.S_struct_name = t.s_cfg__db__class__name
			}

			// init 함수 작성
			{
				pt_root__func__init.t_var_arg.td_n1_code_format = TD_N1_code_format__var__in_func_arg
				pt_root__func__init.t_var_return.td_n1_code_format = TD_N1_code_format__var__in_func_return
				pt_root__func__init.S_func_name = DEF_s_gen_config__go__db__struct__func_name__init

				pt_root__func__init.S_struct_name = DEF_s_gen_config__go__db__func__struct_var_name
				pt_root__func__init.S_struct_type = fmt.Sprintf("*%s", pt_root.S_struct_name)

				// arg
				pt_root__func__init__var := &T_Go__var__item{}
				{
					pt_root__func__init__var.S_var_name = fmt.Sprintf("%s%s",
						DEF_s_gen_config__go__db__func__arg__prefix,
						strings.ToLower(DEF_s_gen_config__go__db__instance_name))
					pt_root__func__init__var.S_var_type = t.s_cfg__db__instance__type
					pt_root__func__init.t_var_arg.Add(pt_root__func__init__var)
				}

				// return
				{
					// 없음 - err 처리 할 것이 없음
				}
			}
		}

		// group 단위 구조체
		for _, pt_gen_group := range _pt_gen_data.arrpt_group {
			// group 구조체 생성
			pt_group := t.gen_source_code__go__group(pt_gen_group.s_group_name)
			t.pt_gen_go.T_global.Add(pt_group)

			// root 구조체 안에 필드 변수 선언 -> group 구조체 사용을 위해
			{
				// root 구조체 안에 group 구조체 포인터 선언
				pt_root__vars := &T_Go__var__item{}
				pt_root__vars.S_var_type = pt_group.S_struct_name
				pt_root__vars.S_var_name = fmt.Sprintf("%s%s", DEF_s_gen_config__go__db__class__prefix, strings.ToLower(pt_gen_group.s_group_name))
				pt_root.t_var_field.Add(pt_root__vars)

				// root init body 작성
				{
					s_code := fmt.Sprintf("%s.%s.%s(%s)\n",
						DEF_s_gen_config__go__db__func__struct_var_name,
						pt_root__vars.S_var_name,
						DEF_s_gen_config__go__db__struct__func_name__init,
						pt_root__func__init.t_var_arg.arrpt[0].S_var_name,
					)
					pt_root__func__init.S_in_body_code += s_code
				}
			}

			// group 구조체 안에 query 함수 생성
			for _, pt_gen_query := range pt_gen_group.arrpt_query {
				t.gen_query(pt_group, pt_gen_query)
			}
		} // end of for pt_group
	}

	// 실제 소스 출력
	s_code, err = t.pt_gen_go.Code()
	if err != nil {
		return "", err
	}

	return s_code, nil
}

func (t *T_BP__gen__go) gen_source_code__go__group(_s_group_name string) (pt_group *T_Go__struct) {
	pt_group = &T_Go__struct{}

	pt_group.t_var_field.td_n1_code_format = TD_N1_code_format__var__in_struct_field
	pt_group.S_struct_name = fmt.Sprintf("%s%s", DEF_s_gen_config__go__db__class__prefix, Util__conv_first_upper_case(_s_group_name))

	// group 구조체 안에
	{
		// root 구조체 연결을 위한 구조체 필드 변수 제작
		pt_group__var := &T_Go__var__item{}
		pt_group.t_var_field.Add(pt_group__var)
		{
			pt_group__var.S_var_name = DEF_s_gen_config__go__db__instance_name
			pt_group__var.S_var_type = t.s_cfg__db__instance__type
		}

		// root 구조체에서 초기화를 요청할 Init 함수 제작
		pt_group__func__init := &T_Go__func{}
		t.pt_gen_go.T_global.Add(pt_group__func__init)
		{
			pt_group__func__init.t_var_arg.td_n1_code_format = TD_N1_code_format__var__in_func_arg
			pt_group__func__init.t_var_return.td_n1_code_format = TD_N1_code_format__var__in_func_return
			pt_group__func__init.S_func_name = DEF_s_gen_config__go__db__struct__func_name__init

			pt_group__func__init.S_struct_name = DEF_s_gen_config__go__db__func__struct_var_name
			pt_group__func__init.S_struct_type = fmt.Sprintf("*%s", pt_group.S_struct_name)

			// init 함수 입력 인자
			pt_group__func__init__arg := &T_Go__var__item{}
			pt_group__func__init.t_var_arg.Add(pt_group__func__init__arg)
			{
				pt_group__func__init__arg.S_var_name = fmt.Sprintf("%s%s",
					DEF_s_gen_config__go__db__func__arg__prefix,
					strings.ToLower(DEF_s_gen_config__go__db__instance_name))
				pt_group__func__init__arg.S_var_type = t.s_cfg__db__instance__type

			}

			// init 함수 리턴 정의
			{
				// 없음 - err 처리 할 것이 없음
			}

			pt_group__func__init.S_in_body_code = fmt.Sprintf("%s.%s = %s",
				pt_group__func__init.S_struct_name,
				pt_group__var.S_var_name,
				pt_group__func__init__arg.S_var_name)
		}
	}
	return pt_group
}

func (t *T_BP__gen__go) gen_query(
	_pt_group *T_Go__struct,
	_pt_gen_query *T_BP__gen__data__query,
) {
	pt_query := &T_Go__func{}

	pt_query.t_var_arg.td_n1_code_format = TD_N1_code_format__var__in_func_arg
	pt_query.t_var_return.td_n1_code_format = TD_N1_code_format__var__in_func_return
	pt_query.S_struct_name = DEF_s_gen_config__go__db__func__struct_var_name
	pt_query.S_struct_type = fmt.Sprintf("*%s", _pt_group.S_struct_name)
	pt_query.S_func_name = Util__conv_first_upper_case(_pt_gen_query.s_query_name)

	switch _pt_gen_query.TD_n1_query_type {
	case TD_N1_query_type__select:
		{
			t.gen_query__select(pt_query, _pt_gen_query)
		}
	case TD_N1_query_type__insert:
		{
			t.gen_query__insert(pt_query, _pt_gen_query)
		}
	case TD_N1_query_type__update:
		{
			t.gen_query__update(pt_query, _pt_gen_query)
		}
	case TD_N1_query_type__delete:
		{
			t.gen_query__delete(pt_query, _pt_gen_query)
		}
	default:
		{
			logf.GC.Fatal("bp,gen", "invalid query type | query type : %v", _pt_gen_query.TD_n1_query_type)
		}
	}

	t.pt_gen_go.T_global.Add(pt_query)
}

func (t *T_BP__gen__go) gen_query__select(
	_pt_query *T_Go__func,
	_pt_gen_query *T_BP__gen__data__query,
) {

	// 1. 함수 입력 인자
	arrs_tpl := t.gen_query__add__func__arg__tpl(_pt_query, _pt_gen_query)
	arrs_arg := t.gen_query__add__func__arg(_pt_query, _pt_gen_query)

	// 2. 함수 리턴
	pt_query__ret := &T_Go__struct{}
	t.pt_gen_go.T_global.Add(pt_query__ret)

	pt_query__ret__struct := &T_Go__var__item{}
	var s_body_code__ret_declare string
	var s_body_code__ret_set string
	{
		// 2-1. 쿼리-리턴 변수 선언
		{
			pt_query__ret.t_var_field.td_n1_code_format = TD_N1_code_format__var__in_struct_field
			pt_query__ret.S_struct_name = fmt.Sprintf("%s%s__%s", DEF_s_gen_config__go__db__struct__prefix, Util__conv_first_upper_case(_pt_gen_query.s_group_name), strings.ToLower(_pt_query.S_func_name))
			for _, pt_field_type := range _pt_gen_query.t_ret.arrpt_pair {
				pt_var__item := &T_Go__var__item{}
				pt_var__item.S_var_name = Util__conv_first_upper_case(pt_field_type.s_key)
				pt_var__item.S_var_type = t.pt_bp_config.T_gen.Conv_field_type__bp_to_lang(pt_field_type.s_value, TD_S_lang_name__go)
				pt_query__ret.t_var_field.Add(pt_var__item)
			}
		}

		// 2-2. 리턴 변수 처리
		{
			// 리턴 변수 선언 - 구조체
			if _pt_gen_query.is_select__single == true {
				pt_query__ret__struct.S_var_name = fmt.Sprintf("pt_%s", strings.ToLower(_pt_query.S_func_name))
				pt_query__ret__struct.S_var_type = fmt.Sprintf("*%s", pt_query__ret.S_struct_name)
				s_body_code__ret_set = fmt.Sprintf("%s = pt_struct\n\tbreak", pt_query__ret__struct.S_var_name)
			} else {
				pt_query__ret__struct.S_var_name = fmt.Sprintf("arrpt_%s", strings.ToLower(_pt_query.S_func_name))
				pt_query__ret__struct.S_var_type = fmt.Sprintf("[]*%s", pt_query__ret.S_struct_name)
				s_body_code__ret_declare = fmt.Sprintf("\n%s = make(%s, 0, 100)", pt_query__ret__struct.S_var_name, pt_query__ret__struct.S_var_type)
				s_body_code__ret_set = fmt.Sprintf("%s = append(%s, pt_struct)", pt_query__ret__struct.S_var_name, pt_query__ret__struct.S_var_name)
			}
			_pt_query.t_var_return.Add(pt_query__ret__struct)

			// error 추가
			t.gen_query__add__func__ret__error(_pt_query)
		}
	}

	// 3. 함수 body
	{
		_pt_query.S_in_body_code = fmt.Sprintf(`
%s
s_sql := fmt.Sprintf(
	"%s",%s
)
pc_ret := %s.%s.Query(
	s_sql,
	arri_arg...,
)
defer pc_ret.Close()
%s
for {
	pt_struct := &%s{}
	is_end, err := pc_ret.Row_next(pt_struct)
	if err != nil {
		return nil, err
	}
	if is_end == true {
		break
	}
	%s
}

return %s, nil
`,
			t.gen_query__add__func__body__set_args(arrs_arg),
			_pt_gen_query.s_query,
			t.gen_query__add__func__body__arg(arrs_tpl),
			DEF_s_gen_config__go__db__func__struct_var_name,
			DEF_s_gen_config__go__db__instance_name,
			s_body_code__ret_declare,
			pt_query__ret.S_struct_name,
			s_body_code__ret_set,
			pt_query__ret__struct.S_var_name,
		)
	}
	return
}

func (t *T_BP__gen__go) gen_query__insert(
	_pt_query *T_Go__func,
	_pt_gen_query *T_BP__gen__data__query,
) {
	// 1. 함수 입력 인자
	arrs_tpl := t.gen_query__add__func__arg__tpl(_pt_query, _pt_gen_query)
	arrs_arg := t.gen_query__add__func__arg(_pt_query, _pt_gen_query)

	// 2. 함수 리턴 변수
	t.gen_query__add__func__ret__last_insert_id(_pt_query)
	t.gen_query__add__func__ret__error(_pt_query)

	// 3. 함수 body
	{
		// body 전처리
		var s_multi_insert__query, s_gen_args string
		if _pt_gen_query.is_insert__multi == true { // multi insert
			s_query_values := Util__export__insert_query_values(_pt_gen_query.s_query)
			if _pt_gen_query.s_query[len(_pt_gen_query.s_query)-1:] == ";" {
				_pt_gen_query.s_query = _pt_gen_query.s_query[:len(_pt_gen_query.s_query)-1]
			}
			_pt_gen_query.s_query += "%s"
			s_gen_args = t.gen_query__add__func__body__insert_multi__proc(arrs_arg)
			s_multi_insert__query = t.gen_query__add__func__body__insert_multi__query(s_query_values)
		} else { // insert
			s_gen_args = t.gen_query__add__func__body__set_args(arrs_arg)
		}

		_pt_query.S_in_body_code = fmt.Sprintf(`
%s
s_sql := fmt.Sprintf(
	"%s",%s%s
)

pc_exec := %s.%s.Exec(
	s_sql,
	arri_arg...,
)

return pc_exec.LastInsertId()
`,
			s_gen_args,
			_pt_gen_query.s_query,
			t.gen_query__add__func__body__arg(arrs_tpl),
			s_multi_insert__query,
			DEF_s_gen_config__go__db__func__struct_var_name,
			DEF_s_gen_config__go__db__instance_name,
		)
	}
	return
}

func (t *T_BP__gen__go) gen_query__update(
	_pt_query *T_Go__func,
	_pt_gen_query *T_BP__gen__data__query,
) {
	// 1. 함수 입력 인자
	arrs_tpl := t.gen_query__add__func__arg__tpl(_pt_query, _pt_gen_query)
	arrs_arg := t.gen_query__add__func__arg(_pt_query, _pt_gen_query)

	// 2. 함수 리턴 변수
	t.gen_query__add__func__ret__row_affected(_pt_query)
	t.gen_query__add__func__ret__error(_pt_query)

	// 3. 함수 body
	{
		// 전처리
		var s_body__set_args string
		if _pt_gen_query.is_update__null_ignore == true {
			s_body__set_args = t.gen_query__add__func__body__set_args__remove_sets(arrs_arg)
		} else {
			s_body__set_args = t.gen_query__add__func__body__set_args(arrs_arg)
		}
		_pt_query.S_in_body_code = fmt.Sprintf(`
s_sql := fmt.Sprintf(
	"%s",%s
)
%s
pc_exec := %s.%s.Exec(
	s_sql,
	arri_arg...,
)

return pc_exec.RowsAffected()
`,
			_pt_gen_query.s_query,
			t.gen_query__add__func__body__arg(arrs_tpl),
			s_body__set_args,
			DEF_s_gen_config__go__db__func__struct_var_name,
			DEF_s_gen_config__go__db__instance_name,
		)
	}
	return
}

func (t *T_BP__gen__go) gen_query__delete(
	_pt_query *T_Go__func,
	_pt_gen_query *T_BP__gen__data__query,
) {

	// 1. 함수 입력 인자
	arrs_tpl := t.gen_query__add__func__arg__tpl(_pt_query, _pt_gen_query)
	arrs_arg := t.gen_query__add__func__arg(_pt_query, _pt_gen_query)

	// 2. 함수 리턴 변수
	t.gen_query__add__func__ret__row_affected(_pt_query)
	t.gen_query__add__func__ret__error(_pt_query)

	// 3. 함수 body
	{
		_pt_query.S_in_body_code = fmt.Sprintf(`
%s
s_sql := fmt.Sprintf(
	"%s",%s
)
		
pc_exec := %s.%s.Exec(
	s_sql,
	arri_arg...,
)

return pc_exec.RowsAffected()
`,
			t.gen_query__add__func__body__set_args(arrs_arg),
			_pt_gen_query.s_query,
			t.gen_query__add__func__body__arg(arrs_tpl),
			DEF_s_gen_config__go__db__func__struct_var_name,
			DEF_s_gen_config__go__db__instance_name,
		)
	}
	return
}

//------------------------------------------------------------------------------------------------------//

func (t *T_BP__gen__go) gen_query__add__func__arg__tpl(
	_pt_query *T_Go__func,
	_pt_gen_query *T_BP__gen__data__query,
) (
	arrs_tpl []string,
) {
	arrs_tpl = make([]string, 0, len(_pt_gen_query.t_tpl.arrpt_pair))
	for _, pt_tpl := range _pt_gen_query.t_tpl.arrpt_pair {
		pt_query__arg := &T_Go__var__item{}
		pt_query__arg.S_var_name = fmt.Sprintf("%s%s%s",
			DEF_s_gen_config__go__db__func__arg__prefix,
			DEF_s_gen_config__go__db__func__arg__prefix__tpl,
			pt_tpl.s_key)
		pt_query__arg.S_var_type = "string"
		_pt_query.t_var_arg.Add(pt_query__arg)

		arrs_tpl = append(arrs_tpl, pt_query__arg.S_var_name)
	}
	return
}

func (t *T_BP__gen__go) gen_query__add__func__arg(
	_pt_query *T_Go__func,
	_pt_gen_query *T_BP__gen__data__query,
) (
	arrs_arg []string,
) {
	arrs_arg = make([]string, 0, len(_pt_gen_query.t_tpl.arrpt_pair))

	for _, pt_field_type := range _pt_gen_query.t_arg.arrpt_pair {
		pt_query__arg := &T_Go__var__item{}
		var s_var_type string
		// 1. type 판정
		{
			if pt_field_type.s_value == "" { // 형을 특정할 수 없을 때
				s_var_type = DEF_s_gen_config__go__db__func__in_arg__type
			} else { // 형을 특정할 수 있을 때
				s_var_type = "*" + t.pt_bp_config.T_gen.Conv_field_type__bp_to_lang(pt_field_type.s_value, TD_S_lang_name__go)
			}
			if _pt_gen_query.is_insert__multi == true {
				s_var_type = "[]" + s_var_type
			}
		}

		// 2. set query arg
		{
			pt_query__arg.S_var_name = fmt.Sprintf("%s%s%s",
				DEF_s_gen_config__go__db__func__arg__prefix,
				DEF_s_gen_config__go__db__func__arg__prefix__arg,
				pt_field_type.s_key)
			pt_query__arg.S_var_type = s_var_type
			_pt_query.t_var_arg.Add(pt_query__arg)
		}
		arrs_arg = append(arrs_arg, pt_query__arg.S_var_name)
	}
	return
}

func (t *T_BP__gen__go) gen_query__add__func__ret__error(
	_pt_query *T_Go__func,
) {
	pt_query__ret__err := &T_Go__var__item{}
	pt_query__ret__err.S_var_name = DEF_s_gen_config__go__db__func__ret__error__name
	pt_query__ret__err.S_var_type = DEF_s_gen_config__go__db__func__ret__error__type
	_pt_query.t_var_return.Add(pt_query__ret__err)
}

func (t *T_BP__gen__go) gen_query__add__func__ret__last_insert_id(
	_pt_query *T_Go__func,
) {
	pt_query__ret__last_id := &T_Go__var__item{}
	pt_query__ret__last_id.S_var_name = "nn_last_insert_id"
	pt_query__ret__last_id.S_var_type = "int64"
	_pt_query.t_var_return.Add(pt_query__ret__last_id)
}

func (t *T_BP__gen__go) gen_query__add__func__ret__row_affected(
	_pt_query *T_Go__func,
) {

	pt_query__ret__row_affected := &T_Go__var__item{}
	pt_query__ret__row_affected.S_var_name = "nn_row_affected"
	pt_query__ret__row_affected.S_var_type = "int64"
	_pt_query.t_var_return.Add(pt_query__ret__row_affected)
}

func (t *T_BP__gen__go) gen_query__add__func__body__arg(_arrs_arg []string) (s_arg string) {
	for _, s_arg__one := range _arrs_arg {
		s_arg += fmt.Sprintf("\n\t%s,", s_arg__one)
	}
	return s_arg
}

func (t *T_BP__gen__go) gen_query__add__func__body__set_args(_arrs_arg []string) (s_gen_arg string) {
	var s_gen_arg__item string
	s_gen_arg__item = t.gen_query__add__func__body__arg(_arrs_arg)
	if s_gen_arg__item != "" {
		s_gen_arg__item += "\n"
	}

	s_gen_arg += fmt.Sprintf(`arri_arg := make([]interface{}, 0, %d)
arri_arg = append(arri_arg, I_to_arri(%s)...)
`,
		len(_arrs_arg),
		s_gen_arg__item,
	)

	return s_gen_arg
}

func (t *T_BP__gen__go) gen_query__add__func__body__insert_multi__query(_s_query_values string) (s_multi_insert__query string) {
	s_multi_insert__query += fmt.Sprintf("\n\tstrings.Repeat(\", (%s)\", n_len_arg-1),", _s_query_values)
	return s_multi_insert__query
}

func (t *T_BP__gen__go) gen_query__add__func__body__insert_multi__proc(_arrs_arg []string) (s_multi_insert__body string) {
	var s_check_len string
	for i, s_arg := range _arrs_arg {
		s_check_len += fmt.Sprintf("n_len_arg != len(%s)", s_arg)
		if i != len(_arrs_arg)-1 {
			s_check_len += fmt.Sprintf(" || ")
		}
	}

	var s_append string
	for i, s_arg := range _arrs_arg {
		s_append += fmt.Sprintf("%s[i]", s_arg)
		if i != len(_arrs_arg)-1 {
			s_append += fmt.Sprintf(",\n\t\t")
		}
	}

	s_multi_insert__body = fmt.Sprintf(`n_len_arg := len(%s)
if n_len_arg == 0 {
	return 0, fmt.Errorf("arg len is zero")
}
if %s {
	return 0, fmt.Errorf("arg len is not same")
}

arri_arg := make([]interface{}, 0, n_len_arg*%d)
for i := 0; i < n_len_arg; i++ {
	arri_arg = append(arri_arg, I_to_arri(
		%s,
	)...)
}
`,
		_arrs_arg[0],
		s_check_len,
		len(_arrs_arg),
		s_append)
	return s_multi_insert__body
}

func (t *T_BP__gen__go) gen_query__add__func__body__set_args__remove_sets(_arrs_arg []string) (s_sql_update__delete_sets string) {
	var s_sql_is_nil string
	for _, s_arg := range _arrs_arg {
		s_field_name := strings.TrimPrefix(s_arg, DEF_s_gen_config__go__db__func__arg__prefix+DEF_s_gen_config__go__db__func__arg__prefix__arg)

		s_sql_is_nil += fmt.Sprintf(`if %s == nil {
	arrs_sets__removed = append(arrs_sets__removed, "%s")
} else {
	arri_arg = append(arri_arg, %s)
}
`,
			s_arg,
			s_field_name,
			s_arg,
		)
	}

	s_sql_update__delete_sets = fmt.Sprintf(`
arri_arg := make([]interface{}, 0, %d)
arrs_sets__removed := make([]string, 0, %d)
%s
if len(arrs_sets__removed) != 0 {
	s_sql, _ = %s(s_sql, arrs_sets__removed)
}
`,
		len(_arrs_arg),
		len(_arrs_arg),
		s_sql_is_nil,
		DEF_s_gen_config__go__db__func__body__func_name__SQL_remove__update_set_field__null,
	)
	return s_sql_update__delete_sets
}
