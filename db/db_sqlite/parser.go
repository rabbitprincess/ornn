package db_sqlite

import (
	"github.com/CovenantSQL/sqlparser"
)

// TODO
type Parser struct {
}

func (t *Parser) Parse(sql string) {
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		return
	}
	switch stmt := stmt.(type) {
	case *sqlparser.Select:
		_ = stmt
	case *sqlparser.Insert:
	}
}
