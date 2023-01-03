package db_mysql

import (
	tiparser "github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
)

// TODO
type Parser struct {
}

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
