package ornn

import (
	"fmt"
	"log"

	"strings"

	"github.com/gokch/ornn/codegen"
	"github.com/gokch/ornn/config"
	"github.com/gokch/ornn/sql"
)

type GenCode struct {
	config  *config.Config
	codeGen *codegen.CodeGen
}

func (t *GenCode) code(config *config.Config, genData *GenData) (genCode string, err error) {
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
		rootFunc.StructName = "t"
		rootFunc.StructType = fmt.Sprintf("*%s", rootStruct.Name)
		rootFunc.FuncName = "Init"

		// arg
		rootFuncInitArg := &codegen.Var{}
		rootFunc.AddArg(rootFuncInitArg)
		rootFuncInitArg.Name = strings.ToLower("Job")
		rootFuncInitArg.Type = "*Job"

		// group 단위 구조체
		for _, genGroup := range genData.groups {
			// group 구조체 생성
			group := t.genGroup(genGroup.Name)
			t.codeGen.AddItem(group)

			// root 구조체 안에 필드 변수 선언 -> group 구조체 사용을 위해
			{
				// root 구조체 안에 group 구조체 포인터 선언
				rootVars := &codegen.Var{}
				rootVars.Type = group.Name
				rootVars.Name = strings.ToLower(genGroup.Name)
				rootStruct.AddField(rootVars)

				// root init body 작성
				code := fmt.Sprintf("%s.%s.%s(%s)\n", "t", rootVars.Name, "Init", rootFunc.Arg.Items[0].Name)
				rootFunc.InlineCode += code
			}

			// group 구조체 안에 query 함수 생성
			for _, query := range genGroup.Queries {
				t.genQuery(group, query)
			}
		} // end of for group
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
			groupVar.Name = "Job"
			groupVar.Type = "*Job"
		}

		// root 구조체에서 초기화를 요청할 Init 함수 제작
		groupFuncInit := &codegen.Function{}
		t.codeGen.AddItem(groupFuncInit)
		{
			groupFuncInit.FuncName = "Init"
			groupFuncInit.StructName = "t"
			groupFuncInit.StructType = fmt.Sprintf("*%s", genGroup.Name)

			// args
			groupFuncInitArg := &codegen.Var{}
			groupFuncInit.AddArg(groupFuncInitArg)
			groupFuncInitArg.Name = strings.ToLower("Job")
			groupFuncInitArg.Type = "*Job"

			// inline code
			groupFuncInit.InlineCode = fmt.Sprintf("%s.%s = %s", groupFuncInit.StructName, groupVar.Name, groupFuncInitArg.Name)
		}
	}
	return genGroup
}

