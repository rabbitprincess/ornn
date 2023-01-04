package parser

import (
	"fmt"
	"testing"

	tiparser "github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/stretchr/testify/require"
)

func TestParseMysql(t *testing.T) {
	ps := tiparser.New()
	stmtNodes, _, err := ps.Parse("select * from test_table where id = ?", "", "")
	require.NoError(t, err)

	for _, stmtNode := range stmtNodes {
		selectStmt := stmtNode.(*ast.SelectStmt)
		for _, tbl := range selectStmt.TableHints {
			fmt.Println(tbl, "ef")
		}

		// from
		fmt.Println(selectStmt.From.TableRefs, "tblRefs")
		fmt.Println(selectStmt.From.TableRefs, "tblRefs")
		fmt.Println(selectStmt.From.Text())
	}
}
