package parser_postgres

import (
	"fmt"
	"testing"

	sqlparser "github.com/cockroachdb/cockroachdb-parser/pkg/sql/parser"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree"
	"github.com/stretchr/testify/require"
)

func TestParseMysqlSelect(t *testing.T) {
	sql := "select * from test where a = $1 and b = 2"
	stmts, err := sqlparser.Parse(sql)
	require.NoError(t, err)
	require.Equal(t, 1, len(stmts))

	stmt := stmts[0].AST.(*tree.Select)
	// selectStmt := stmt.Select.(*tree.SelectClause)
	fmt.Println(stmt)
}

func TestParseMysqlInsert(t *testing.T) {
	sql := "insert into test (a,b,c) values (1,$1,$2), (2,3,4)"
	stmts, err := sqlparser.Parse(sql)
	require.NoError(t, err)
	require.Equal(t, 1, len(stmts))

	insertStmt := stmts[0].AST.(*tree.Insert)
	fmt.Println(insertStmt.Table)
}