func (t *GenCode) genQuery(structGroup *codegen.Struct, query *GenDataQuery) {
	funcQuery := &codegen.Function{}

	funcQuery.StructName = "t"
	funcQuery.StructType = "*" + structGroup.Name
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

func (t *GenCode) genQuerySelect(funcQuery *codegen.Function, query *GenDataQuery) {
	// 1. 함수 입력 인자
	tpls := t.genQuery_tpls(funcQuery, query)
	args := t.genQuery_args(funcQuery, query)

	// 2. 함수 리턴
	ret := &codegen.Struct{}
	t.codeGen.AddItem(ret)

	retItem := &codegen.Var{}
	var bodyRetDeclare string
	var bodyRetSet string
	{
		// 2-1. 리턴 변수 선언
		ret.Name = fmt.Sprintf("%s_%s", sql.Util_ConvFirstToUpper(query.groupName), strings.ToLower(funcQuery.FuncName))
		for _, pair := range query.ret.pairs {
			item := &codegen.Var{}
			item.Name = sql.Util_ConvFirstToUpper(pair.Key)
			item.Type = pair.Value
			ret.AddField(item)
		}

		// 2-2. 리턴 변수 처리
		// 리턴 변수 선언 - 구조체
		if query.SelectSingle == true {
			retItem.Name = fmt.Sprintf("pt_%s", strings.ToLower(funcQuery.FuncName))
			retItem.Type = fmt.Sprintf("*%s", ret.Name)
			bodyRetSet = fmt.Sprintf("%s = scan\n\tbreak", retItem.Name)
		} else {
			retItem.Name = fmt.Sprintf("%ss", strings.ToLower(funcQuery.FuncName))
			retItem.Type = fmt.Sprintf("[]*%s", ret.Name)
			bodyRetDeclare = fmt.Sprintf("\n%s = make(%s, 0, 100)", retItem.Name, retItem.Type)
			bodyRetSet = fmt.Sprintf("%s = append(%s, scan)", retItem.Name, retItem.Name)
		}
		funcQuery.AddRet(retItem)

		// error 추가
		t.genQuery_ret_error(funcQuery)
	}

	// 3. 함수 body
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
	scan := &%s{}
	err := ret.Scan(scan)
	if err != nil {
		return nil, err
	}
	%s
}

return %s, nil
`,
		t.genQuery_body_setArgs(args),
		query.query,
		t.genQuery_body_arg(tpls),
		"t",
		"Job",
		bodyRetDeclare,
		ret.Name,
		bodyRetSet,
		retItem.Name,
	)
	return
}

func (t *GenCode) genQueryInsert(funcQuery *codegen.Function, query *GenDataQuery) {
	// 1. 함수 입력 인자
	tpl := t.genQuery_tpls(funcQuery, query)
	arg := t.genQuery_args(funcQuery, query)

	// 2. 함수 리턴 변수
	t.genQuery_ret_lastInsertId(funcQuery)
	t.genQuery_ret_error(funcQuery)

	// 3. 함수 body
	var multiInsert, genArgs string
	if query.InsertMulti == true { // multi insert
		s_query_values := sql.Util_ExportInsertQueryValues(query.query)
		if query.query[len(query.query)-1:] == ";" {
			query.query = query.query[:len(query.query)-1]
		}
		query.query += "%s"
		genArgs = t.genQuery_body_multiInsertProc(arg)
		multiInsert = t.genQuery_body_multiInsert(s_query_values)
	} else { // insert
		genArgs = t.genQuery_body_setArgs(arg)
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
		t.genQuery_body_arg(tpl),
		multiInsert,
		"t",
		"Job",
	)
	return
}

func (t *GenCode) genQueryUpdate(funcQuery *codegen.Function, query *GenDataQuery) {
	// 1. 함수 입력 인자
	tpl := t.genQuery_tpls(funcQuery, query)
	arg := t.genQuery_args(funcQuery, query)

	// 2. 함수 리턴 변수
	t.genQuery_ret_rowAffected(funcQuery)
	t.genQuery_ret_error(funcQuery)

	// 3. 함수 body
	var body string
	if query.UpdateNullIgnore == true {
		body = t.genQuery_body_removeSets(arg)
	} else {
		body = t.genQuery_body_setArgs(arg)
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
		t.genQuery_body_arg(tpl),
		body,
		"t",
		"Job",
	)
}

func (t *GenCode) genQueryDelete(funcQuery *codegen.Function, query *GenDataQuery) {
	// 1. 함수 입력 인자
	tpls := t.genQuery_tpls(funcQuery, query)
	args := t.genQuery_args(funcQuery, query)

	// 2. 함수 리턴 변수
	t.genQuery_ret_rowAffected(funcQuery)
	t.genQuery_ret_error(funcQuery)

	// 3. 함수 body
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
		t.genQuery_body_setArgs(args),
		query.query,
		t.genQuery_body_arg(tpls),
		"t",
		"Job",
	)
}

func (t *GenCode) genQuery_tpls(funcQuery *codegen.Function, query *GenDataQuery) (tpls []string) {
	tpls = make([]string, 0, len(query.tpl.pairs))
	for _, tpl := range query.tpl.pairs {
		arg := &codegen.Var{}
		funcQuery.AddArg(arg)
		arg.Name = fmt.Sprintf("%s%s", t.config.Global.TplPrefix, tpl.Key)
		arg.Type = "string"

		tpls = append(tpls, arg.Name)
	}
	return
}

func (t *GenCode) genQuery_args(funcQuery *codegen.Function, query *GenDataQuery) (args []string) {
	args = make([]string, 0, len(query.tpl.pairs))

	for _, pair := range query.arg.pairs {
		arg := &codegen.Var{}
		arg.Name = fmt.Sprintf("%s%s", t.config.Global.ArgPrefix, pair.Key)
		if pair.Value != "" { // 형을 특정할 수 없을 때
			arg.Type = "*" + pair.Value
		} else { // 형을 특정할 수 있을 때
			arg.Type = "interface{}"
		}
		if query.InsertMulti == true {
			arg.Type = "[]" + arg.Type
		}
		funcQuery.AddArg(arg)
		args = append(args, arg.Name)
	}
	return
}

func (t *GenCode) genQuery_ret_error(funcQuery *codegen.Function) {
	funcQuery.AddRet(&codegen.Var{
		Name: "err",
		Type: "error",
	})
}

func (t *GenCode) genQuery_ret_lastInsertId(funcQuery *codegen.Function) {
	funcQuery.AddRet(&codegen.Var{
		Name: "lastInsertId",
		Type: "int64",
	})
}

func (t *GenCode) genQuery_ret_rowAffected(funcQuery *codegen.Function) {
	funcQuery.AddRet(&codegen.Var{
		Name: "rowAffected",
		Type: "int64",
	})
}

func (t *GenCode) genQuery_body_arg(args []string) (ret string) {
	for _, arg := range args {
		ret += fmt.Sprintf("\n\t%s,", arg)
	}
	return ret
}

func (t *GenCode) genQuery_body_setArgs(args []string) (arg string) {
	genArgItem := t.genQuery_body_arg(args)
	if genArgItem != "" {
		genArgItem += "\n"
	}

	arg += fmt.Sprintf(`args := make([]interface{}, 0, %d)
args = append(args, I_to_arri(%s)...)
`,
		len(args),
		genArgItem,
	)

	return arg
}

func (t *GenCode) genQuery_body_multiInsert(query string) (multiInsert string) {
	return fmt.Sprintf("\n\tstrings.Repeat(\", (%s)\", argLen-1),", query)
}

func (t *GenCode) genQuery_body_multiInsertProc(args []string) (multiInsertProc string) {
	var checkLen string
	for i, arg := range args {
		checkLen += fmt.Sprintf("argLen != len(%s)", arg)
		if i != len(args)-1 {
			checkLen += fmt.Sprintf(" || ")
		}
	}

	var append string
	for i, arg := range args {
		append += fmt.Sprintf("%s[i]", arg)
		if i != len(args)-1 {
			append += fmt.Sprintf(",\n\t\t")
		}
	}

	multiInsertProc = fmt.Sprintf(`argLen := len(%s)
if argLen == 0 {
	return 0, fmt.Errorf("arg len is zero")
}
if %s {
	return 0, fmt.Errorf("arg len is not same")
}

args := make([]interface{}, 0, argLen*%d)
for i := 0; i < argLen; i++ {
	args = append(args, I_to_arri(
		%s,
	)...)
}
`,
		args[0],
		checkLen,
		len(args),
		append)
	return multiInsertProc
}

func (t *GenCode) genQuery_body_removeSets(args []string) (removeSets string) {
	var isNil string
	for _, arg := range args {
		fieldName := strings.TrimPrefix(arg, t.config.Global.ArgPrefix)

		isNil += fmt.Sprintf(`if %s == nil {
	setsRemoved = append(setsRemoved, "%s")
} else {
	args = append(args, %s)
}
`,
			arg,
			fieldName,
			arg,
		)
	}

	removeSets = fmt.Sprintf(`
args := make([]interface{}, 0, %d)
setsRemoved := make([]string, 0, %d)
%s
if len(setsRemoved) != 0 {
	sql, _ = RemoveNull(sql, setsRemoved)
}
`,
		len(args),
		len(args),
		isNil,
	)
	return removeSets
}
