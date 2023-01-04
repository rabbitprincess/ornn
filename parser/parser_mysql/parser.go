package parser_mysql

import (
	"fmt"

	"ariga.io/atlas/sql/schema"
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
			p.parseSelect(stmt, parseQuery)
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

// parseQuery 를 stmt 를 이용해 초기화, p.sch 를 이용해 타입 설정
func (p *Parser) parseSelect(stmt *ast.SelectStmt, parseQuery *parser.ParsedQuery) error {
	parseQuery.QueryType = parser.QueryTypeSelect

	// from ( TODO : join )
	var table *schema.Table
	if tbl, ok := stmt.From.TableRefs.Left.(*ast.TableSource); ok {
		if tblName, ok := tbl.Source.(*ast.TableName); ok {
			table, ok = p.sch.Table(tblName.Name.O)
			if ok != true {
				// need more programming
				return fmt.Errorf("parser error | not found table %s", tblName.Name.O)
			}
		}
	}

	// select
	fields := stmt.Fields.Fields
	if len(fields) == 1 && fields[0].WildCard != nil { // select * 일 경우 schema 의 모든 필드 추출
		for _, col := range table.Columns {
			parseQuery.Ret[col.Name] = p.ConvType(col.Type.Raw)
		}
	} else {
		for _, field := range stmt.Fields.Fields {
			switch fieldExpr := field.Expr.(type) {
			case *ast.ColumnNameExpr:
				colName := fieldExpr.Name
				col, ok := table.Column(fieldExpr.Name.Name.O)
				if ok != true {
					panic(fmt.Sprintf("parser error | not found column %s.%s", table.Name, colName.Name))
				}
				parseQuery.Ret[col.Name] = p.ConvType(col.Type.Raw)
			}
		}
	}
	// where
	p.parseSelectWhere(stmt.Where, table, parseQuery)

	return nil
}

// set parseQuery rets recursive where expr
func (p *Parser) parseSelectWhere(whereExpr ast.ExprNode, table *schema.Table, parseQuery *parser.ParsedQuery) {
	switch whereExpr := whereExpr.(type) {
	case *ast.BinaryOperationExpr:
		// left
		switch leftExpr := whereExpr.L.(type) {
		case *ast.ColumnNameExpr:
			col := leftExpr
			colName := col.Name.Name.O
			colType, _ := p.sch.GetFieldType(table.Name, colName)
			parseQuery.Arg[colName] = p.ConvType(colType)
		}
		// right
		switch rightExpr := whereExpr.R.(type) {
		case *ast.ColumnNameExpr:
			col := rightExpr
			colName := col.Name.Name.O
			colType, _ := p.sch.GetFieldType(table.Name, colName)
			parseQuery.Arg[colName] = p.ConvType(colType)
		}
		// recursive
		p.parseSelectWhere(whereExpr.L, table, parseQuery)
		p.parseSelectWhere(whereExpr.R, table, parseQuery)
	}
}
