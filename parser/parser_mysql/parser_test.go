package parser_mysql

import (
	"fmt"
	"testing"

	"github.com/gokch/ornn/atlas"
	"github.com/gokch/ornn/config"
	"github.com/gokch/ornn/db/db_mysql"
	tiparser "github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/test_driver"
	"github.com/stretchr/testify/require"
)

func TestParseMysql(t *testing.T) {
	db, err := db_mysql.New("127.0.0.1", "3306", "root", "1234", "test")
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
	stmtNodes, _, err := tiparser.New().Parse("select seq, id from user where id = a limit 123 offset 456;", "", "")
	// stmtNodes, _, err := tiparser.New().Parse("select * from user where id = ?", "", "")

	require.NoError(t, err)
	for _, stmtNode := range stmtNodes {
		selectStmt := stmtNode.(*ast.SelectStmt)
		// from
		tableName := selectStmt.From.TableRefs.Left.(*ast.TableSource).Source.(*ast.TableName).Name.O
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
			continue
		}
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

		// where 필드 추출
		// fmt.Printf("%+v\n", selectStmt.Where)
		switch whereExpr := selectStmt.Where.(type) {
		case *ast.BinaryOperationExpr:
			binOpExpr := whereExpr
			switch leftExpr := binOpExpr.L.(type) {
			case *ast.ColumnNameExpr:
				col := leftExpr
				var colName, colType string
				colName = col.Name.Name.O
				colType, _ = sch.GetFieldType(tableName, colName)
				fmt.Printf("where | name : %s | db type : %s | golang Type : %s\n", colName, colType, myps.ConvType(colType))
			}
			/*
				switch rightExpr := binOpExpr.R.(type) {
				case *ast.ValuesExpr:
					valExpr := rightExpr
					fmt.Printf("where | value : %s | golang Type : %s\n", valExpr.Type.String(), myps.ConvType(valExpr.GetType()))
				}
			*/
		}
		// limit, offset 출력

		// 와.. 개에반데? ㅋㅋㅋ
		fmt.Printf(" ( %v ) : limit\n", selectStmt.Limit.Count.(*test_driver.ValueExpr).Datum.GetInt64())
		fmt.Printf(" ( %v ) : offset\n", selectStmt.Limit.Offset.(*test_driver.ValueExpr).Datum.GetInt64())

	}
}
