package ornn

import (
	"fmt"
	"log"

	"strings"

	"github.com/gokch/ornn/codegen"
	"github.com/gokch/ornn/config"
	"github.com/gokch/ornn/sql"
)

const (
	// TODO
	DEF_s_gen_config__go__db__func__body__func_name__SQL_remove__update_set_field__null string = "SQL_remove__update_set_field__null"
)

type GenCode struct {
	config  *config.Config
	codeGen *codegen.CodeGen
}

func (t *GenCode) gen(config *config.Config, genData *GenData) (genCode string, err error) {
	// config 설정
	t.config = config
	t.codeGen = &codegen.CodeGen{}

	// gen_go 에 소스 생성을 위한 데이터 넣기
	{
		t.codeGen.DoNotEdit = t.config.Global.DoNotEdit
		t.codeGen.Package = t.config.Global.PackageName

		// import 경로 추가
		for _, imp := range config.Global.Import {
			t.codeGen.AddImport(&codegen.Import{
				Path:  imp.Path,
				Alias: imp.Alias,
			})
		}

		// 루트 구조체 작성
		rootStruct := &codegen.Struct{}
		t.codeGen.AddItem(rootStruct)
		rootStruct.Name = t.config.Global.ClassName

		// 루트 함수 작성
		rootFunc := &codegen.Function{}
		t.codeGen.AddItem(rootFunc)
		rootFunc.StructName = t.config.Global.StructName
		rootFunc.StructType = fmt.Sprintf("*%s", rootStruct.Name)
		rootFunc.FuncName = "Init"

		// arg
		rootFuncInitArg := &codegen.Var{}
		rootFunc.AddArg(rootFuncInitArg)
		rootFuncInitArg.Name = strings.ToLower(t.config.Global.InstanceName)
		rootFuncInitArg.Type = t.config.Global.InstanceType

		// group 단위 구조체
		for _, gen_group := range genData.groups {
			// group 구조체 생성
			group := t.genGroup(gen_group.Name)
			t.codeGen.AddItem(group)

			// root 구조체 안에 필드 변수 선언 -> group 구조체 사용을 위해
			{
				// root 구조체 안에 group 구조체 포인터 선언
				rootVars := &codegen.Var{}
				rootVars.Type = group.Name
				rootVars.Name = strings.ToLower(gen_group.Name)
				rootStruct.AddField(rootVars)

				// root init body 작성
				s_code := fmt.Sprintf("%s.%s.%s(%s)\n", t.config.Global.StructName, rootVars.Name, "Init", rootFunc.Arg.Items[0].Name)
				rootFunc.InlineCode += s_code
			}

			// group 구조체 안에 query 함수 생성
			for _, pt_gen_query := range gen_group.Queries {
				t.genQuery(group, pt_gen_query)
			}
		} // end of for pt_group
	}

	// 실제 소스 출력
	genCode = t.codeGen.Code()

	return genCode, nil
}

func (t *GenCode) genGroup(group string) (genGroup *codegen.Struct) {
	genGroup = &codegen.Struct{}
	genGroup.Name = sql.Util_ConvFirstToUpper(group)

	// group 구조체 안에
	{
		// root 구조체 연결을 위한 구조체 필드 변수 제작
		groupVar := &codegen.Var{}
		genGroup.AddField(groupVar)
		{
			groupVar.Name = t.config.Global.InstanceName
			groupVar.Type = t.config.Global.InstanceType
		}

		// root 구조체에서 초기화를 요청할 Init 함수 제작
		groupFuncInit := &codegen.Function{}
		t.codeGen.AddItem(groupFuncInit)
		{
			groupFuncInit.FuncName = "Init"
			groupFuncInit.StructName = t.config.Global.StructName
			groupFuncInit.StructType = fmt.Sprintf("*%s", genGroup.Name)

			// args
			groupFuncInitArg := &codegen.Var{}
			groupFuncInit.AddArg(groupFuncInitArg)
			groupFuncInitArg.Name = strings.ToLower(t.config.Global.InstanceName)
			groupFuncInitArg.Type = t.config.Global.InstanceType

			// inline code
			groupFuncInit.InlineCode = fmt.Sprintf("%s.%s = %s", groupFuncInit.StructName, groupVar.Name, groupFuncInitArg.Name)
		}
	}
	return genGroup
}

func (t *GenCode) genQuery(structGroup *codegen.Struct, query *GenDataQuery) {
	funcQuery := &codegen.Function{}

	funcQuery.StructName = t.config.Global.StructName
	funcQuery.StructType = fmt.Sprintf("*%s", structGroup.Name)
	funcQuery.FuncName = sql.Util_ConvFirstToUpper(query.queryName)

	switch query.queryType {
	case QueryTypeSelect:
		t.genQuerySelect(funcQuery, query)
	case QueryTypeInsert:
		t.genQueryInsert(funcQuery, query)
	case QueryTypeUpdate:
		t.genQueryUpdate(funcQuery, query)
	case QueryTypeDelete:
		t.genQueryDelete(funcQuery, query)
	default:
		log.Fatalf("invalid query type | query type : %v", query.queryType)
	}

	t.codeGen.AddItem(funcQuery)
}

