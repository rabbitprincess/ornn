package bp

import (
	"module/debug/logf"

	. "github.com/blastrain/vitess-sqlparser/sqlparser"
)

type T_SQL__delete struct {
	Arrpt_table []*T_SQL__table_as
}

func (t *T_SQL__delete) i_sql() {}

func (t *T_SQL__delete) parser_result__to__struct(_pt_parser *Delete) error {
	for _, i_table := range _pt_parser.TableExprs {
		pt_table_as := &T_SQL__table_as{}

		switch data := i_table.(type) {
		case *AliasedTableExpr: // 단순 테이블
			pt_table_as.S_table_name = data.Expr.(TableName).Name.String()
			pt_table_as.S_as = data.As.String()
		case *ParenTableExpr:
			// 임시 - 작업필요
			// -> 반드시 sub query 를 재귀호출로 해체하여 제일 외부 에 있는 () 에 대해 서만 table list 에 남긴다. = *  타입 지정 문제
			logf.GC.Fatal("bp", "parser error - need more programming")
		case *JoinTableExpr:
			// 임시 - 작업필요
			// -> 반드시 sub query 를 재귀호출로 해체하여 제일 외부 에 있는 () 에 대해 서만 table list 에 남긴다. = *  타입 지정 문제
			logf.GC.Fatal("bp", "parser error - need more programming")
		}
		t.table__add(pt_table_as)
	}

	/*
		// where
		// 재귀를 통해 모든 where 문을 배열로 전처리
		// select 일 경우 select 문 재귀 처리
		// in 또는 not in 이면 multi arg 처리
		switch data := _pt_parser.Where.Expr.(type) {
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

func (t *T_SQL__delete) table__add(_pt_table_as *T_SQL__table_as) {
	if t.Arrpt_table == nil {
		t.Arrpt_table = make([]*T_SQL__table_as, 0, 10)
	}
	t.Arrpt_table = append(t.Arrpt_table, _pt_table_as)
}
