package db_postgres

import (
	"github.com/auxten/postgresql-parser/pkg/sql/parser"
	"github.com/auxten/postgresql-parser/pkg/walk"
)

// TODO
type Parser struct {
}

func (p *Parser) Parse(sql string) error {
	w := &walk.AstWalker{
		// Fn: p.walker,
	}

	stmts, err := parser.Parse(sql)
	if err != nil {
		return err
	}

	_, err = w.Walk(stmts, nil)
	return err
}
