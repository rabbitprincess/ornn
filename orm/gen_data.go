package orm

import (
	"fmt"
	"strings"

	"github.com/gokch/go-orm-gen/config"
	"github.com/gokch/go-orm-gen/db"
	"github.com/gokch/go-orm-gen/db/db_mysql"
	"github.com/gokch/go-orm-gen/sql"
)

type GenData struct {
	db     *db.Conn
	vendor db.Vendor

	groups []*GenDataGroup
}

func (t *GenData) Init(db *db.Conn) {
	t.db = db
	t.vendor = db_mysql.NewVendor(db)
	t.groups = make([]*GenDataGroup, 0, 10)
}

func (t *GenData) SetData(config *config.Config) (err error) {
	job, err := t.db.Begin()
	if err != nil {
		return err
	}
	job.Exec("Insert from test_table")

	// group
	schema := &config.Schema

	for _, table := range config.Schema.Tables {
		genGroup := &GenDataGroup{}
		genGroup.Set(table.Name)
		t.Add(genGroup)

		// func
		for _, query := range table.Queries {
			genQuery := &GenDataQuery{}

			// set args
			// tpl args ( # name # )를 배열로 추출
			tpls, err := sql.Util_ExportBetweenDelimiter(query.Sql, sql.DEF_s_sql__tpl__delimiter)
			if err != nil {
				return err
			}

			for _, tpl := range tpls {
				tmps := strings.Split(tpl, sql.DEF_s_sql__tpl__split)
				var argName string
				var argData string
				if len(tmps) == 1 {
					argName = tmps[0]
					argData = ""
				} else if len(tmps) == 2 {
					argName = tmps[0]
					argData = tmps[1]
				} else {
					return fmt.Errorf("tpl format is wrong - %s", tpl)
				}

				genQuery.tpl.setKV(argName, argData)
			}

			// args ( % name % )를 배열로 추출
			args, err := sql.Util_ExportBetweenDelimiter(query.Sql, sql.DEF_s_sql__prepare_statement__delimiter)
			if err != nil {
				return err
			}
			genQuery.arg.setKs(args)

			// %arg% -> ? # # +  /
			sqlAfterArg := sql.Util_ReplaceBetweenDelimiter(query.Sql, sql.DEF_s_sql__prepare_statement__delimiter, sql.DEF_s_sql__prepare_statement__after)

			// 쿼리 분석 후 struct 화
			// #tpl# -> tpl
			sqlAfterArgClearTpl := sql.Util_ReplaceInDelimiter(sqlAfterArg, sql.DEF_s_sql__tpl__delimiter, sql.DEF_s_sql__tpl__split)

			isql, err := sql.New(sqlAfterArgClearTpl)
			if err != nil {
				query.ErrParser = fmt.Sprintf("%v", err)
				continue
			}

			switch data_parser := isql.(type) {
			case *sql.Select:
				err = t.Select(config, table, query, genQuery, data_parser)
			case *sql.Insert:
				err = t.Insert(config, schema, table, query, genQuery, data_parser)
			case *sql.Update:
				err = t.Update(config, schema, table, query, genQuery, data_parser)
			case *sql.Delete:
				err = t.Delete(config, table, query, genQuery, data_parser)
			}

			if err != nil {
				query.ErrQuery = fmt.Sprintf("%v", err)
				continue
			}

			// pt_gen__query 데이터 구성 후처리
			{
				// 그룹 이름 복사
				genQuery.tableName = table.Name

				// 쿼리 이름 복사
				genQuery.queryName = query.Name

				// sql 문 복사 ( #이름# -> %s 로 변경 )
				sqlAfterArgTpl := sql.Util_ReplaceBetweenDelimiter(sqlAfterArg, sql.DEF_s_sql__tpl__delimiter, sql.DEF_s_sql__tpl__after)
				genQuery.query = sqlAfterArgTpl

				// group list 에 func 추가
				genGroup.AddQuery(genQuery)
			}
		}
	}

	return nil
}

