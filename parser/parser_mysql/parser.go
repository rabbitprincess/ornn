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

func (p *Parser) Parse(sql string) (*parser.ParsedQuery, error) {
	sqlParser := sqlparser.New()
	stmtNodes, _, err := sqlParser.Parse(sql, "", "")
	if err != nil {
		return nil, err
	}

	parseQuery := &parser.ParsedQuery{}
	parseQuery.Init(sql)

	for _, stmtNode := range stmtNodes {
		switch stmt := stmtNode.(type) {
		case *ast.SelectStmt:
			p.parseSelect(stmt, parseQuery)
		case *ast.InsertStmt:
			p.parseInsert(stmt, parseQuery)
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
	tableSources := ParseJoinToTables(stmt.From.TableRefs)
	if len(tableSources) != 1 {
		// todo - select join
		panic("need more programming")
	}
	tableName := ParseTableName(tableSources[0])
	table, ok := p.sch.Table(tableName)
	if ok != true {
		return fmt.Errorf("parser error | not found table %s", tableName)
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
	whereFields := ParseWhereToFields(stmt.Where)
	for left, right := range whereFields {
		// left 의 column 을 인자로 추출
		if _, paramMarkerExpr, _ := ParseDriverValue(right); paramMarkerExpr != nil {
			if data, ok := left.(*ast.ColumnNameExpr); ok == true {
				colName := data.Name.Name.O
				col, ok := table.Column(colName)
				if ok != true {
					parseQuery.Arg[col.Name] = "interface{}"
				}
				parseQuery.Arg[col.Name] = p.ConvType(col.Type.Raw)
			}
		}

		// right 의 column 을 인자로 추출
		if _, paramMarkerExpr, _ := ParseDriverValue(left); paramMarkerExpr != nil {
			if data, ok := right.(*ast.ColumnNameExpr); ok == true {
				colName := data.Name.Name.O
				col, ok := table.Column(colName)
				if ok != true {
					parseQuery.Arg[col.Name] = "interface{}"
				}
				parseQuery.Arg[col.Name] = p.ConvType(col.Type.Raw)
			}
		}
	}

	return nil
}

// parseQuery 를 stmt 를 이용해 초기화, p.sch 를 이용해 타입 설정
func (p *Parser) parseInsert(stmt *ast.InsertStmt, parseQuery *parser.ParsedQuery) error {
	parseQuery.QueryType = parser.QueryTypeInsert

	// from
	tableSources := ParseJoinToTables(stmt.Table.TableRefs)
	if len(tableSources) != 1 {
		// todo - insert join
		panic("need more programming")
	}
	tableName := ParseTableName(tableSources[0])
	table, ok := p.sch.Table(tableName)
	if ok != true {
		return fmt.Errorf("parser error | not found table %s", tableName)
	}

	// insert fields
	if len(stmt.Lists) != 1 {
		panic("bulk query is invalid, use bulk options")
	}
	colNames := make([]string, len(table.Columns))
	if len(stmt.Columns) == 0 { // insert all fields
		for i, col := range table.Columns {
			colNames[i] = col.Name
		}
		if len(table.Columns) != len(stmt.Lists[0]) {
			panic("not same column and value count")
		}

		for i, list := range stmt.Lists[0] {
			if _, paramMarkerExpr, ok := ParseDriverValue(list); !ok {
				panic("need more programming")
			} else if paramMarkerExpr != nil {
				parseQuery.Arg[table.Columns[i].Name] = p.ConvType(table.Columns[i].Type.Raw)
			}
		}
	} else { // insert specific fields
		if len(stmt.Columns) != len(stmt.Lists[0]) {
			panic("not same column and value count")
		}
		for i, list := range stmt.Lists[0] {
			if _, paramMarkerExpr, ok := ParseDriverValue(list); !ok {
				panic("need more programming")
			} else if paramMarkerExpr != nil {
				colName := stmt.Columns[i].Name.O
				col, exist := table.Column(colName)
				if exist != true {
					panic("not found column")
				}
				parseQuery.Arg[colName] = p.ConvType(col.Type.Raw)
			}
		}
	}

	if stmt.Select != nil {
		// todo - insert select
		panic("need more programming")
	}

	if stmt.OnDuplicate != nil {
		// todo - on duplicate
		panic("need more programming")
	}
	return nil
}
