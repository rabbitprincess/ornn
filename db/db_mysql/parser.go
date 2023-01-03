package db_mysql

import (
	tiparser "github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
)

// TODO
type Parser struct {
}

// ParserMysql ...
func (p *Parser) Parse(sql string) error {
	ps := tiparser.New()
	stmtNodes, _, err := ps.Parse(sql, "", "")
	if err != nil {
		return err
	}

	for _, n := range stmtNodes {
		switch n.(type) {
		case ast.DDLNode:
			break
		}
	}

	return nil
}