func (t *GenCode) genQuerySelect(
	funcQuery *codegen.Function,
	query *GenDataQuery,
) {
	// 1. 함수 입력 인자
	tpl := t.gen_query__add__func__arg__tpl(funcQuery, query)
	arg := t.gen_query__add__func__arg(funcQuery, query)

	// 2. 함수 리턴
	ret := &codegen.Struct{}
	t.codeGen.AddItem(ret)

	retItem := &codegen.Var{}
	var s_body_code__ret_declare string
	var s_body_code__ret_set string
	{
		// 2-1. 쿼리-리턴 변수 선언
		{
			ret.Name = fmt.Sprintf("%s_%s", sql.Util_ConvFirstToUpper(query.tableName), strings.ToLower(funcQuery.FuncName))
			for _, pt_field_type := range query.ret.pairs {
				item := &codegen.Var{}
				item.Name = sql.Util_ConvFirstToUpper(pt_field_type.Key)
				item.Type = pt_field_type.Value
				ret.AddField(item)
			}
		}

		// 2-2. 리턴 변수 처리
		{
			// 리턴 변수 선언 - 구조체
			if query.isSelectSingle == true {
				retItem.Name = fmt.Sprintf("pt_%s", strings.ToLower(funcQuery.FuncName))
				retItem.Type = fmt.Sprintf("*%s", ret.Name)
				s_body_code__ret_set = fmt.Sprintf("%s = pt_struct\n\tbreak", retItem.Name)
			} else {
				retItem.Name = fmt.Sprintf("%ss", strings.ToLower(funcQuery.FuncName))
				retItem.Type = fmt.Sprintf("[]*%s", ret.Name)
				s_body_code__ret_declare = fmt.Sprintf("\n%s = make(%s, 0, 100)", retItem.Name, retItem.Type)
				s_body_code__ret_set = fmt.Sprintf("%s = append(%s, pt_struct)", retItem.Name, retItem.Name)
			}
			funcQuery.AddRet(retItem)

			// error 추가
			t.gen_query__add__func__ret__error(funcQuery)
		}
	}

	// 3. 함수 body
	{
		funcQuery.InlineCode = fmt.Sprintf(`
%s
sql := fmt.Sprintf(
	"%s",%s
)
ret, err := %s.%s.Query(
	sql,
	args...,
)
if err != nil {
	return nil, err
}
defer ret.Close()
%s
for ret.Next() {
	pt_struct := &%s{}
	err := ret.Scan(pt_struct)
	if err != nil {
		return nil, err
	}
	%s
}

return %s, nil
`,
			t.gen_query__add__func__body__set_args(arg),
			query.query,
			t.gen_query__add__func__body__arg(tpl),
			t.config.Global.StructName,
			t.config.Global.InstanceName,
			s_body_code__ret_declare,
			ret.Name,
			s_body_code__ret_set,
			retItem.Name,
		)
	}
	return
}

func (t *GenCode) genQueryInsert(funcQuery *codegen.Function, query *GenDataQuery) {
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
		if query.InsertMulti == true { // multi insert
			s_query_values := sql.Util_ExportInsertQueryValues(query.query)
			if query.query[len(query.query)-1:] == ";" {
				query.query = query.query[:len(query.query)-1]
			}
			query.query += "%s"
			genArgs = t.gen_query__add__func__body__insert_multi__proc(arg)
			multiInsert = t.gen_query__add__func__body__insert_multi__query(s_query_values)
		} else { // insert
			genArgs = t.gen_query__add__func__body__set_args(arg)
		}

		funcQuery.InlineCode = fmt.Sprintf(`
%s
sql := fmt.Sprintf(
	"%s",%s%s
)

exec, err := %s.%s.Exec(
	sql,
	args...,
)
if err != nil {
	return 0, err
}

return exec.LastInsertId()
`,
			genArgs,
			query.query,
			t.gen_query__add__func__body__arg(tpl),
			multiInsert,
			t.config.Global.StructName,
			t.config.Global.InstanceName,
		)
	}
	return
}

