package db_postgres

import (
	"github.com/auxten/postgresql-parser/pkg/sql/parser"
	"github.com/auxten/postgresql-parser/pkg/walk"
	"github.com/gokch/ornn/db"
)

// TODO
type Parser struct {
}

func (p *Parser) Parse(sql string) (*db.ParseQuery, error) {
	w := &walk.AstWalker{
		// Fn: p.walker,
	}

	stmts, err := parser.Parse(sql)
	if err != nil {
		return nil, err
	}

	_, err = w.Walk(stmts, nil)
	if err != nil {
		return nil, err
	}

	parseQuery := &db.ParseQuery{}

	return parseQuery, nil
}
