package ornn

import (
	"fmt"
	"strings"

	"github.com/gokch/ornn/config"
	"github.com/gokch/ornn/db"
	"github.com/gokch/ornn/sql"
	"github.com/gokch/ornn/sql/parser"
)

type GenData struct {
	conf   *config.Config
	db     *db.Conn
	vendor db.Vendor

	groups []*GenDataGroup
}

type GenDataGroup struct {
	Name    string
	Queries []*GenDataQuery
}

type GenDataQuery struct {
	queryType QueryType
	groupName string
	queryName string
	query     string

	tpl Pairs
	arg Pairs
	ret Pairs

	// options
	SelectSingle     bool
	InsertMulti      bool
	UpdateNullIgnore bool
}

type QueryType int8

const (
	QueryTypeSelect QueryType = iota + 1
	QueryTypeInsert
	QueryTypeUpdate
	QueryTypeDelete
)

type Pairs struct {
	pairs []*Pair
}

type Pair struct {
	Key   string
	Value string
}

func (t *GenData) Init(conf *config.Config, db *db.Conn, vendor db.Vendor) {
	t.conf = conf
	t.db = db
	t.vendor = vendor
	t.groups = make([]*GenDataGroup, 0, 10)
}

func (t *GenData) SetData() (err error) {
	// default
	for _, group := range t.conf.Schema.Tables {
		Queries, ok := t.conf.Queries.Default[group.Name]
		if ok != true {
			continue
		}

		genGroup, err := t.SetDataGroup(group.Name, Queries)
		if err != nil {
			return err
		}
		t.Add(genGroup)
	}

	// custom
	Queries := t.conf.Queries.Custom
	genGroup, err := t.SetDataGroup("custom", Queries)
	if err != nil {
		return err
	}
	t.Add(genGroup)

	return nil
}

func (t *GenData) SetDataGroup(groupName string, queries []*config.Query) (genGroup *GenDataGroup, err error) {
	genGroup = &GenDataGroup{}
	genGroup.Init(groupName)

	for _, query := range queries {
		genQuery, err := t.SetDataQuery(groupName, query)
		if err != nil {
			return nil, err
		}
		genGroup.AddQuery(genQuery)
	}
	return genGroup, nil
}

func (t *GenData) SetDataQuery(groupName string, query *config.Query) (genQuery *GenDataQuery, err error) {
	genQuery = &GenDataQuery{}

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

		genQuery.tpl.setKV(argName, argData)
	}

	// args ( % name % )를 배열로 추출
	args, err := sql.Util_ExportBetweenDelimiter(query.Sql, sql.PrepareStatementDelimeter)
	if err != nil {
		return nil, err
	}
	genQuery.arg.setKs(args)

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

	switch data := psr.(type) {
	case *parser.Select:
		genQuery.queryType = QueryTypeSelect
		err = t.Select(t.conf, query, genQuery, data)
	case *parser.Insert:
		genQuery.queryType = QueryTypeInsert
		err = t.Insert(t.conf, query, genQuery, data)
	case *parser.Update:
		genQuery.queryType = QueryTypeUpdate
		err = t.Update(t.conf, query, genQuery, data)
	case *parser.Delete:
		genQuery.queryType = QueryTypeDelete
		err = t.Delete(t.conf, query, genQuery, data)
	}

	if err != nil {
		query.ErrQuery = fmt.Sprintf("%v", err)
		return nil, nil
	}

	// query 데이터 구성 후처리
	{
		// 그룹 이름 복사
		genQuery.groupName = groupName

		// 쿼리 이름 복사
		genQuery.queryName = query.Name

		// sql 문 복사 ( #이름# -> %s 로 변경 )
		sqlAfterArgTpl := sql.Util_ReplaceBetweenDelimiter(sqlAfterArg, sql.TplDelimiter, sql.TplAfter)
		genQuery.query = sqlAfterArgTpl

	}
	return genQuery, nil
}

func (t *GenData) Select(conf *config.Config, query *config.Query, genQuery *GenDataQuery, sqlSelect *parser.Select) error {
	// 필드 정보를 얻어온다.
	sqlWithoutWhere, _ := sql.Util_SplitByDelimiter(query.Sql, "where")
	sqlAfterArg := sql.Util_ReplaceBetweenDelimiter(sqlWithoutWhere, sql.PrepareStatementDelimeter, sql.PrepareStatementAfter)
	sqlAfterArgClearTpl := sql.Util_ReplaceInDelimiter(sqlAfterArg, sql.TplDelimiter, sql.TplSplit)

	job, err := t.db.Begin()
	if err != nil {
		return err
	}
	rows, err := job.Query(sqlAfterArgClearTpl)
	if err != nil {
		return err
	}

	cols, err := rows.ColumnTypes()
	if err != nil {
		return err
	}

	for _, col := range cols {
		var fieldName, fieldType string
		fieldName = col.Name()
		fieldType = query.GetFieldType(fieldName)

		// if custom type is not defined, get database type
		if fieldType == "" {
			colType := col.DatabaseTypeName()
			fieldType = t.vendor.ConvType(colType)
		}
		genQuery.ret.setKV(fieldName, fieldType)
	}

	// single select 처리
	// 코드 생성 시 단일 구조체 반환 목적
	if sqlSelect.Limit != nil && *(sqlSelect.Limit) == 1 {
		genQuery.SelectSingle = true
	}
	return nil
}

