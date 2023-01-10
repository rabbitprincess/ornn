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
			err = p.parseSelect(stmt, parseQuery)
		case *ast.InsertStmt:
			err = p.parseInsert(stmt, parseQuery)
		case *ast.UpdateStmt:
			err = p.parseUpdate(stmt, parseQuery)
		case *ast.DeleteStmt:
			err = p.parseDelete(stmt, parseQuery)
		default:
			err = fmt.Errorf("parser error | not support query statement %T", stmt)
		}
		if err != nil {
			return nil, err
		}
	}

	return parseQuery, nil
}

// parseQuery 를 stmt 를 이용해 초기화, p.sch 를 이용해 타입 설정
func (p *Parser) parseSelect(stmt *ast.SelectStmt, parsedQuery *parser.ParsedQuery) error {
	parsedQuery.QueryType = parser.QueryTypeSelect

	// from
	tableSources := ParseJoinToTables(stmt.From.TableRefs)
	if len(tableSources) != 1 {
		// TODO : select join
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
			parsedQuery.Ret = append(parsedQuery.Ret, parser.NewField(col.Name, p.ConvType(col.Type)))
		}
	} else {
		for _, field := range stmt.Fields.Fields {
			switch fieldExpr := field.Expr.(type) {
			case *ast.ColumnNameExpr:
				colName := fieldExpr.Name.Name.O
				col, ok := table.Column(colName)
				if ok != true {
					parsedQuery.Ret = append(parsedQuery.Ret, parser.NewField(colName, "interface{}"))
				} else {
					parsedQuery.Ret = append(parsedQuery.Ret, parser.NewField(colName, p.ConvType(col.Type)))
				}
			}
		}
	}
	// where
	whereFields := ParseWhereToFields(stmt.Where)
	for _, where := range whereFields {
		if where.right == nil || where.left == nil {
			continue
		}

		// left 의 column 을 인자로 추출
		if _, paramMarkerExpr, _ := ParseDriverValue(where.right); paramMarkerExpr != nil {
			if data, ok := where.left.(*ast.ColumnNameExpr); ok == true {
				colName := data.Name.Name.O
				col, ok := table.Column(colName)
				if ok != true {
					parsedQuery.Arg = append(parsedQuery.Arg, parser.NewField(colName, "interface{}"))
				} else {
					parsedQuery.Arg = append(parsedQuery.Arg, parser.NewField(colName, p.ConvType(col.Type)))
				}
			}
		}

		// right 의 column 을 인자로 추출
		if _, paramMarkerExpr, _ := ParseDriverValue(where.left); paramMarkerExpr != nil {
			if data, ok := where.right.(*ast.ColumnNameExpr); ok == true {
				colName := data.Name.Name.O
				col, ok := table.Column(colName)
				if ok != true {
					parsedQuery.Arg = append(parsedQuery.Arg, parser.NewField(colName, "interface{}"))
				} else {
					parsedQuery.Arg = append(parsedQuery.Arg, parser.NewField(colName, p.ConvType(col.Type)))
				}
			}
		}
	}

	return nil
}

// parseQuery 를 stmt 를 이용해 초기화, p.sch 를 이용해 타입 설정
func (p *Parser) parseInsert(stmt *ast.InsertStmt, parsedQuery *parser.ParsedQuery) error {
	parsedQuery.QueryType = parser.QueryTypeInsert

	// from
	tableSources := ParseJoinToTables(stmt.Table.TableRefs)
	if len(tableSources) != 1 {
		// TODO : insert join
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
				parsedQuery.Arg = append(parsedQuery.Arg, parser.NewField(colNames[i], p.ConvType(table.Columns[i].Type)))
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
				col, ok := table.Column(colName)
				if ok != true {
					parsedQuery.Arg = append(parsedQuery.Arg, parser.NewField(colName, "interface{}"))
				} else {
					parsedQuery.Arg = append(parsedQuery.Arg, parser.NewField(colName, p.ConvType(col.Type)))
				}
			}
		}
	}

	if stmt.Select != nil {
		// TODO : insert select
		panic("need more programming")
	}

	if stmt.OnDuplicate != nil {
		// TODO : on duplicate
		panic("need more programming")
	}
	return nil
}

