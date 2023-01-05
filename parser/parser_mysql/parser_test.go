package parser_mysql

import (
	"fmt"
	"testing"

	"ariga.io/atlas/sql/schema"
	"github.com/gokch/ornn/atlas"
	"github.com/gokch/ornn/config"
	"github.com/gokch/ornn/db/db_mysql"
	tiparser "github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/stretchr/testify/require"
)

func TestParseMysqlSelect(t *testing.T) {
	db, err := db_mysql.New("127.0.0.1", "3306", "root", "951753ck", "test")
	require.NoError(t, err)

	atlas := atlas.New(atlas.DbTypeMaria, db)
	sc, err := atlas.InspectSchema()
	require.NoError(t, err)
	sch := config.Schema{
		DbType: atlas.DbType,
		Schema: sc,
	}

	myps := New(&sch)

	// parse
	// stmtNodes, _, err := tiparser.New().Parse("select seq, id from user where id = ? and seq = b limit 123 offset 456;", "", "")
	// stmtNodes, _, err := tiparser.New().Parse("SELECT Orders.OrderID, Customers.CustomerName, Orders.OrderDate FROM Orders INNER JOIN Customers ON Orders.CustomerID=Customers.CustomerID;", "", "")

	stmtNodes, _, err := tiparser.New().Parse("select * from user where id = ?", "", "")

	require.NoError(t, err)
	for _, stmtNode := range stmtNodes {
		selectStmt := stmtNode.(*ast.SelectStmt)
		// from
		tables := ParseJoinToTables(selectStmt.From.TableRefs)
		for _, tableExpr := range tables {
			tblName := ParseTableName(tableExpr)
			fmt.Println(tblName)
		}

		tableName := selectStmt.From.TableRefs.Left.(*ast.TableSource).Source.(*ast.TableName).Name.O
		table, _ := sch.Table(tableName)
		fmt.Println("table name :", tableName)

		// select
		// select * 일 경우 schema 의 모든 필드 추출
		if selectStmt.Fields.Fields[0].WildCard != nil {
			tbl, ok := sch.Table(tableName)
			if ok == true {
				for _, col := range tbl.Columns {
					fmt.Printf("name : %s | db type : %s | golang Type : %s\n", col.Name, col.Type.Raw, myps.ConvType(col.Type.Raw))
				}
			}
		} else {
			// select * 외의 경우
			for i, field := range selectStmt.Fields.Fields {
				switch fieldExpr := field.Expr.(type) {
				case *ast.ColumnNameExpr:
					col := fieldExpr
					var colName, colType string
					colName = col.Name.Name.O
					colType, _ = sch.GetFieldType(tableName, colName)
					fmt.Printf("select %d | name : %s | db type : %s | golang Type : %s\n", i, colName, colType, myps.ConvType(colType))
				}
			}
		}

		// visit 하면서 재귀적으로 where 필드 추출
		parseSelectWhere(selectStmt.Where, table)
	}
}

// set parseQuery rets recursive where expr
func parseSelectWhere(whereExpr ast.ExprNode, table *schema.Table) {
	switch whereExpr := whereExpr.(type) {
	case *ast.BinaryOperationExpr:
		// left
		switch leftExpr := whereExpr.L.(type) {
		case *ast.ColumnNameExpr:
			col := leftExpr
			var colName, colType string
			colName = col.Name.Name.O
			colTable, _ := table.Column(colName)
			if colTable != nil {
				colType = colTable.Type.Raw
			}
			fmt.Printf("where name : %s | db type : %s\n", colName, colType)
		}
		// right
		switch rightExpr := whereExpr.R.(type) {
		case *ast.ColumnNameExpr:
			col := rightExpr
			var colName, colType string
			colName = col.Name.Name.O
			colTable, _ := table.Column(colName)
			if colTable != nil {
				colType = colTable.Type.Raw
			}
			fmt.Printf("where name : %s | db type : %s\n", colName, colType)
		}
		// recursive
		parseSelectWhere(whereExpr.L, table)
		parseSelectWhere(whereExpr.R, table)
	}
}

func TestParseMysqlInsert(t *testing.T) {
	db, err := db_mysql.New("127.0.0.1", "3306", "root", "951753ck", "test")
	require.NoError(t, err)
	atlas := atlas.New(atlas.DbTypeMaria, db)
	sc, err := atlas.InspectSchema()
	require.NoError(t, err)

	myps := New(&config.Schema{
		DbType: atlas.DbType,
		Schema: sc,
	})

	parsedQuery, err := myps.Parse("insert into user VALUES(1, ?, ?, ?);")
	require.NoError(t, err)
	fmt.Println(parsedQuery)
}
