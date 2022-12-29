package go_orm_gen

import (
	"fmt"
	"log"

	"strings"

	"github.com/gokch/go-orm-gen/codegen"
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

type GenCode struct {
	codeGen *codegen.CodeGen

	s_cfg__comment__do_not_edit string
	s_cfg__db__package__name    string
	s_cfg__db__class__name      string
	s_cfg__db__instance__type   string

	pt_bp_config *T_BP__config
}

func (t *GenCode) init(bpConfig *T_BP__config, mapConfig map[string]string) {
	t.pt_bp_config = bpConfig
	t.codeGen.Global = &codegen.Global{}

	// config
	if mapConfig == nil {
		mapConfig = make(map[string]string)
	}

	if mapConfig[DEF_s_gen_config__go__comment__do_not_edit] == "" {
		mapConfig[DEF_s_gen_config__go__comment__do_not_edit] = DEF_s_gen_config__go__comment__do_not_edit
	}
	if mapConfig[DEF_s_gen_config__go__db__package_name] == "" {
		mapConfig[DEF_s_gen_config__go__db__package_name] = DEF_s_gen_config__default__go__db__package_name
	}
	if mapConfig[DEF_s_gen_config__go__db__class__name] == "" {
		mapConfig[DEF_s_gen_config__go__db__class__name] = DEF_s_gen_config__default__go__db__class_name
	}
	if mapConfig[DEF_s_gen_config__go__db__instance__type] == "" {
		mapConfig[DEF_s_gen_config__go__db__instance__type] = DEF_s_gen_config__default__go__db__instance_type
	}

	t.s_cfg__comment__do_not_edit = mapConfig[DEF_s_gen_config__go__comment__do_not_edit]
	t.s_cfg__db__package__name = mapConfig[DEF_s_gen_config__go__db__package_name]
	t.s_cfg__db__class__name = mapConfig[DEF_s_gen_config__go__db__class__name]
	t.s_cfg__db__instance__type = mapConfig[DEF_s_gen_config__go__db__instance__type]
	return
}

func (t *GenCode) gen_source_code(bpConfig *T_BP__config, mapConfig map[string]string, genData *GenData) (genCode string, err error) {
	// config 설정
	t.init(bpConfig, mapConfig)

	// gen_go 에 소스 생성을 위한 데이터 넣기
	{
		t.codeGen.Init()
		t.codeGen.DoNotEdit = t.s_cfg__comment__do_not_edit
		t.codeGen.PackageName = t.s_cfg__db__package__name

		// import 경로 추가
		for _, pt_lang := range t.pt_bp_config.T_gen.Arrpt_lang_type {
			for _, pt_import := range pt_lang.Imports {
				t.codeGen.AddImport(&codegen.ImportItem{
					Path:  pt_import.S_path,
					Alias: pt_import.S_name,
				})
			}
		}

		// 루트 구조체 작성
		rootStruct := &codegen.Struct{}
		t.codeGen.AddItem(rootStruct)
		rootStruct.Init()
		rootStruct.Name = t.s_cfg__db__class__name

		// 루트 함수 작성
		rootFunc := &codegen.Func{}
		t.codeGen.AddItem(rootFunc)
		rootFunc.Init()
		rootFunc.StructName = DEF_s_gen_config__go__db__func__struct_var_name
		rootFunc.StructType = fmt.Sprintf("*%s", rootStruct.Name)
		rootFunc.FuncName = DEF_s_gen_config__go__db__struct__func_name__init
		rootFuncInitArg := &codegen.VarItem{} // arg
		rootFunc.Arg.Add(rootFuncInitArg)
		rootFuncInitArg.Name = fmt.Sprintf("%s%s", DEF_s_gen_config__go__db__func__arg__prefix, strings.ToLower(DEF_s_gen_config__go__db__instance_name))
		rootFuncInitArg.Type = t.s_cfg__db__instance__type

		// group 단위 구조체
		for _, pt_gen_group := range genData.arrpt_group {
			// group 구조체 생성
			pt_group := t.genGroup(pt_gen_group.s_group_name)
			t.codeGen.AddItem(pt_group)

			// root 구조체 안에 필드 변수 선언 -> group 구조체 사용을 위해
			{
				// root 구조체 안에 group 구조체 포인터 선언
				rootVars := &codegen.VarItem{}
				rootVars.Type = pt_group.Name
				rootVars.Name = fmt.Sprintf("%s%s", DEF_s_gen_config__go__db__class__prefix, strings.ToLower(pt_gen_group.s_group_name))
				rootStruct.Field.Add(rootVars)

				// root init body 작성
				s_code := fmt.Sprintf("%s.%s.%s(%s)\n", DEF_s_gen_config__go__db__func__struct_var_name, rootVars.Name, DEF_s_gen_config__go__db__struct__func_name__init, rootFunc.Arg.Items[0].Name)
				rootFunc.InlineCode += s_code
			}

			// group 구조체 안에 query 함수 생성
			for _, pt_gen_query := range pt_gen_group.arrpt_query {
				t.genQuery(pt_group, pt_gen_query)
			}
		} // end of for pt_group
	}

	// 실제 소스 출력
	genCode = t.codeGen.Code()

	return genCode, nil
}

func (t *GenCode) genGroup(group string) (genGroup *codegen.Struct) {
	genGroup = &codegen.Struct{}
	genGroup.Name = fmt.Sprintf("%s%s", DEF_s_gen_config__go__db__class__prefix, Util__conv_first_upper_case(group))

	// group 구조체 안에
	{
		// root 구조체 연결을 위한 구조체 필드 변수 제작
		groupVar := &codegen.VarItem{}
		genGroup.Field.Add(groupVar)
		{
			groupVar.Name = DEF_s_gen_config__go__db__instance_name
			groupVar.Type = t.s_cfg__db__instance__type
		}

		// root 구조체에서 초기화를 요청할 Init 함수 제작
		groupFuncInit := &codegen.Func{}
		t.codeGen.AddItem(groupFuncInit)
		{
			groupFuncInit.FuncName = DEF_s_gen_config__go__db__struct__func_name__init
			groupFuncInit.StructName = DEF_s_gen_config__go__db__func__struct_var_name
			groupFuncInit.StructType = fmt.Sprintf("*%s", genGroup.Name)

			// args
			groupFuncInitArg := &codegen.VarItem{}
			groupFuncInit.Arg.Add(groupFuncInitArg)
			groupFuncInitArg.Name = fmt.Sprintf("%s%s", DEF_s_gen_config__go__db__func__arg__prefix, strings.ToLower(DEF_s_gen_config__go__db__instance_name))
			groupFuncInitArg.Type = t.s_cfg__db__instance__type

			// inline code
			groupFuncInit.InlineCode = fmt.Sprintf("%s.%s = %s", groupFuncInit.StructName, groupVar.Name, groupFuncInitArg.Name)
		}
	}
	return genGroup
}

func (t *GenCode) genQuery(structGroup *codegen.Struct, query *T_BP__gen__data__query) {
	funcQuery := &codegen.Func{}
	funcQuery.Init()

	funcQuery.StructName = DEF_s_gen_config__go__db__func__struct_var_name
	funcQuery.StructType = fmt.Sprintf("*%s", structGroup.Name)
	funcQuery.FuncName = Util__conv_first_upper_case(query.s_query_name)

	switch query.TD_n1_query_type {
	case TD_N1_query_type__select:
		t.genQuerySelect(funcQuery, query)
	case TD_N1_query_type__insert:
		t.genQueryInsert(funcQuery, query)
	case TD_N1_query_type__update:
		t.genQueryUpdate(funcQuery, query)
	case TD_N1_query_type__delete:
		t.genQueryDelete(funcQuery, query)
	default:
		log.Fatal("invalid query type | query type : %v", query.TD_n1_query_type)
	}

	t.codeGen.AddItem(funcQuery)
}

func (t *GenCode) genQuerySelect(
	funcQuery *codegen.Func,
	query *T_BP__gen__data__query,
) {
	// 1. 함수 입력 인자
	tpl := t.gen_query__add__func__arg__tpl(funcQuery, query)
	arg := t.gen_query__add__func__arg(funcQuery, query)

	// 2. 함수 리턴
	ret := &codegen.Struct{}
	ret.Init()
	t.codeGen.AddItem(ret)

	retItem := &codegen.VarItem{}
	var s_body_code__ret_declare string
	var s_body_code__ret_set string
	{
		// 2-1. 쿼리-리턴 변수 선언
		{
			ret.Name = fmt.Sprintf("%s%s__%s", DEF_s_gen_config__go__db__struct__prefix, Util__conv_first_upper_case(query.s_group_name), strings.ToLower(funcQuery.FuncName))
			for _, pt_field_type := range query.t_ret.arrpt_pair {
				item := &codegen.VarItem{}
				item.Name = Util__conv_first_upper_case(pt_field_type.s_key)
				item.Type = t.pt_bp_config.T_gen.Conv_field_type__bp_to_lang(pt_field_type.s_value, LangTypeGo)
				ret.Field.Add(item)
			}
		}

		// 2-2. 리턴 변수 처리
		{
			// 리턴 변수 선언 - 구조체
			if query.is_select__single == true {
				retItem.Name = fmt.Sprintf("pt_%s", strings.ToLower(funcQuery.FuncName))
				retItem.Type = fmt.Sprintf("*%s", ret.Name)
				s_body_code__ret_set = fmt.Sprintf("%s = pt_struct\n\tbreak", retItem.Name)
			} else {
				retItem.Name = fmt.Sprintf("arrpt_%s", strings.ToLower(funcQuery.FuncName))
				retItem.Type = fmt.Sprintf("[]*%s", ret.Name)
				s_body_code__ret_declare = fmt.Sprintf("\n%s = make(%s, 0, 100)", retItem.Name, retItem.Type)
				s_body_code__ret_set = fmt.Sprintf("%s = append(%s, pt_struct)", retItem.Name, retItem.Name)
			}
			funcQuery.Ret.Add(retItem)

			// error 추가
			t.gen_query__add__func__ret__error(funcQuery)
		}
	}

	// 3. 함수 body
	{
		funcQuery.InlineCode = fmt.Sprintf(`
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
			t.gen_query__add__func__body__set_args(arg),
			query.s_query,
			t.gen_query__add__func__body__arg(tpl),
			DEF_s_gen_config__go__db__func__struct_var_name,
			DEF_s_gen_config__go__db__instance_name,
			s_body_code__ret_declare,
			ret.Name,
			s_body_code__ret_set,
			retItem.Name,
		)
	}
	return
}

func (t *GenCode) genQueryInsert(funcQuery *codegen.Func, query *T_BP__gen__data__query) {
	// 1. 함수 입력 인자
	tpl := t.gen_query__add__func__arg__tpl(funcQuery, query)
	arg := t.gen_query__add__func__arg(funcQuery, query)

	// 2. 함수 리턴 변수
	t.gen_query__add__func__ret__last_insert_id(funcQuery)
	t.gen_query__add__func__ret__error(funcQuery)

	// 3. 함수 body
	{
		// body 전처리
		var multiInsert, genArgs string
		if query.is_insert__multi == true { // multi insert
			s_query_values := Util__export__insert_query_values(query.s_query)
			if query.s_query[len(query.s_query)-1:] == ";" {
				query.s_query = query.s_query[:len(query.s_query)-1]
			}
			query.s_query += "%s"
			genArgs = t.gen_query__add__func__body__insert_multi__proc(arg)
			multiInsert = t.gen_query__add__func__body__insert_multi__query(s_query_values)
		} else { // insert
			genArgs = t.gen_query__add__func__body__set_args(arg)
		}

		funcQuery.InlineCode = fmt.Sprintf(`
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
			genArgs,
			query.s_query,
			t.gen_query__add__func__body__arg(tpl),
			multiInsert,
			DEF_s_gen_config__go__db__func__struct_var_name,
			DEF_s_gen_config__go__db__instance_name,
		)
	}
	return
}

func (t *GenCode) genQueryUpdate(funcQuery *codegen.Func, query *T_BP__gen__data__query) {
	// 1. 함수 입력 인자
	tpl := t.gen_query__add__func__arg__tpl(funcQuery, query)
	arg := t.gen_query__add__func__arg(funcQuery, query)

	// 2. 함수 리턴 변수
	t.gen_query__add__func__ret__row_affected(funcQuery)
	t.gen_query__add__func__ret__error(funcQuery)

	// 3. 함수 body
	{
		// 전처리
		var s_body__set_args string
		if query.is_update__null_ignore == true {
			s_body__set_args = t.gen_query__add__func__body__set_args__remove_sets(arg)
		} else {
			s_body__set_args = t.gen_query__add__func__body__set_args(arg)
		}
		funcQuery.InlineCode = fmt.Sprintf(`
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
			query.s_query,
			t.gen_query__add__func__body__arg(tpl),
			s_body__set_args,
			DEF_s_gen_config__go__db__func__struct_var_name,
			DEF_s_gen_config__go__db__instance_name,
		)
	}
	return
}

func (t *GenCode) genQueryDelete(funcQuery *codegen.Func, query *T_BP__gen__data__query) {
	// 1. 함수 입력 인자
	arrs_tpl := t.gen_query__add__func__arg__tpl(funcQuery, query)
	arrs_arg := t.gen_query__add__func__arg(funcQuery, query)

	// 2. 함수 리턴 변수
	t.gen_query__add__func__ret__row_affected(funcQuery)
	t.gen_query__add__func__ret__error(funcQuery)

	// 3. 함수 body
	{
		funcQuery.InlineCode = fmt.Sprintf(`
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
			query.s_query,
			t.gen_query__add__func__body__arg(arrs_tpl),
			DEF_s_gen_config__go__db__func__struct_var_name,
			DEF_s_gen_config__go__db__instance_name,
		)
	}
	return
}

//------------------------------------------------------------------------------------------------------//

func (t *GenCode) gen_query__add__func__arg__tpl(funcQuery *codegen.Func, query *T_BP__gen__data__query) (arrs_tpl []string) {
	arrs_tpl = make([]string, 0, len(query.t_tpl.arrpt_pair))
	for _, pt_tpl := range query.t_tpl.arrpt_pair {
		arg := &codegen.VarItem{}
		funcQuery.Arg.Add(arg)
		arg.Name = fmt.Sprintf("%s%s%s", DEF_s_gen_config__go__db__func__arg__prefix, DEF_s_gen_config__go__db__func__arg__prefix__tpl, pt_tpl.s_key)
		arg.Type = "string"

		arrs_tpl = append(arrs_tpl, arg.Name)
	}
	return
}

func (t *GenCode) gen_query__add__func__arg(funcQuery *codegen.Func, query *T_BP__gen__data__query) (arrs_arg []string) {
	arrs_arg = make([]string, 0, len(query.t_tpl.arrpt_pair))

	for _, pt_field_type := range query.t_arg.arrpt_pair {
		arg := &codegen.VarItem{}
		var varType string
		// 1. type 판정
		{
			if pt_field_type.s_value == "" { // 형을 특정할 수 없을 때
				varType = DEF_s_gen_config__go__db__func__in_arg__type
			} else { // 형을 특정할 수 있을 때
				varType = "*" + t.pt_bp_config.T_gen.Conv_field_type__bp_to_lang(pt_field_type.s_value, LangTypeGo)
			}
			if query.is_insert__multi == true {
				varType = "[]" + varType
			}
		}

		// 2. set query arg
		{
			arg.Name = fmt.Sprintf("%s%s%s",
				DEF_s_gen_config__go__db__func__arg__prefix,
				DEF_s_gen_config__go__db__func__arg__prefix__arg,
				pt_field_type.s_key)
			arg.Type = varType
			funcQuery.Arg.Add(arg)
		}
		arrs_arg = append(arrs_arg, arg.Name)
	}
	return
}

func (t *GenCode) gen_query__add__func__ret__error(funcQuery *codegen.Func) {
	retErr := &codegen.VarItem{}
	retErr.Name = DEF_s_gen_config__go__db__func__ret__error__name
	retErr.Type = DEF_s_gen_config__go__db__func__ret__error__type
	funcQuery.Ret.Add(retErr)
}

func (t *GenCode) gen_query__add__func__ret__last_insert_id(funcQuery *codegen.Func) {
	retLastid := &codegen.VarItem{}
	retLastid.Name = "nn_last_insert_id"
	retLastid.Type = "int64"
	funcQuery.Ret.Add(retLastid)
}

func (t *GenCode) gen_query__add__func__ret__row_affected(funcQuery *codegen.Func) {
	pt_query__ret__row_affected := &codegen.VarItem{}
	pt_query__ret__row_affected.Name = "nn_row_affected"
	pt_query__ret__row_affected.Type = "int64"
	funcQuery.Ret.Add(pt_query__ret__row_affected)
}

func (t *GenCode) gen_query__add__func__body__arg(_arrs_arg []string) (s_arg string) {
	for _, s_arg__one := range _arrs_arg {
		s_arg += fmt.Sprintf("\n\t%s,", s_arg__one)
	}
	return s_arg
}

func (t *GenCode) gen_query__add__func__body__set_args(_arrs_arg []string) (s_gen_arg string) {
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

func (t *GenCode) gen_query__add__func__body__insert_multi__query(_s_query_values string) (s_multi_insert__query string) {
	s_multi_insert__query += fmt.Sprintf("\n\tstrings.Repeat(\", (%s)\", n_len_arg-1),", _s_query_values)
	return s_multi_insert__query
}

func (t *GenCode) gen_query__add__func__body__insert_multi__proc(_arrs_arg []string) (s_multi_insert__body string) {
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

func (t *GenCode) gen_query__add__func__body__set_args__remove_sets(_arrs_arg []string) (s_sql_update__delete_sets string) {
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
