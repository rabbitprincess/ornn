package bp

import (
	"fmt"
	"module/debug/logf"
	"strconv"

	"github.com/blastrain/vitess-sqlparser/sqlparser"
	. "github.com/blastrain/vitess-sqlparser/sqlparser"
)

const DEF_s_delimeter__between__table_name__to__field_name string = "__"

type T_SQL__select struct {
	Arrpt_select_as []*T_SQL__select_as
	Arrpt_table     []*T_SQL__table_as
	pn8_offset      *int64
	pn8_limit       *int64
}

func (t *T_SQL__select) i_sql() {}

func (t *T_SQL__select) parser_result__to__struct(_pt_parser *Select) error {

	// from
	{
		for _, i_from := range _pt_parser.From {
			pt_table_as := &T_SQL__table_as{}

			switch data := i_from.(type) {
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
	}

	// select
	{
		for _, i_select := range _pt_parser.SelectExprs {
			pt_field_as := &T_SQL__select_as{}

			switch data := i_select.(type) {
			case *StarExpr:
				pt_field_as.S_field_name = "*"
				pt_field_as.S_table_name = data.TableName.Name.String()
				if pt_field_as.S_table_name == "" {
					// 임시 - 작업중 pt_field_as.S_table_name__sql
				}
			case *AliasedExpr:
				pt_field_as.S_field_name = data.Expr.(*ColName).Name.String()
				pt_field_as.S_as = data.As.String()
				pt_field_as.S_table_name = data.Expr.(*ColName).Qualifier.Name.String()

			case Nextval:
				logf.GC.Fatal("bp", "parser error - need more programming")
			}

			t.field__add(pt_field_as)
		}
	}

	/*
		// where
		// 재귀를 통해 모든 where 를 배열로 얻기
		// where in or not in 이면 multi arg 처리
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

	// limit, offset
	if _pt_parser.Limit != nil {
		pt_val_offset, is_ok := _pt_parser.Limit.Offset.(*sqlparser.SQLVal)
		if is_ok == true {
			if pt_val_offset.Type == sqlparser.IntVal {
				n8_offset, err := strconv.ParseInt(string(pt_val_offset.Val), 10, 64)
				if err == nil {
					t.pn8_offset = &n8_offset
				}
			}
		}
		pt_val_limit, is_ok := _pt_parser.Limit.Rowcount.(*sqlparser.SQLVal)
		if is_ok == true {
			if pt_val_limit.Type == sqlparser.IntVal {
				n8_limit, err := strconv.ParseInt(string(pt_val_limit.Val), 10, 64)
				if err == nil {
					t.pn8_limit = &n8_limit
				}
			}
		}
	}

	return nil
}

func (t *T_SQL__select) get_table_name() []string {
	arrs_table_name := make([]string, len(t.Arrpt_table))

	for i, pt_table := range t.Arrpt_table {
		arrs_table_name[i] = pt_table.S_table_name
	}
	return arrs_table_name
}

func (t *T_SQL__select) get_table_name__in_sql_name__to__schema_name(_s_table_name__in_sql string) (s_table_name string, err error) {
	// 지정된 table_name 이 없는경우 해석하지 않는다.
	if _s_table_name__in_sql == "" {
		return "", nil
	}

	// 지정된 table_name 이 있는 경우 as 인지 schema_name 인지 구분하여 리턴
	for _, pt_table := range t.Arrpt_table {
		s_table_name__in_sql := pt_table.get_table_name__in_sql()
		if s_table_name__in_sql == _s_table_name__in_sql {
			return pt_table.S_table_name, nil // schema table_name 을 리턴
		}
	}

	// 테이블 이름이 매칭이 안되는 케이스 - 입력값(json) 오류
	return "", fmt.Errorf("not exist table_name in select expr - table_name : %s", _s_table_name__in_sql)
}

//--------------------------------------------------------------------------------------------------------//

func (t *T_SQL__select) field__add(_pt_field_as *T_SQL__select_as) {
	if t.Arrpt_select_as == nil {
		t.Arrpt_select_as = make([]*T_SQL__select_as, 0, 10)
	}

	t.Arrpt_select_as = append(t.Arrpt_select_as, _pt_field_as)
}

func (t *T_SQL__select) table__add(_pt_table_as *T_SQL__table_as) {
	if t.Arrpt_table == nil {
		t.Arrpt_table = make([]*T_SQL__table_as, 0, 10)
	}
	t.Arrpt_table = append(t.Arrpt_table, _pt_table_as)
}
