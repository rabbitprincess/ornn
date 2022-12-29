package sql

import (
	"log"

	parser "github.com/blastrain/vitess-sqlparser/sqlparser"
)

type Delete struct {
	TableAs []*TableAs
}

func (t *Delete) parse(psr *parser.Delete) error {
	for _, tableExprs := range psr.TableExprs {
		tableAs := &TableAs{}

		switch data := tableExprs.(type) {
		case *parser.AliasedTableExpr: // 단순 테이블
			tableAs.Tbl = data.Expr.(parser.TableName).Name.String()
			tableAs.As = data.As.String()
		case *parser.ParenTableExpr:
			// 임시 - 작업필요
			// -> 반드시 sub query 를 재귀호출로 해체하여 제일 외부 에 있는 () 에 대해 서만 table list 에 남긴다. = *  타입 지정 문제
			log.Fatal("need more programming")
		case *parser.JoinTableExpr:
			// 임시 - 작업필요
			// -> 반드시 sub query 를 재귀호출로 해체하여 제일 외부 에 있는 () 에 대해 서만 table list 에 남긴다. = *  타입 지정 문제
			log.Fatal("need more programming")
		}
		t.AddTable(tableAs)
	}

	/*
		// where
		// 재귀를 통해 모든 where 문을 배열로 전처리
		// select 일 경우 select 문 재귀 처리
		// in 또는 not in 이면 multi arg 처리
		switch data := parser.Where.Expr.(type) {
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

	return nil
}

func (t *Delete) AddTable(as *TableAs) {
	if t.TableAs == nil {
		t.TableAs = make([]*TableAs, 0, 10)
	}
	t.TableAs = append(t.TableAs, as)
}
