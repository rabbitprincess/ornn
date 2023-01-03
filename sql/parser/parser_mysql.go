package parser

import (
	"github.com/gokch/ornn/db"
	tiparser "github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
)

// TODO
type ParserMysql struct {
}

func (p *ParserMysql) Parse(sql string) (*db.ParseQuery, error) {
	ps := tiparser.New()
	stmtNodes, _, err := ps.Parse(sql, "", "")
	if err != nil {
		return nil, err
	}

	for _, n := range stmtNodes {
		switch n.(type) {
		case ast.DDLNode:
			break
		}
	}
	parseQuery := &db.ParseQuery{}

	return parseQuery, nil
}
