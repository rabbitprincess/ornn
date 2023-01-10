package parser_sqlite

import (
	"fmt"
	"testing"

	"github.com/CovenantSQL/sqlparser"
	"github.com/stretchr/testify/require"
)

func TestParseSqliteSelect(t *testing.T) {
	sql := "select a,b,c from test where a=1 and b=? and c=?"
	stmtNodes, err := sqlparser.Parse(sql)
	require.NoError(t, err)

	selectStmt := stmtNodes.(*sqlparser.Select)
	/*
		ret := ParseWhereToFields(selectStmt.Where.Expr)
		for _, v := range ret {
			fmt.Printf("ret : %d %T | %d %T\n", v.left, v.left, v.right, v.right)
			fmt.Println(string(v.right.(*sqlparser.SQLVal).Val), v.right.(*sqlparser.SQLVal).Type)
		}
	*/
	// select
	switch sel := selectStmt.SelectExprs[0].(type) {
	case *sqlparser.StarExpr: // select * 일 경우 schema 의 모든 인자 추출
		fmt.Println(sel)
	case *sqlparser.AliasedExpr:
		fmt.Printf("%T\n", sel.Expr)
		switch expr := sel.Expr.(type) {
		case *sqlparser.ColName:
			fmt.Println(expr)

		}
	default:
		panic("need more programming")
	}
}

func TestParseSqliteInsert(t *testing.T) {
	sql := "insert into test (a,b,c) values (?,?,?)"
	stmtNodes, err := sqlparser.Parse(sql)
	require.NoError(t, err)

	insertStmt := stmtNodes.(*sqlparser.Insert)
	fmt.Println(insertStmt.Table.Name.String())
	fmt.Println(insertStmt.OnDup)

}
