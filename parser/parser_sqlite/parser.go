package parser_sqlite

import (
	"fmt"

	"github.com/CovenantSQL/sqlparser"
	"github.com/gokch/ornn/config"
	"github.com/gokch/ornn/parser"
)

// TODO
type Parser struct {
	sch *config.Schema
}

func (p *Parser) Init(sch *config.Schema) {
	p.sch = sch
}

func (t *Parser) Parse(sql string) (*parser.ParsedQuery, error) {
	stmtNodes, err := sqlparser.Parse(sql)
	if err != nil {
		return nil, err
	}

	switch stmt := stmtNodes.(type) {
	case *sqlparser.Select:
		_ = stmt
	case *sqlparser.Insert:
		_ = stmt
	default:
		return nil, fmt.Errorf("parser error | not support query statement %T", stmt)
	}

	parsedQuery := &parser.ParsedQuery{}
	return parsedQuery, nil
}
