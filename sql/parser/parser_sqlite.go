package parser

import (
	"github.com/CovenantSQL/sqlparser"
)

// TODO
type ParserSqlite struct {
}

func (t *ParserSqlite) Parse(sql string) (*ParseQuery, error) {
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		return nil, err
	}
	switch stmt := stmt.(type) {
	case *sqlparser.Select:
		_ = stmt
	case *sqlparser.Insert:
	}
	parseQuery := &ParseQuery{}
	return parseQuery, nil
}
