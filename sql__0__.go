package bp

import (
	"module/debug/logf"

	. "github.com/blastrain/vitess-sqlparser/sqlparser"
)

//------------------------------------------------------------------------------------------------//
// SQL 파싱 결과 구조체
type I_SQL interface {
	i_sql()
}

type T_SQL struct {
}

func (t *T_SQL) Get_parser(_s_sql string) (i_sql I_SQL, err error) {
	stmt, err := Parse(_s_sql)
	if err != nil {
		return nil, err
	}

	switch data := stmt.(type) {
	case *Select:
		pt_sql__select := &T_SQL__select{}
		err = pt_sql__select.parser_result__to__struct(data)
		i_sql = pt_sql__select
	case *Insert:
		pt_sql__insert := &T_SQL__insert{}
		err = pt_sql__insert.parser_result__to__struct(data)
		i_sql = pt_sql__insert
	case *Update:
		pt_sql__update := &T_SQL__update{}
		err = pt_sql__update.parser_result__to__struct(data)
		i_sql = pt_sql__update
	case *Delete:
		pt_sql__delete := &T_SQL__delete{}
		err = pt_sql__delete.parser_result__to__struct(data)
		i_sql = pt_sql__delete
	default:
		logf.GC.Fatal("bp", "parser error")
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return i_sql, nil
}
