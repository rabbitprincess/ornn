package parser_mysql

import (
	"fmt"

	"github.com/gokch/ornn/config"
	"github.com/gokch/ornn/parser"
	sqlparser "github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	_ "github.com/pingcap/tidb/parser/test_driver"
)

func New(sch *config.Schema) parser.Parser {
	return &Parser{
		sch: sch,
	}
}

type Parser struct {
	sch *config.Schema
}

// TODO
func (p *Parser) Parse(sql string) (*parser.ParsedQuery, error) {
	sqlParser := sqlparser.New()
	stmtNodes, _, err := sqlParser.Parse(sql, "", "")
	if err != nil {
		return nil, err
	}

	parseQuery := &parser.ParsedQuery{}
	parseQuery.Init()

	for _, stmtNode := range stmtNodes {
		switch stmt := stmtNode.(type) {
		case *ast.SelectStmt:

			// p.parseSelect(stmt, parseQuery)
		case *ast.InsertStmt:
			parseQuery.QueryType = parser.QueryTypeInsert
			_ = stmt
		case *ast.UpdateStmt:
			_ = stmt
		case *ast.DeleteStmt:
			_ = stmt
		default:
			return nil, fmt.Errorf("parser error | not support query statement %T", stmt)
		}
	}

	return parseQuery, nil
}

/*
// parseQuery 를 stmt 를 이용해 초기화, p.sch 를 이용해 타입 설정
func (p *Parser) parseSelect(stmt *ast.SelectStmt, parseQuery *parser.ParsedQuery) {
	parseQuery.QueryType = parser.QueryTypeSelect
	// select
	for _, field := range stmt.Fields.Fields {
		fmt.Println(field.Expr.GetType().String())
		switch fieldExpr := field.Expr.(type) {
		// case *ast:
		// parseQuery.SelectAll = true
		// return
		case *ast.ColumnNameExpr:
			colName := fieldExpr.Name
			col := p.sch.Table(colName.Table.Schema, colName.Table.Name).Column(colName.Name)
			if col == nil {
				panic(fmt.Sprintf("parser error | not found column %s.%s.%s", colName.Table.Schema, colName.Table.Name, colName.Name))
			}
			parseQuery.SelectColumns = append(parseQuery.SelectColumns, &parser.ParsedColumn{
				ColumnName: colName.Name,
				TableName:  colName.Table.Name,
				ColumnType: col.Type,
			})
		}
	}

	// where

	// single select 처리
	// 코드 생성 시 단일 구조체 반환 목적
	if stmt.Limit.Count.Text() == "1" {
		parseQuery.SelectSingle = true
	}

	return
}
*/