func (t *GenCode) genQueryUpdate(funcQuery *codegen.Function, query *GenDataQuery) {
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
		if query.UpdateNullIgnore == true {
			s_body__set_args = t.gen_query__add__func__body__set_args__remove_sets(arg)
		} else {
			s_body__set_args = t.gen_query__add__func__body__set_args(arg)
		}
		funcQuery.InlineCode = fmt.Sprintf(`
sql := fmt.Sprintf(
	"%s",%s
)
%s
exec, err := %s.%s.Exec(
	sql,
	args...,
)
if err != nil {
	return 0, err
}

return exec.RowsAffected()
`,
			query.query,
			t.gen_query__add__func__body__arg(tpl),
			s_body__set_args,
			t.config.Global.StructName,
			t.config.Global.InstanceName,
		)
	}
	return
}

func (t *GenCode) genQueryDelete(funcQuery *codegen.Function, query *GenDataQuery) {
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
sql := fmt.Sprintf(
	"%s",%s
)
		
exec, err := %s.%s.Exec(
	sql,
	args...,
)
if err != nil {
	return 0, err
}

return exec.RowsAffected()
`,
			t.gen_query__add__func__body__set_args(arrs_arg),
			query.query,
			t.gen_query__add__func__body__arg(arrs_tpl),
			t.config.Global.StructName,
			t.config.Global.InstanceName,
		)
	}
	return
}

func (t *GenCode) gen_query__add__func__arg__tpl(funcQuery *codegen.Function, query *GenDataQuery) (arrs_tpl []string) {
	arrs_tpl = make([]string, 0, len(query.tpl.pairs))
	for _, pt_tpl := range query.tpl.pairs {
		arg := &codegen.Var{}
		funcQuery.AddArg(arg)
		arg.Name = fmt.Sprintf("%s%s", t.config.Global.TplPrefix, pt_tpl.Key)
		arg.Type = "string"

		arrs_tpl = append(arrs_tpl, arg.Name)
	}
	return
}

func (t *GenCode) gen_query__add__func__arg(funcQuery *codegen.Function, query *GenDataQuery) (arrs_arg []string) {
	arrs_arg = make([]string, 0, len(query.tpl.pairs))

	for _, pt_field_type := range query.arg.pairs {
		arg := &codegen.Var{}
		var varType string
		// 1. type 판정
		{
			if pt_field_type.Value == "" { // 형을 특정할 수 없을 때
				varType = "interface{}"
			} else { // 형을 특정할 수 있을 때
				varType = "*" + pt_field_type.Value
			}
			if query.InsertMulti == true {
				varType = "[]" + varType
			}
		}

		// 2. set query arg
		{
			arg.Name = fmt.Sprintf("%s%s", t.config.Global.ArgPrefix, pt_field_type.Key)
			arg.Type = varType
			funcQuery.AddArg(arg)
		}
		arrs_arg = append(arrs_arg, arg.Name)
	}
	return
}

func (t *GenCode) gen_query__add__func__ret__error(funcQuery *codegen.Function) {
	retErr := &codegen.Var{}
	retErr.Name = "err"
	retErr.Type = "error"
	funcQuery.AddRet(retErr)
}

func (t *GenCode) gen_query__add__func__ret__last_insert_id(funcQuery *codegen.Function) {
	retLastid := &codegen.Var{}
	retLastid.Name = "lastInsertId"
	retLastid.Type = "int64"
	funcQuery.AddRet(retLastid)
}

func (t *GenCode) gen_query__add__func__ret__row_affected(funcQuery *codegen.Function) {
	pt_query__ret__row_affected := &codegen.Var{}
	pt_query__ret__row_affected.Name = "rowAffected"
	pt_query__ret__row_affected.Type = "int64"
	funcQuery.AddRet(pt_query__ret__row_affected)
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

	s_gen_arg += fmt.Sprintf(`args := make([]interface{}, 0, %d)
args = append(args, I_to_arri(%s)...)
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

args := make([]interface{}, 0, n_len_arg*%d)
for i := 0; i < n_len_arg; i++ {
	args = append(args, I_to_arri(
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

func (t *GenCode) gen_query__add__func__body__set_args__remove_sets(_arrs_arg []string) (sql_update__delete_sets string) {
	var sql_is_nil string
	for _, s_arg := range _arrs_arg {
		s_field_name := strings.TrimPrefix(s_arg, t.config.Global.ArgPrefix)

		sql_is_nil += fmt.Sprintf(`if %s == nil {
	arrs_sets__removed = append(arrs_sets__removed, "%s")
} else {
	args = append(args, %s)
}
`,
			s_arg,
			s_field_name,
			s_arg,
		)
	}

	sql_update__delete_sets = fmt.Sprintf(`
args := make([]interface{}, 0, %d)
arrs_sets__removed := make([]string, 0, %d)
%s
if len(arrs_sets__removed) != 0 {
	sql, _ = %s(sql, arrs_sets__removed)
}
`,
		len(_arrs_arg),
		len(_arrs_arg),
		sql_is_nil,
		DEF_s_gen_config__go__db__func__body__func_name__SQL_remove__update_set_field__null,
	)
	return sql_update__delete_sets
}
