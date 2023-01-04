package parser_mysql

import (
	"fmt"

	"github.com/gokch/ornn/config"
	"github.com/gokch/ornn/parser"
	sqlparser "github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	_ "github.com/pingcap/tidb/parser/test_driver"
)

// TODO
type ParserMysql struct {
	sch *config.Schema
}

func (p *ParserMysql) Init(sch *config.Schema) {
	p.sch = sch
}

func (p *ParserMysql) Parse(sql string) (*parser.ParsedQuery, error) {
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

func (p *ParserMysql) parseSelect(stmt *ast.SelectStmt, parseQuery *parser.ParsedQuery) {
	// select
	for _, field := range stmt.Fields.Fields {
		var fieldName, fieldType string
		fieldName = field.AsName.O
		fmt.Println(field.AsName.O)
		fmt.Println(field.Expr.GetType().String())

		parseQuery.Ret[fieldName] = fieldType
	}

	// where

	// single select 처리
	// 코드 생성 시 단일 구조체 반환 목적
	if stmt.Limit.Count.Text() == "1" {
		parseQuery.SelectSingle = true
	}
	return
}
