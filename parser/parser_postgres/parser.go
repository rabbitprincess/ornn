package parser_postgres

import (
	sqlparser "github.com/auxten/postgresql-parser/pkg/sql/parser"
	"github.com/auxten/postgresql-parser/pkg/walk"
	"github.com/gokch/ornn/config"
	"github.com/gokch/ornn/parser"
)

// TODO
type ParserPostgres struct {
	sch *config.Schema
}

func (p *ParserPostgres) Init(sch *config.Schema) {
	p.sch = sch
}

func (p *ParserPostgres) Parse(sql string) (*parser.ParsedQuery, error) {
	stmtNodes, err := sqlparser.Parse(sql)
	if err != nil {
		return nil, err
	}

	// TODO
	w := &walk.AstWalker{
		// Fn: p.walker,
	}
	_, err = w.Walk(stmtNodes, nil)
	if err != nil {
		return nil, err
	}

	// for _, stmtNode := range stmtNodes {
	// }

	parsedQuery := &parser.ParsedQuery{}
	return parsedQuery, nil
}
