package bp

import (
	"module/debug/logf"
	"strings"

	. "github.com/blastrain/vitess-sqlparser/sqlparser"
)

type T_SQL__update struct {
	Arrpt_table []*T_SQL__table_as
	Arrpt_field []*T_SQL__field_value
}

func (t *T_SQL__update) i_sql() {}

func (t *T_SQL__update) get_table_name() []string {
	arrs_table_name := make([]string, len(t.Arrpt_table))

	for i, pt_table := range t.Arrpt_table {
		arrs_table_name[i] = pt_table.S_table_name
	}
	return arrs_table_name
}

func (t *T_SQL__update) parser_result__to__struct(_pt_parser *Update) error {
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

	// field value
	for _, pt_expr := range _pt_parser.Exprs {

		pt_field := &T_SQL__field_value{}
		pt_field.S_table_name = pt_expr.Name.Qualifier.Name.String()
		pt_field.S_field_name = pt_expr.Name.Name.String()
		switch data := pt_expr.Expr.(type) {
		case *SQLVal:
			{
				pt_field.BT_val = data.Val
			}
		case *BinaryExpr:
			{
				data_left, is_ok_left := data.Left.(*SQLVal)
				data_right, is_ok_right := data.Right.(*SQLVal)
				if is_ok_left == true && is_ok_right == true {
					// 양쪽 다 val 일 경우 - 에러 ( set u8_num = %u8_num% + %u8_num% )
					logf.GC.Fatal("bp", "parser error - need more programming")
				} else if is_ok_left == true {
					// 왼쪽이 val 일 경우 ( set u8_num = %u8_num% + u8_num )
					pt_field.BT_val = data_left.Val
				} else if is_ok_right == true {
					// 오른쪽이 val 일 경우 ( set u8_num = u8_num + %u8_num% )
					pt_field.BT_val = data_right.Val
				}
			}
		case *FuncExpr:
			{
				s_func_name := strings.ToLower(data.Name.String())
				switch s_func_name {
				case "ifnull": // 임시 - 개선 필요 - 현재 ifnull 만 적용
					{
						// update 시 nil 값을 입력하면 업데이트 하지 않도록 하기 위함
						// sql - "seq = ifnull(%seq%, seq)" 식으로 작성
						pt_alised_expr, is_ok := data.Exprs[0].(*AliasedExpr)
						if is_ok == false {
							logf.GC.Fatal("bp", "parser error - need more programming")
						}
						pt_val, is_ok := pt_alised_expr.Expr.(*SQLVal)
						pt_field.BT_val = pt_val.Val
						// is_pointer 플래그 추가
						// nil 을 입력받을 수 있어야 함으로 포인터로 세팅
					}
				default:
					{
						logf.GC.Fatal("bp", "parser error - need more programming")
					}
				}
			}
		default:
			logf.GC.Fatal("bp", "parser error - need more programming")
		}

		t.field__add(pt_field)
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

func (t *T_SQL__update) field__add(_pt_field_value *T_SQL__field_value) {
	if t.Arrpt_field == nil {
		t.Arrpt_field = make([]*T_SQL__field_value, 0, 10)
	}

	t.Arrpt_field = append(t.Arrpt_field, _pt_field_value)
}

func (t *T_SQL__update) table__add(_pt_table_as *T_SQL__table_as) {
	if t.Arrpt_table == nil {
		t.Arrpt_table = make([]*T_SQL__table_as, 0, 10)
	}
	t.Arrpt_table = append(t.Arrpt_table, _pt_table_as)
}
