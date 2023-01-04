package parser_postgres

import (
	"ariga.io/atlas/sql/schema"
	sqlparser "github.com/auxten/postgresql-parser/pkg/sql/parser"
	"github.com/auxten/postgresql-parser/pkg/walk"
	"github.com/gokch/ornn/parser"
)

func New(sch *schema.Schema) parser.Parser {
	return &Parser{
		sch: sch,
	}
}

// TODO
type Parser struct {
	sch *schema.Schema
}

func (p *Parser) Parse(sql string) (*parser.ParsedQuery, error) {
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