func (p *Parser) parseUpdate(stmt *ast.UpdateStmt, parsedQuery *parser.ParsedQuery) error {
	parsedQuery.QueryType = parser.QueryTypeUpdate
	// update
	tableSources := ParseJoinToTables(stmt.TableRefs.TableRefs)
	if len(tableSources) != 1 {
		// TODO : update join
		panic("need more programming")
	}
	tableName := ParseTableName(tableSources[0])
	table, ok := p.sch.Table(tableName)
	if ok != true {
		return fmt.Errorf("parser error | not found table %s", tableName)
	}

	// set
	for _, set := range stmt.List {
		colName := set.Column.Name.O
		col, ok := table.Column(colName)
		if ok != true {
			parsedQuery.Arg = append(parsedQuery.Arg, parser.NewField(colName, "interface{}"))
		} else {
			parsedQuery.Arg = append(parsedQuery.Arg, parser.NewField(colName, p.ConvType(col.Type)))
		}
	}

	// where
	whereFields := ParseWhereToFields(stmt.Where)
	for _, where := range whereFields {
		if where.right == nil || where.left == nil {
			continue
		}

		// left 의 column 을 인자로 추출
		if _, paramMarkerExpr, _ := ParseDriverValue(where.right); paramMarkerExpr != nil {
			if data, ok := where.left.(*ast.ColumnNameExpr); ok == true {
				colName := data.Name.Name.O
				col, ok := table.Column(colName)
				if ok != true {
					parsedQuery.Arg = append(parsedQuery.Arg, parser.NewField("where_"+colName, "interface{}"))
				} else {
					parsedQuery.Arg = append(parsedQuery.Arg, parser.NewField("where_"+colName, p.ConvType(col.Type)))
				}
			}
		}

		// right 의 column 을 인자로 추출
		if _, paramMarkerExpr, _ := ParseDriverValue(where.left); paramMarkerExpr != nil {
			if data, ok := where.right.(*ast.ColumnNameExpr); ok == true {
				colName := data.Name.Name.O
				col, ok := table.Column(colName)
				if ok != true {
					parsedQuery.Arg = append(parsedQuery.Arg, parser.NewField("where_"+colName, "interface{}"))
				} else {
					parsedQuery.Arg = append(parsedQuery.Arg, parser.NewField("where_"+colName, p.ConvType(col.Type)))
				}
			}
		}
	}

	return nil
}

func (p *Parser) parseDelete(stmt *ast.DeleteStmt, parsedQuery *parser.ParsedQuery) error {
	parsedQuery.QueryType = parser.QueryTypeDelete

	// from
	tableSources := ParseJoinToTables(stmt.TableRefs.TableRefs)
	if len(tableSources) != 1 {
		// TODO : delete join
		panic("need more programming")
	}
	tableName := ParseTableName(tableSources[0])
	table, ok := p.sch.Table(tableName)
	if ok != true {
		return fmt.Errorf("parser error | not found table %s", tableName)
	}

	// where
	whereFields := ParseWhereToFields(stmt.Where)
	for _, where := range whereFields {
		if where.right == nil || where.left == nil {
			continue
		}
		// left 의 column 을 인자로 추출
		if _, paramMarkerExpr, _ := ParseDriverValue(where.right); paramMarkerExpr != nil {
			if data, ok := where.left.(*ast.ColumnNameExpr); ok == true {
				colName := data.Name.Name.O
				col, ok := table.Column(colName)
				if ok != true {
					parsedQuery.Arg = append(parsedQuery.Arg, parser.NewField(colName, "interface{}"))
				} else {
					parsedQuery.Arg = append(parsedQuery.Arg, parser.NewField(colName, p.ConvType(col.Type)))
				}
			}
		}

		// right 의 column 을 인자로 추출
		if _, paramMarkerExpr, _ := ParseDriverValue(where.left); paramMarkerExpr != nil {
			if data, ok := where.right.(*ast.ColumnNameExpr); ok == true {
				colName := data.Name.Name.O
				col, ok := table.Column(colName)
				if ok != true {
					parsedQuery.Arg = append(parsedQuery.Arg, parser.NewField(colName, "interface{}"))
				} else {
					parsedQuery.Arg = append(parsedQuery.Arg, parser.NewField(colName, p.ConvType(col.Type)))
				}
			}
		}
	}

	return nil
}
