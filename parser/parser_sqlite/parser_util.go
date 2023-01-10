package parser_sqlite

import (
	"github.com/CovenantSQL/sqlparser"
)

type binaryExpr struct {
	left  sqlparser.Expr
	right sqlparser.Expr
	op    string
}

func ParseWhereToFields(whereExpr sqlparser.Expr) []*binaryExpr {
	if whereExpr == nil {
		return nil
	}
	fields := make([]*binaryExpr, 0, 100)

	switch data := whereExpr.(type) {
	case *sqlparser.BinaryExpr:
		fields = append(fields, &binaryExpr{
			left:  data.Left,
			right: data.Right,
			op:    data.Operator,
		})
	case *sqlparser.AndExpr:
		fields = append(ParseWhereToFields(data.Left), fields...)
		fields = append(fields, ParseWhereToFields(data.Right)...)
	case *sqlparser.OrExpr:
		fields = append(ParseWhereToFields(data.Left), fields...)
		fields = append(fields, ParseWhereToFields(data.Right)...)
	case *sqlparser.ComparisonExpr:
		fields = append(fields, &binaryExpr{
			left:  data.Left,
			right: data.Right,
		})
	case *sqlparser.ParenExpr:
		fields = append(fields, ParseWhereToFields(data.Expr)...)
	case *sqlparser.NotExpr:
		fields = append(fields, ParseWhereToFields(data.Expr)...)
	case *sqlparser.ExistsExpr:
		fields = append(fields, ParseWhereToFields(data.Subquery)...)
	case *sqlparser.SQLVal:
		// do nothing
	case *sqlparser.NullVal:
		// do nothing
	case *sqlparser.ColName:
		// do nothing
	case *sqlparser.Subquery:
		// do nothing
	case *sqlparser.ListArg:
		// do nothing
	default:
		panic("parser error | not support where type")
	}
	return fields
}

func ParseDriverValue(node sqlparser.Expr) (*sqlparser.ColName, *sqlparser.SQLVal, bool) {
	switch data := node.(type) {
	case *sqlparser.SQLVal:
		return nil, data, true
	default:
		return nil, nil, false
	}
}
