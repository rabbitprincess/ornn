package ornn

import (
	"fmt"
	"log"

	"strings"

	"github.com/gokch/ornn/codegen"
	"github.com/gokch/ornn/config"
	"github.com/gokch/ornn/sql"
	"github.com/gokch/ornn/sql/template"
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
		log.Fatalf("need more programming | invalid query type | query type : %v", query.queryType)
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
	funcQuery.AddRet(retItem)

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
		retItem.Name = fmt.Sprintf("%s", strings.ToLower(funcQuery.FuncName))
		retItem.Type = fmt.Sprintf("*%s", ret.Name)
	} else {
		retItem.Name = fmt.Sprintf("%ss", strings.ToLower(funcQuery.FuncName))
		retItem.Type = fmt.Sprintf("[]*%s", ret.Name)
	}
	// error 추가
	t.genQuery_ret_error(funcQuery)

	// 3. 함수 body
	funcQuery.InlineCode = template.Select(args, tpls, query.query, query.SelectSingle, "t", "Job", ret.Name, retItem.Name, retItem.Type)
}

func (t *GenCode) genQueryInsert(funcQuery *codegen.Function, query *GenDataQuery) {
	// 1. 함수 입력 인자
	args := t.genQuery_args(funcQuery, query)
	tpls := t.genQuery_tpls(funcQuery, query)

	// 2. 함수 리턴 변수
	t.genQuery_ret_lastInsertId(funcQuery)
	t.genQuery_ret_error(funcQuery)

	// 3. 함수 body
	funcQuery.InlineCode = template.Insert(args, tpls, query.query, query.InsertMulti, "t", "Job")
}

func (t *GenCode) genQueryUpdate(funcQuery *codegen.Function, query *GenDataQuery) {
	// 1. 함수 입력 인자
	args := t.genQuery_args(funcQuery, query)
	tpls := t.genQuery_tpls(funcQuery, query)

	// 2. 함수 리턴 변수
	t.genQuery_ret_rowAffected(funcQuery)
	t.genQuery_ret_error(funcQuery)

	// 3. 함수 body
	funcQuery.InlineCode = template.Update(args, tpls, query.query, query.UpdateNullIgnore, "t", "Job")
}

func (t *GenCode) genQueryDelete(funcQuery *codegen.Function, query *GenDataQuery) {
	// 1. 함수 입력 인자
	tpls := t.genQuery_tpls(funcQuery, query)
	args := t.genQuery_args(funcQuery, query)

	// 2. 함수 리턴 변수
	t.genQuery_ret_rowAffected(funcQuery)
	t.genQuery_ret_error(funcQuery)

	// 3. 함수 body
	funcQuery.InlineCode = template.Delete(args, query.query, tpls, "t", "Job")
}

func (t *GenCode) genQuery_tpls(funcQuery *codegen.Function, query *GenDataQuery) (tpls []string) {
	tpls = make([]string, 0, len(query.tpl.pairs))
	for _, tpl := range query.tpl.pairs {
		arg := &codegen.Var{
			Name: fmt.Sprintf("tpl_%s", tpl.Key),
			Type: "string",
		}
		funcQuery.AddArg(arg)
		tpls = append(tpls, arg.Name)
	}
	return tpls
}

func (t *GenCode) genQuery_args(funcQuery *codegen.Function, query *GenDataQuery) (args []string) {
	args = make([]string, 0, len(query.tpl.pairs))

	for _, pair := range query.arg.pairs {
		arg := &codegen.Var{}
		arg.Name = fmt.Sprintf("arg_%s", pair.Key)
		if pair.Value != "" { // 형을 특정할 수 있을 때
			arg.Type = "*" + pair.Value
		} else { // 형을 특정할 수 없을 때
			arg.Type = "interface{}"
		}
		if query.InsertMulti == true {
			arg.Type = "[]" + arg.Type
		}
		funcQuery.AddArg(arg)
		args = append(args, arg.Name)
	}
	return args
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
