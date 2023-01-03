package parser

import (
	"fmt"
	"log"
	"strconv"

	"github.com/blastrain/vitess-sqlparser/sqlparser"
)

type Select struct {
	FieldAs []*FieldAs
	TableAs []*TableAs
	Offset  *int64
	Limit   *int64
}

func (t *Select) parse(psr *sqlparser.Select) error {
	// from
	for _, from := range psr.From {
		tableAs := &TableAs{}

		switch data := from.(type) {
		case *sqlparser.AliasedTableExpr: // 단순 테이블
			tableAs.Table = data.Expr.(sqlparser.TableName).Name.String()
			tableAs.As = data.As.String()
		case *sqlparser.ParenTableExpr:
			// 임시 - 작업필요
			// -> 반드시 sub query 를 재귀호출로 해체하여 제일 외부 에 있는 () 에 대해 서만 table list 에 남긴다. = *  타입 지정 문제
			log.Fatal("need more programming")
		case *sqlparser.JoinTableExpr:
			// 임시 - 작업필요
			// -> 반드시 sub query 를 재귀호출로 해체하여 제일 외부 에 있는 () 에 대해 서만 table list 에 남긴다. = *  타입 지정 문제
			log.Fatal("need more programming")
		}

		t.addTable(tableAs)
	}

	// select
	for _, selectExpr := range psr.SelectExprs {
		fieldAs := &FieldAs{}

		switch data := selectExpr.(type) {
		case *sqlparser.StarExpr:
			fieldAs.Field = "*"
			fieldAs.Table = data.TableName.Name.String()
			if fieldAs.Table == "" {
				// 임시 - 작업중
			}
		case *sqlparser.AliasedExpr:
			fieldAs.Field = data.Expr.(*sqlparser.ColName).Name.String()
			fieldAs.As = data.As.String()
			fieldAs.Table = data.Expr.(*sqlparser.ColName).Qualifier.Name.String()

		case sqlparser.Nextval:
			log.Fatal("need more programming")
		}

		t.addField(fieldAs)
	}

	/*
		// where
		// 재귀를 통해 모든 where 를 배열로 얻기
		// where in or not in 이면 multi arg 처리
		switch data := psr.Where.Expr.(type) {
		case *AndExpr:
			{
				// 임시 - 작업 예정 - multi where 일 경우 재귀를 통해 타입을 하나씩 얻어낸다
			}
		case *ComparisonExpr:
			{
				if data.Operator == InStr || data.Operator == NotInStr {

				}
			}
		}
	*/

	// limit, offset
	if psr.Limit != nil {
		if offset, ok := psr.Limit.Offset.(*sqlparser.SQLVal); ok == true {
			if offset.Type == sqlparser.IntVal {
				n, err := strconv.ParseInt(string(offset.Val), 10, 64)
				if err == nil {
					t.Offset = &n
				}
			}
		}
		if limit, ok := psr.Limit.Rowcount.(*sqlparser.SQLVal); ok == true {
			if limit.Type == sqlparser.IntVal {
				n, err := strconv.ParseInt(string(limit.Val), 10, 64)
				if err == nil {
					t.Limit = &n
				}
			}
		}
	}

	return nil
}

func (t *Select) GetTableNames() []string {
	tables := make([]string, len(t.TableAs))

	for i, table := range t.TableAs {
		tables[i] = table.Table
	}
	return tables
}

func (t *Select) getTableName(tableNameSql string) (tableName string, err error) {
	// 지정된 table name 이 없는경우 해석하지 않는다.
	if tableNameSql == "" {
		return "", nil
	}

	// 지정된 table name 이 있는 경우 as 인지 schema 인지 구분하여 리턴
	for _, table := range t.TableAs {
		tableNameSchema := table.get()
		if tableNameSchema == tableNameSql {
			return table.Table, nil // schema table 리턴
		}
	}

	// 테이블 이름이 매칭이 안되는 케이스 - 입력값(json) 오류
	return "", fmt.Errorf("not exist table name in select expr - table name : %s", tableNameSql)
}

func (t *Select) addField(as *FieldAs) {
	if t.FieldAs == nil {
		t.FieldAs = make([]*FieldAs, 0, 10)
	}
	t.FieldAs = append(t.FieldAs, as)
}

func (t *Select) addTable(as *TableAs) {
	if t.TableAs == nil {
		t.TableAs = make([]*TableAs, 0, 10)
	}
	t.TableAs = append(t.TableAs, as)
}
