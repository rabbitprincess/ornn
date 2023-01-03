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

func (t *GenCode) code(config *config.Config, genQueries *GenQueries) (genCode string, err error) {
	t.config = config
	t.codeGen = &codegen.CodeGen{}
	t.codeGen.DoNotEdit = t.config.Global.DoNotEdit
	t.codeGen.Package = t.config.Global.PackageName
	for _, imp := range config.Global.Import {
		t.codeGen.AddImport(&codegen.Import{
			Path:  imp.Path,
			Alias: imp.Alias,
		})
	}

	// root struct
	rootStruct := &codegen.Struct{
		Name: t.config.Global.ClassName,
	}
	t.codeGen.AddItem(rootStruct)

	// init function
	rootFunc := &codegen.Function{
		StructName: "t",
		StructType: "*" + rootStruct.Name,
		FuncName:   "Init",
	}
	t.codeGen.AddItem(rootFunc)

	// arg
	rootFuncInitArg := &codegen.Var{
		Name: "job",
		Type: "*Job",
	}
	rootFunc.AddArg(rootFuncInitArg)

	for groupName, queryGroup := range genQueries.groups {
		genClass := t.genClass(groupName)
		t.codeGen.AddItem(genClass)

		// root 구조체 안에 queries 구조체 포인터 선언
		rootStruct.AddField(&codegen.Var{
			Type: genClass.Name,
			Name: genClass.Name,
		})
		rootFunc.InlineCode += fmt.Sprintf("%s.%s.%s(%s)\n", "t", genClass.Name, "Init", rootFunc.Args.Items[0].Name)

		for _, query := range queryGroup {
			genFunc := t.genFunc(genClass.Name, query)
			t.codeGen.AddItem(genFunc)
		}
	}

	// 소스 출력
	genCode = t.codeGen.Code()
	return genCode, nil
}

func (t *GenCode) genClass(name string) (genGroup *codegen.Struct) {
	genGroup = &codegen.Struct{
		Name: sql.Util_ConvFirstToUpper(name),
	}

	// root 구조체 연결을 위한 구조체 필드 변수 제작
	groupVar := &codegen.Var{
		Name: "job",
		Type: "*Job",
	}
	genGroup.AddField(groupVar)

	// root 구조체에서 초기화를 요청할 Init 함수 제작
	groupFuncInit := &codegen.Function{
		FuncName:   "Init",
		StructName: "t",
		StructType: fmt.Sprintf("*%s", genGroup.Name),
	}
	t.codeGen.AddItem(groupFuncInit)

	// args
	groupFuncInitArg := &codegen.Var{
		Name: "job",
		Type: "*Job",
	}
	groupFuncInit.AddArg(groupFuncInitArg)

	// inline code
	groupFuncInit.InlineCode = fmt.Sprintf("%s.%s = %s", groupFuncInit.StructName, groupVar.Name, groupFuncInitArg.Name)

	return genGroup
}

func (t *GenCode) genFunc(groupName string, query *GenQuery) (funcQuery *codegen.Function) {
	funcQuery = &codegen.Function{
		StructName: "t",
		StructType: "*" + groupName,
		FuncName:   sql.Util_ConvFirstToUpper(query.queryName),
	}

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
	return funcQuery
}

func (t *GenCode) genQuerySelect(funcQuery *codegen.Function, query *GenQuery) {
	// struct for select
	structName := t.genQuery_struct_select(funcQuery, query)

	// args
	tpls := t.genQuery_tpls(funcQuery, query)
	args := t.genQuery_args(funcQuery, query)

	// rets
	retItemName, retItemType := t.genQuery_ret_select(funcQuery, structName, query.SelectSingle)
	t.genQuery_ret_error(funcQuery)

	// body
	funcQuery.InlineCode = template.Select(args, tpls, query.query, query.SelectSingle, "t", "job", structName, retItemName, retItemType)
}

func (t *GenCode) genQueryInsert(funcQuery *codegen.Function, query *GenQuery) {
	// args
	args := t.genQuery_args(funcQuery, query)
	tpls := t.genQuery_tpls(funcQuery, query)

	// rets
	t.genQuery_ret_lastInsertId(funcQuery)
	t.genQuery_ret_error(funcQuery)

	// body
	funcQuery.InlineCode = template.Insert(args, tpls, query.query, query.InsertMulti, "t", "job")
}

func (t *GenCode) genQueryUpdate(funcQuery *codegen.Function, query *GenQuery) {
	// args
	args := t.genQuery_args(funcQuery, query)
	tpls := t.genQuery_tpls(funcQuery, query)

	// rets
	t.genQuery_ret_rowAffected(funcQuery)
	t.genQuery_ret_error(funcQuery)

	// body
	funcQuery.InlineCode = template.Update(args, tpls, query.query, query.UpdateNullIgnore, "t", "job")
}

func (t *GenCode) genQueryDelete(funcQuery *codegen.Function, query *GenQuery) {
	// args
	tpls := t.genQuery_tpls(funcQuery, query)
	args := t.genQuery_args(funcQuery, query)

	// rets
	t.genQuery_ret_rowAffected(funcQuery)
	t.genQuery_ret_error(funcQuery)

	// body
	funcQuery.InlineCode = template.Delete(args, query.query, tpls, "t", "job")
}

func (t *GenCode) genQuery_tpls(funcQuery *codegen.Function, query *GenQuery) (tpls []string) {
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

func (t *GenCode) genQuery_args(funcQuery *codegen.Function, query *GenQuery) (args []string) {
	args = make([]string, 0, len(query.tpl.pairs))

	for _, pair := range query.arg.pairs {
		arg := &codegen.Var{}
		arg.Name = fmt.Sprintf("arg_%s", pair.Key)
		if pair.Value != "" { // 형을 특정할 수 있을 때
			arg.Type = pair.Value
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

func (t *GenCode) genQuery_struct_select(funcQuery *codegen.Function, query *GenQuery) (retStructName string) {
	retStruct := &codegen.Struct{
		Name: fmt.Sprintf("%s_%s", sql.Util_ConvFirstToUpper(query.groupName), strings.ToLower(funcQuery.FuncName)),
	}
	for _, pair := range query.ret.pairs {
		item := &codegen.Var{
			Name: sql.Util_ConvFirstToUpper(pair.Key),
			Type: pair.Value,
		}
		retStruct.AddField(item)
	}
	t.codeGen.AddItem(retStruct)
	return retStruct.Name
}

func (t *GenCode) genQuery_ret_select(funcQuery *codegen.Function, retStructName string, selectSingle bool) (retItemName, retItemType string) {
	retItem := &codegen.Var{
		Name: strings.ToLower(funcQuery.FuncName),
		Type: "*" + retStructName,
	}
	if selectSingle != true {
		retItem.Name = retItem.Name + "s"
		retItem.Type = "[]" + retItem.Type
	}
	funcQuery.AddRet(retItem)
	return retItem.Name, retItem.Type
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