func (t *GenData) Insert(conf *config.Config, query *config.Query, genQuery *GenDataQuery, sqlInsert *parser.Insert) error {
	// 필드 정보를 얻어온다.
	schemaTable, exist := query.Schema.Table(sqlInsert.TableName)
	if exist != true {
		return fmt.Errorf("table name is not exist | table name - %s", sqlInsert.TableName)
	}

	// 스키마와 파서의 전체 필드 숫자가 다르면 -> 파서에서 모든 필드 이름이 제공되어야 함 -> 하나라도 없으면 에러
	if len(sqlInsert.Fields) != len(schemaTable.Columns) {
		for _, field := range sqlInsert.Fields {
			if field.FieldName == "" {
				return fmt.Errorf("field name is empty")
			}
		}
	} else {
		// 스키마와 파서의 전체 필드수가 같으면 -> 파서에서 모든 필드 이름이 없어도 가능 -> 스키마에서 추출하여 모든 필드명을 채움
		for i, field := range sqlInsert.Fields {
			field.FieldName = schemaTable.Columns[i].Name
		}
	}

	// 필드 이름을 모두 채운 상태에서 처리 시작
	for _, field := range sqlInsert.Fields {
		// 입력값이 ? (arg) 형식이 아니면 func arg 를 만들 필요가 없음으로 continue
		if sql.Util_IsParserValArg(field.Val) == false {
			continue
		}

		// 입력값이 ? (arg) 일 때만 필드이름 조사 = func arg 의 name 으로 활용
		schemaField, exist := schemaTable.Column(field.FieldName)
		if exist != true {
			return fmt.Errorf("not exist field in schema | field name : %s", field.FieldName)
		}

		genQuery.arg.setKV(field.FieldName, t.vendor.ConvType(schemaField.Type.Raw))
	}

	// multi insert 처리
	genQuery.InsertMulti = query.InsertMulti

	return nil
}

func (t *GenData) Update(conf *config.Config, query *config.Query, genQuery *GenDataQuery, sqlUpdate *parser.Update) error {
	// set
	for _, field := range sqlUpdate.Field {
		// 입력값이 ? (arg) 형식이 아니면 func arg 를 만들 필요가 없음으로 continue
		if sql.Util_IsParserValArg(field.Val) == false {
			continue
		}

		fieldName := field.FieldName
		tableName := field.TableName

		// 정의된 table name 이 없으면 update 대상 테이블 중 매칭되는 테이블을 찾는다
		/*
			if tableName == "" {
				tables := sqlUpdate.GetTableNames()
				tablesMatch, err := query.Schema.GetTableFieldMatched(fieldName, tables)
				if err != nil {
					return err
				}

				// parse 에러 처리
				{
					// 두개 이상의 테이블이 매칭됨
					if len(tablesMatch) > 1 {
						var dup string
						for _, table := range tablesMatch {
							dup += fmt.Sprintf("%s, ", table)
						}
						dup = dup[:len(dup)-2]
						return fmt.Errorf("duplicated field name in multiple table | field name - %s | tables name - %s", fieldName, dup)
					}
					// 매칭되는 테이블이 한개도 없음
					if len(tablesMatch) == 0 {
						return fmt.Errorf("no tables match the field | field name - %s", fieldName)
					}
				}

				// 테이블 이름 설정
				tableName = tablesMatch[0]
			}
		*/
		// 테이블과 필드 이름을 이용해 필드 타입을 찾아낸다
		var genType string
		{
			schemaTable, exist := query.Schema.Table(tableName)
			if exist != true {
				return fmt.Errorf("not exist table | table name - %s", tableName)
			}
			schemaField, exist := schemaTable.Column(fieldName)
			if exist != true {
				return fmt.Errorf("not exist field | field name - %s", field.FieldName)
			}
			genType = t.vendor.ConvType(schemaField.Type.Raw)
		}

		genQuery.arg.setKV(field.FieldName, genType)
	}
	// update 시 null 값 ignore 처리
	genQuery.UpdateNullIgnore = query.UpdateNullIgnore

	return nil
}

func (t *GenData) Delete(conf *config.Config, query *config.Query, genQuery *GenDataQuery, sqlDelete *parser.Delete) error {
	return nil
}

func (t *GenData) Add(group *GenDataGroup) {
	if t.groups == nil {
		t.groups = make([]*GenDataGroup, 0, 10)
	}
	t.groups = append(t.groups, group)
}

func (t *GenDataGroup) Init(Name string) {
	t.Name = Name
}

func (t *GenDataGroup) AddQuery(query *GenDataQuery) {
	if t.Queries == nil {
		t.Queries = make([]*GenDataQuery, 0, 10)
	}
	t.Queries = append(t.Queries, query)
}

func (t *Pairs) setKs(Keys []string) {
	for _, key := range Keys {
		t.setKV(key, "")
	}
}

func (t *Pairs) setKV(key string, valueNew string) {
	if t.pairs == nil {
		t.pairs = make([]*Pair, 0, 10)
	}

	for _, field := range t.pairs {
		if field.Key == key {
			field.Value = valueNew
			return
		}
	}

	t.pairs = append(t.pairs, &Pair{
		Key:   key,
		Value: valueNew,
	})
}
