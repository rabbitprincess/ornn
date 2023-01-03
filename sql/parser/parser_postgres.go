package parser

import (
	"github.com/auxten/postgresql-parser/pkg/sql/parser"
	"github.com/auxten/postgresql-parser/pkg/walk"
)

// TODO
type ParserPostgres struct {
}

func (p *ParserPostgres) Parse(sql string) (*ParseQuery, error) {
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

	parseQuery := &ParseQuery{}

	return parseQuery, nil
}