func (t *GenData) Select(conf *config.Config, table *config.Table, query *config.Query, genQuery *GenDataQuery, sqlSelect *sql.Select) error {
	pt_rds := &DbVendor{}
	pt_rds.Init(db_mysql.NewVendor(t.db))

	genQuery.queryType = QueryTypeSelect

	// 필드 정보를 얻어온다.
	{
		s_sql, _ := sql.Util_SplitByDelimiter(query.Sql, "where")
		s_sql__after_arg := sql.Util_ReplaceBetweenDelimiter(s_sql, sql.DEF_s_sql__prepare_statement__delimiter, sql.DEF_s_sql__prepare_statement__after)
		s_sql__after_arg_clear_tpl := sql.Util_ReplaceInDelimiter(s_sql__after_arg, sql.DEF_s_sql__tpl__delimiter, sql.DEF_s_sql__tpl__split)

		job, err := t.db.Begin()
		if err != nil {
			return err
		}
		rows, err := job.Query(s_sql__after_arg_clear_tpl)
		if err != nil {
			return err
		}

		cols, err := rows.ColumnTypes()
		if err != nil {
			return err
		}

		for _, col := range cols {
			var fldName, fldType string
			fldName = col.Name()
			fldType = query.GetFieldType(fldName)

			// if custom type is not defined, get database type
			if fldType == "" {
				colType := col.DatabaseTypeName()
				fldType = pt_rds.vendor.ConvType(colType)
			}
			genQuery.ret.setKV(fldName, fldType)
		}
	}
	// single select 처리
	// 코드 생성 시 단일 구조체 반환 목적
	if sqlSelect.Limit != nil && *(sqlSelect.Limit) == 1 {
		genQuery.isSelectSingle = true
	}
	return nil
}

func (t *GenData) Insert(conf *config.Config, schema *config.Schema, table *config.Table, query *config.Query, genQuery *GenDataQuery, sqlInsert *sql.Insert) error {

	genQuery.queryType = QueryTypeInsert

	// 필드 정보를 얻어온다.
	{
		schemaTbl := schema.GetTable(sqlInsert.TblName)
		if schemaTbl == nil {
			return fmt.Errorf("table name is not exist | table name - %s", sqlInsert.TblName)
		}

		// 스키마와 파서의 전체 필드 숫자가 다르면 -> 파서에서 모든 필드 이름이 제공되어야 함 -> 하나라도 없으면 에러
		if len(sqlInsert.Fields) != len(schemaTbl.Fields) {
			for _, pt_field_value := range sqlInsert.Fields {
				if pt_field_value.FldName == "" {
					return fmt.Errorf("field name is empty")
				}
			}
		} else {
			// 스키마와 파서의 전체 필드수가 같으면 -> 파서에서 모든 필드 이름이 없어도 가능 -> 스키마에서 추출하여 모든 필드명을 채움
			for i, field := range sqlInsert.Fields {
				field.FldName = schemaTbl.Fields[i].Name
			}
		}

		// 필드 이름을 모두 채운 상태에서 처리 시작
		for _, field := range sqlInsert.Fields {
			// 입력값이 ? (arg) 형식이 아니면 func arg 를 만들 필요가 없음으로 continue
			if sql.Util_IsParserValArg(field.Val) == false {
				continue
			}

			// 입력값이 ? (arg) 일 때만 필드이름 조사 = func arg 의 name 으로 활용
			schemaFld := schemaTbl.GetField(field.FldName)
			if schemaFld == nil {
				return fmt.Errorf("not exist field in schema | field name : %s", field.FldName)
			}

			genQuery.arg.setKV(field.FldName, schemaFld.TypeGen)
		}
	}
	// multi insert 처리
	genQuery.InsertMulti = query.InsertMulti

	return nil
}

