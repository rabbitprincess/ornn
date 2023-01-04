package parser_mysql

import (
	"fmt"
	"testing"

	"github.com/gokch/ornn/atlas"
	"github.com/gokch/ornn/config"
	"github.com/gokch/ornn/db/db_mysql"
	tiparser "github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
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

	ps := tiparser.New()
	stmtNodes, _, err := ps.Parse("select seq,id,name from user where id = ?", "", "")
	// stmtNodes, _, err := ps.Parse("select * from user where id = ?", "", "")
	require.NoError(t, err)
	for _, stmtNode := range stmtNodes {
		selectStmt := stmtNode.(*ast.SelectStmt)
		// from
		tableName := selectStmt.From.TableRefs.Left.(*ast.TableSource).Source.(*ast.TableName).Name.O
		fmt.Println(tableName)
		// select
		// select * 일 경우 schema 의 모든 필드 추출
		// fmt.Printf("fields = %+v\n", selectStmt.TableHints)
		for _, field := range selectStmt.Fields.Fields {
			// fmt.Println(field.AsName)
			// fmt.Println("wildcard : ", field.WildCard)
			// fmt.Printf("%+v\n", field)
			switch fieldExpr := field.Expr.(type) {
			case *ast.ColumnNameExpr:
				col := fieldExpr
				var colName, colType string
				colName = col.Name.Name.O
				colType, _ = sch.GetFieldType(tableName, colName)
				fmt.Printf("name : %s | type : %s\n", colName, colType)
			}
		}

		// where
	}
}
