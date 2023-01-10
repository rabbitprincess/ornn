package parser_sqlite

import (
	"fmt"

	"ariga.io/atlas/sql/schema"
	"github.com/CovenantSQL/sqlparser"
	"github.com/gokch/ornn/config"
	"github.com/gokch/ornn/parser"
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
	stmtNode, err := sqlparser.Parse(sql)
	if err != nil {
		return nil, err
	}
	parsedQuery := &parser.ParsedQuery{}
	parsedQuery.Init(sql)

	switch stmt := stmtNode.(type) {
	case *sqlparser.Select:
		err = p.parseSelect(stmt, parsedQuery)
	case *sqlparser.Insert:
		err = p.parseInsert(stmt, parsedQuery)
	case *sqlparser.Update:
		err = p.parseUpdate(stmt, parsedQuery)
	case *sqlparser.Delete:
		err = p.parseDelete(stmt, parsedQuery)
	default:
		err = fmt.Errorf("parser error | not support query statement %T", stmt)
	}
	if err != nil {
		return nil, err
	}

	return parsedQuery, nil
}

func (p *Parser) parseSelect(stmt *sqlparser.Select, parsedQuery *parser.ParsedQuery) error {
	parsedQuery.QueryType = parser.QueryTypeSelect

	// from
	if len(stmt.From) != 1 {
		// TODO: select join
		return fmt.Errorf("need more programming")
	}
	var tableName string
	var table *schema.Table
	switch from := stmt.From[0].(type) {
	case *sqlparser.AliasedTableExpr:
		tableName = from.Expr.(sqlparser.TableName).Name.String()
	default:
		panic("need more programming")
	}
	if table, _ = p.sch.Table(tableName); table == nil {
		return fmt.Errorf("table not found | %s", tableName)
	}

	// select
	for _, selectExpr := range stmt.SelectExprs {
		switch data := selectExpr.(type) {
		case *sqlparser.StarExpr:
			for _, col := range table.Columns {
				parsedQuery.Ret = append(parsedQuery.Ret, parser.NewField(col.Name, p.ConvType(col.Type)))
			}
			break
		case *sqlparser.AliasedExpr:
			switch data2 := data.Expr.(type) {
			case *sqlparser.ColName:
				colName := data2.Name.String()
				if col, _ := table.Column(colName); col != nil {
					parsedQuery.Ret = append(parsedQuery.Ret, parser.NewField(col.Name, p.ConvType(col.Type)))
				} else {
					parsedQuery.Ret = append(parsedQuery.Ret, parser.NewField(col.Name, "interface{}"))
				}
			default:
				panic("need more programming")
			}
		default:
			panic("need more programming")
		}
	}

	// where
	whereFields := ParseWhereToFields(stmt.Where.Expr)
	for _, where := range whereFields {
		if where.right == nil || where.left == nil {
			continue
		}
		// left 의 column 을 인자로 추출
		if paramMarkerExpr, _ := where.right.(*sqlparser.SQLVal); paramMarkerExpr != nil && paramMarkerExpr.Type == sqlparser.ValArg {
			if data, ok := where.left.(*sqlparser.ColName); ok == true {
				colName := data.Name.String()
				col, ok := table.Column(colName)
				if ok != true {
					parsedQuery.Arg = append(parsedQuery.Arg, parser.NewField(colName, "interface{}"))
				} else {
					parsedQuery.Arg = append(parsedQuery.Arg, parser.NewField(colName, p.ConvType(col.Type)))
				}
			}
		}
		// right 의 column 을 인자로 추출
		if paramMarkerExpr, _ := where.left.(*sqlparser.SQLVal); paramMarkerExpr != nil && paramMarkerExpr.Type == sqlparser.ValArg {
			if data, ok := where.right.(*sqlparser.ColName); ok == true {
				colName := data.Name.String()
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

func (p *Parser) parseInsert(stmt *sqlparser.Insert, parsedQuery *parser.ParsedQuery) error {
	parsedQuery.QueryType = parser.QueryTypeInsert
	// from
	var tableName string = stmt.Table.Name.String()
	var table *schema.Table
	if table, _ = p.sch.Table(tableName); table == nil {
		return fmt.Errorf("table not found | %s", tableName)
	}

	// values
	// insert fields
	vals, _ := stmt.Rows.(sqlparser.Values)
	if len(vals) != 1 {
		panic("bulk query is invalid, use bulk options")
	}
	colNames := make([]string, len(table.Columns))
	if len(stmt.Columns) == 0 { // insert all fields
		for i, col := range table.Columns {
			colNames[i] = col.Name
		}
		if len(table.Columns) != len(vals[0]) {
			panic("not same column and value count")
		}

		for i, list := range vals[0] {
			if _, paramMarkerExpr, ok := ParseDriverValue(list); !ok {
				panic("need more programming")
			} else if paramMarkerExpr != nil {
				parsedQuery.Arg = append(parsedQuery.Arg, parser.NewField(colNames[i], p.ConvType(table.Columns[i].Type)))
			}
		}
	} else { // insert specific fields
		if len(stmt.Columns) != len(vals[0]) {
			panic("not same column and value count")
		}
		for i, list := range vals[0] {
			if _, paramMarkerExpr, ok := ParseDriverValue(list); !ok {
				panic("need more programming")
			} else if paramMarkerExpr != nil {
				colName := stmt.Columns[i].String()
				col, ok := table.Column(colName)
				if ok != true {
					parsedQuery.Arg = append(parsedQuery.Arg, parser.NewField(colName, "interface{}"))
				} else {
					parsedQuery.Arg = append(parsedQuery.Arg, parser.NewField(colName, p.ConvType(col.Type)))
				}
			}
		}
	}
	// ondup
	if len(stmt.OnDup) != 0 {
		panic("need more programming")
	}

	return nil
}

func (p *Parser) parseUpdate(stmt *sqlparser.Update, parsedQuery *parser.ParsedQuery) error {
	parsedQuery.QueryType = parser.QueryTypeUpdate

	// from
	if len(stmt.TableExprs) != 1 {
		// TODO: select join
		return fmt.Errorf("need more programming")
	}
	var tableName string
	var table *schema.Table
	switch from := stmt.TableExprs[0].(type) {
	case *sqlparser.AliasedTableExpr:
		tableName = from.Expr.(sqlparser.TableName).Name.String()
	default:
		panic("need more programming")
	}
	if table, _ = p.sch.Table(tableName); table == nil {
		return fmt.Errorf("table not found | %s", tableName)
	}

	// select
	for _, updateExpr := range stmt.Exprs {
		switch data := updateExpr.Expr.(type) {
		case *sqlparser.SQLVal:
			if data.Type == sqlparser.ValArg {
				colName := updateExpr.Name.Name.String()
				if col, _ := table.Column(colName); col != nil {
					parsedQuery.Ret = append(parsedQuery.Ret, parser.NewField(col.Name, p.ConvType(col.Type)))
				} else {
					parsedQuery.Ret = append(parsedQuery.Ret, parser.NewField(col.Name, "interface{}"))
				}
			}
		default:
			panic("need more programming")
		}
	}

	// where
	whereFields := ParseWhereToFields(stmt.Where.Expr)
	for _, where := range whereFields {
		if where.right == nil || where.left == nil {
			continue
		}
		// left 의 column 을 인자로 추출
		if paramMarkerExpr, _ := where.right.(*sqlparser.SQLVal); paramMarkerExpr != nil && paramMarkerExpr.Type == sqlparser.ValArg {
			if data, ok := where.left.(*sqlparser.ColName); ok == true {
				colName := data.Name.String()
				col, ok := table.Column(colName)
				if ok != true {
					parsedQuery.Arg = append(parsedQuery.Arg, parser.NewField(colName, "interface{}"))
				} else {
					parsedQuery.Arg = append(parsedQuery.Arg, parser.NewField(colName, p.ConvType(col.Type)))
				}
			}
		}
		// right 의 column 을 인자로 추출
		if paramMarkerExpr, _ := where.left.(*sqlparser.SQLVal); paramMarkerExpr != nil && paramMarkerExpr.Type == sqlparser.ValArg {
			if data, ok := where.right.(*sqlparser.ColName); ok == true {
				colName := data.Name.String()
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

func (p *Parser) parseDelete(stmt *sqlparser.Delete, parsedQuery *parser.ParsedQuery) error {
	parsedQuery.QueryType = parser.QueryTypeDelete

	// from
	if len(stmt.TableExprs) != 1 {
		// TODO: select join
		return fmt.Errorf("need more programming")
	}
	var tableName string
	var table *schema.Table
	switch from := stmt.TableExprs[0].(type) {
	case *sqlparser.AliasedTableExpr:
		tableName = from.Expr.(sqlparser.TableName).Name.String()
	default:
		panic("need more programming")
	}
	if table, _ = p.sch.Table(tableName); table == nil {
		return fmt.Errorf("table not found | %s", tableName)
	}

	// where
	whereFields := ParseWhereToFields(stmt.Where.Expr)
	for _, where := range whereFields {
		if where.right == nil || where.left == nil {
			continue
		}
		// left 의 column 을 인자로 추출
		if paramMarkerExpr, _ := where.right.(*sqlparser.SQLVal); paramMarkerExpr != nil && paramMarkerExpr.Type == sqlparser.ValArg {
			if data, ok := where.left.(*sqlparser.ColName); ok == true {
				colName := data.Name.String()
				col, ok := table.Column(colName)
				if ok != true {
					parsedQuery.Arg = append(parsedQuery.Arg, parser.NewField(colName, "interface{}"))
				} else {
					parsedQuery.Arg = append(parsedQuery.Arg, parser.NewField(colName, p.ConvType(col.Type)))
				}
			}
		}
		// right 의 column 을 인자로 추출
		if paramMarkerExpr, _ := where.left.(*sqlparser.SQLVal); paramMarkerExpr != nil && paramMarkerExpr.Type == sqlparser.ValArg {
			if data, ok := where.right.(*sqlparser.ColName); ok == true {
				colName := data.Name.String()
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
