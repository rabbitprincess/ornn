package ornn

import (
	"fmt"

	"github.com/gokch/ornn/config"
	"github.com/gokch/ornn/parser"
)

type GenQueries struct {
	conf *config.Config
	psr  parser.Parser

	class map[string]map[string]*parser.ParseQuery
}

func (t *GenQueries) Init(conf *config.Config, psr parser.Parser) {
	t.conf = conf
	t.psr = psr
	t.class = make(map[string]map[string]*parser.ParseQuery)
}

func (t *GenQueries) SetData() (err error) {
	// schema
	for _, group := range t.conf.Schema.Tables {
		Queries, ok := t.conf.Queries.Tables[group.Name]
		if ok != true {
			continue
		}
		err := t.SetDataGroup(group.Name, Queries)
		if err != nil {
			return err
		}
	}

	// custom
	for groupName, custom := range t.conf.Queries.Custom {
		err := t.SetDataGroup(groupName, custom)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *GenQueries) SetDataGroup(groupName string, queries []*config.Query) (err error) {
	if t.class == nil {
		t.class = make(map[string]map[string]*parser.ParseQuery)
	} else if t.class[groupName] == nil {
		t.class[groupName] = make(map[string]*parser.ParseQuery)
	}

	for _, query := range queries {
		parseQuery, err := t.SetDataQuery(groupName, query)
		if err != nil {
			return err
		}
		t.class[groupName][query.Name] = parseQuery
	}
	return nil
}

func (t *GenQueries) SetDataQuery(groupName string, query *config.Query) (parseQuery *parser.ParseQuery, err error) {
	parseQuery, err = t.psr.Parse(query.Sql)
	if err != nil {
		query.ErrParser = fmt.Sprintf("%v", err)
		return nil, nil
	}
	return parseQuery, nil

	/*
		parseQuery = &parser.ParseQuery{}
		parseQuery.Init()

		// set args
		// tpl args ( # name # )를 배열로 추출
		tpls, err := sql.Util_ExportBetweenDelimiter(query.Sql, sql.TplDelimiter)
		if err != nil {
			return nil, err
		}

		for _, tpl := range tpls {
			tmps := strings.Split(tpl, sql.TplSplit)
			var argName string
			var argData string
			if len(tmps) == 1 {
				argName = tmps[0]
				argData = ""
			} else if len(tmps) == 2 {
				argName = tmps[0]
				argData = tmps[1]
			} else {
				return nil, fmt.Errorf("tpl format is wrong - %s", tpl)
			}

			parseQuery.Tpl[argName] = argData
		}

		// args ( % name % )를 배열로 추출
		args, err := sql.Util_ExportBetweenDelimiter(query.Sql, sql.PrepareStatementDelimeter)
		if err != nil {
			return nil, err
		}
		for _, arg := range args {
			parseQuery.Arg[arg] = "" // default
		}

		// %arg% -> ? # # +  /
		sqlAfterArg := sql.Util_ReplaceBetweenDelimiter(query.Sql, sql.PrepareStatementDelimeter, sql.PrepareStatementAfter)

		// 쿼리 분석 후 struct 화
		// #tpl# -> tpl
		sqlAfterArgClearTpl := sql.Util_ReplaceInDelimiter(sqlAfterArg, sql.TplDelimiter, sql.TplSplit)

		psr, err := parser.New(sqlAfterArgClearTpl)
		if err != nil {
			query.ErrParser = fmt.Sprintf("%v", err)
			return nil, nil
		}
	*/

	/*
		switch data := psr.(type) {
		case *parser.Select:
			parseQuery.QueryType = parser.QueryTypeSelect
			err = Select(t.db, t.conf, query, parseQuery, data)
		case *parser.Insert:
			parseQuery.QueryType = parser.QueryTypeInsert
			err = Insert(t.conf, query, parseQuery, data)
		case *parser.Update:
			parseQuery.QueryType = parser.QueryTypeUpdate
			err = Update(t.conf, query, parseQuery, data)
		case *parser.Delete:
			parseQuery.QueryType = parser.QueryTypeDelete
			err = Delete(t.conf, query, parseQuery, data)
		}

		if err != nil {
			query.ErrQuery = fmt.Sprintf("%v", err)
			return nil, nil
		}

		// query 데이터 구성 후처리
		{
			// sql 문 복사 ( #이름# -> %s 로 변경 )
			sqlAfterArgTpl := sql.Util_ReplaceBetweenDelimiter(sqlAfterArg, sql.TplDelimiter, sql.TplAfter)
			parseQuery.Query = sqlAfterArgTpl
		}
		return parseQuery, nil
	*/
}