func (t *GenData) Update(conf *config.Config, schema *config.Schema, table *config.Table, query *config.Query, genQuery *GenDataQuery, sqlUpdate *sql.Update) error {
	genQuery.queryType = QueryTypeUpdate

	// set
	{
		for _, field := range sqlUpdate.Field {
			// 입력값이 ? (arg) 형식이 아니면 func arg 를 만들 필요가 없음으로 continue
			if sql.Util_IsParserValArg(field.Val) == false {
				continue
			}

			fldName := field.FldName
			tblName := field.TblName

			// 정의된 table name 이 없으면 update 대상 테이블 중 매칭되는 테이블을 찾는다
			if tblName == "" {
				tables := sqlUpdate.GetTableNames()
				tablesMatch, err := schema.GetTableFieldMatched(fldName, tables)
				if err != nil {
					return err
				}

				// parse 에러 처리
				{
					// 두개 이상의 테이블이 매칭됨
					if len(tablesMatch) > 1 {
						var dup string
						for _, s_table_name__match := range tablesMatch {
							dup += fmt.Sprintf("%s, ", s_table_name__match)
						}
						dup = dup[:len(dup)-2]
						return fmt.Errorf("duplicated field name in multiple table | field name - %s | tables name - %s", fldName, dup)
					}
					// 매칭되는 테이블이 한개도 없음
					if len(tablesMatch) == 0 {
						return fmt.Errorf("no tables match the field | field name - %s", fldName)
					}
				}

				// 테이블 이름 설정
				tblName = tablesMatch[0]
			}

			// 테이블과 필드 이름을 이용해 필드 타입을 찾아낸다
			var genType string
			{
				schemaTbl := schema.GetTable(tblName)
				if schemaTbl == nil {
					return fmt.Errorf("not exist table | table name - %s", tblName)
				}
				schemaFld := schemaTbl.GetField(fldName)
				if schemaFld == nil {
					return fmt.Errorf("not exist field | field name - %s", field.FldName)
				}
				genType = string(schemaFld.TypeGen)
			}

			genQuery.arg.setKV(field.FldName, genType)
		}
	}
	// update 시 null 값 ignore 처리
	genQuery.UpdateNullIgnore = query.UpdateNullIgnore

	return nil
}

func (t *GenData) Delete(conf *config.Config, table *config.Table, query *config.Query, genQuery *GenDataQuery, sqlDelete *sql.Delete) error {

	genQuery.queryType = QueryTypeDelete

	// 임시 - 할게 없음
	return nil
}

func (t *GenData) Add(group *GenDataGroup) {
	if t.groups == nil {
		t.groups = make([]*GenDataGroup, 0, 10)
	}
	t.groups = append(t.groups, group)
}

type GenDataGroup struct {
	Name    string
	Queries []*GenDataQuery
}

func (t *GenDataGroup) Set(Name string) {
	if t.Queries == nil {
		t.Queries = make([]*GenDataQuery, 0, 10)
	}
	t.Name = Name
}

func (t *GenDataGroup) AddQuery(query *GenDataQuery) {
	t.Queries = append(t.Queries, query)
}

type QueryType int8

const (
	QueryTypeSelect QueryType = iota + 1
	QueryTypeInsert
	QueryTypeUpdate
	QueryTypeDelete
)

type GenDataQuery struct {
	queryType QueryType
	tableName string
	queryName string
	query     string

	tpl genDataStruct
	arg genDataStruct
	ret genDataStruct

	isSelectSingle   bool
	InsertMulti      bool
	UpdateNullIgnore bool
}

type Pair struct {
	Key   string
	Value string
}

type genDataStruct struct {
	arrpt_pair []*Pair
}

func (t *genDataStruct) setKs(Keys []string) {
	for _, s_field_name := range Keys {
		t.setKV(s_field_name, "")
	}
}

func (t *genDataStruct) setKV(key string, valueNew string) {
	if t.arrpt_pair == nil {
		t.arrpt_pair = make([]*Pair, 0, 10)
	}

	for _, fld := range t.arrpt_pair {
		if fld.Key == key {
			fld.Value = valueNew
			return
		}
	}

	t.arrpt_pair = append(t.arrpt_pair, &Pair{
		Key:   key,
		Value: valueNew,
	})
}
