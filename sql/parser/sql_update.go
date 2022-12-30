package parser

import (
	"log"
	"strings"

	"github.com/blastrain/vitess-sqlparser/sqlparser"
)

type Update struct {
	TableAs []*TableAs
	Field   []*Field
}

func (t *Update) parse(psr *sqlparser.Update) error {
	for _, tableExpr := range psr.TableExprs {
		tableAs := &TableAs{}

		switch data := tableExpr.(type) {
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
		t.AddTable(tableAs)
	}

	// field value
	for _, pt_expr := range psr.Exprs {

		pt_field := &Field{}
		pt_field.TableName = pt_expr.Name.Qualifier.Name.String()
		pt_field.FieldName = pt_expr.Name.Name.String()
		switch data := pt_expr.Expr.(type) {
		case *sqlparser.SQLVal:
			{
				pt_field.Val = data.Val
			}
		case *sqlparser.BinaryExpr:
			{
				data_left, is_ok_left := data.Left.(*sqlparser.SQLVal)
				data_right, is_ok_right := data.Right.(*sqlparser.SQLVal)
				if is_ok_left == true && is_ok_right == true {
					// 양쪽 다 val 일 경우 - 에러 ( set u8_num = %u8_num% + %u8_num% )
					log.Fatal("need more programming")
				} else if is_ok_left == true {
					// 왼쪽이 val 일 경우 ( set u8_num = %u8_num% + u8_num )
					pt_field.Val = data_left.Val
				} else if is_ok_right == true {
					// 오른쪽이 val 일 경우 ( set u8_num = u8_num + %u8_num% )
					pt_field.Val = data_right.Val
				}
			}
		case *sqlparser.FuncExpr:
			{
				s_func_name := strings.ToLower(data.Name.String())
				switch s_func_name {
				case "ifnull": // 임시 - 개선 필요 - 현재 ifnull 만 적용
					{
						// update 시 nil 값을 입력하면 업데이트 하지 않도록 하기 위함
						// sql - "seq = ifnull(%seq%, seq)" 식으로 작성
						pt_alised_expr, is_ok := data.Exprs[0].(*sqlparser.AliasedExpr)
						if is_ok == false {
							log.Fatal("need more programming")
						}
						pt_val, is_ok := pt_alised_expr.Expr.(*sqlparser.SQLVal)
						pt_field.Val = pt_val.Val
						// is_pointer 플래그 추가
						// nil 을 입력받을 수 있어야 함으로 포인터로 세팅
					}
				default:
					{
						log.Fatal("need more programming")
					}
				}
			}
		default:
			log.Fatal("need more programming")
		}

		t.AddField(pt_field)
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

func (t *Update) GetTableNames() []string {
	tableName := make([]string, len(t.TableAs))

	for i, pt_table := range t.TableAs {
		tableName[i] = pt_table.Table
	}
	return tableName
}

func (t *Update) AddField(field *Field) {
	if t.Field == nil {
		t.Field = make([]*Field, 0, 10)
	}

	t.Field = append(t.Field, field)
}

func (t *Update) AddTable(tableAs *TableAs) {
	if t.TableAs == nil {
		t.TableAs = make([]*TableAs, 0, 10)
	}
	t.TableAs = append(t.TableAs, tableAs)
}
