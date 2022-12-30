package parser

import (
	"fmt"
	"log"

	"github.com/blastrain/vitess-sqlparser/sqlparser"
)

type Insert struct {
	TableName string
	Fields    []*Field
}

func (t *Insert) parse(psr *sqlparser.Insert) error {
	// table name
	t.TableName = psr.Table.Name.String()

	switch row := psr.Rows.(type) {
	case sqlparser.Values:
		if len(row) != 1 {
			// bp config 에서 multi insert 쿼리 직접 입력 금지
			// multi insert 는 bp config 의 mutli insert option 을 true 로 설정하고, 단일 쿼리를 그대로 사용
			log.Fatal("parser error - multi insert query is forbidden. use multi insert option")
		}

		lenField := len(psr.Columns)
		lenVal := len(row[0])
		if lenField != lenVal && lenField != 0 { // 0 일때는 필드명이 모두 비어있는 상태
			return fmt.Errorf("field name length is not same as value length")
		}

		for i, val := range row[0] {
			field := &Field{}
			if psr.Columns != nil {
				field.FieldName = psr.Columns[i].String()
			}

			// value 값이 NULL 이면 처리하지 않음
			if _, ok := val.(*sqlparser.NullVal); ok == true {
				field.Val = []byte("")
				t.addField(field)
				continue
			}

			if sqlVal, ok := val.(*sqlparser.SQLVal); ok == true {
				field.Val = sqlVal.Val
				t.addField(field)
				continue
			}

			return fmt.Errorf("parser error - unexpacted type of field value")
		}
	case *sqlparser.Select:
		// 임시 - 작업필요
		// -> insert 에 입력값으로 select query 를 사용하는 경우로
		// -> 추후 작업할때 where 절에 ? 를 처리해주면 됨
		log.Fatal("need more programming")
	default:
		log.Fatal("need more programming")
	}

	return nil
}

func (t *Insert) addField(field *Field) {
	if t.Fields == nil {
		t.Fields = make([]*Field, 0, 10)
	}

	t.Fields = append(t.Fields, field)
}
