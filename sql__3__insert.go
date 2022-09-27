package bp

import (
	"fmt"
	"module/debug/logf"

	. "github.com/blastrain/vitess-sqlparser/sqlparser"
)

type T_SQL__insert struct {
	S_table_name string
	Arrpt_field  []*T_SQL__field_value
}

func (t *T_SQL__insert) i_sql() {}

func (t *T_SQL__insert) parser_result__to__struct(_pt_parser *Insert) error {
	// table name
	t.S_table_name = _pt_parser.Table.Name.String()

	switch data_row := _pt_parser.Rows.(type) {
	case Values:
		if len(data_row) != 1 {
			// bp config 에서 multi insert 쿼리 직접 입력 금지
			// multi insert 는 bp config 의 mutli insert option 을 true 로 설정하고, 단일 쿼리를 그대로 사용
			logf.GC.Fatal("bp", "parser error - multi insert query is forbidden. use multi insert option")
		}

		n_len_field := len(_pt_parser.Columns)
		n_len_val := len(data_row[0])
		if n_len_field != n_len_val && n_len_field != 0 { // 0 일때는 필드명이 모두 비어있는 상태
			return fmt.Errorf("field name length is not same as value length")
		}

		for i, val := range data_row[0] {
			pt_field := &T_SQL__field_value{}
			if _pt_parser.Columns != nil {
				pt_field.S_field_name = _pt_parser.Columns[i].String()
			}

			// value 값이 NULL 이면 처리하지 않음
			_, is_ok__null := val.(*NullVal)
			if is_ok__null == true {
				pt_field.BT_val = []byte("")
				t.field__add(pt_field)
				continue
			}

			i_sqlval, is_ok__sqlval := val.(*SQLVal)
			if is_ok__sqlval == true {
				pt_field.BT_val = i_sqlval.Val
				t.field__add(pt_field)
				continue
			}

			return fmt.Errorf("parser error - unexpacted type of field value")
		}
	case *Select:
		// 임시 - 작업필요
		// -> insert 에 입력값으로 select query 를 사용하는 경우로
		// -> 추후 작업할때 where 절에 ? 를 처리해주면 됨
		logf.GC.Fatal("bp", "parser error - need more programming")
	default:
		logf.GC.Fatal("bp", "parser error - need more programming")
	}

	return nil
}

func (t *T_SQL__insert) field__add(_pt_field_value *T_SQL__field_value) {
	if t.Arrpt_field == nil {
		t.Arrpt_field = make([]*T_SQL__field_value, 0, 10)
	}

	t.Arrpt_field = append(t.Arrpt_field, _pt_field_value)
}
