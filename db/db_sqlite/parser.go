package db_sqlite

import (
	"github.com/CovenantSQL/sqlparser"
	"github.com/gokch/ornn/db"
)

// TODO
type Parser struct {
}

func (t *Parser) Parse(sql string) (*db.ParseQuery, error) {
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		return nil, err
	}
	switch stmt := stmt.(type) {
	case *sqlparser.Select:
		_ = stmt
	case *sqlparser.Insert:
	}
	parseQuery := &db.ParseQuery{}
	return parseQuery, nil
}
