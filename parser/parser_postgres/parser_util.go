package parser_postgres

import (
	"github.com/auxten/postgresql-parser/pkg/sql/sem/tree"
)

func ParseDriverValue(node tree.Expr) (*tree.NumVal, *tree.StrVal, *tree.Placeholder, bool) {
	switch data := node.(type) {
	case *tree.NumVal:
		return data, nil, nil, true
	case *tree.StrVal:
		return nil, data, nil, true
	case *tree.Placeholder:
		return nil, nil, data, true
	default:
		return nil, nil, nil, false
	}
}

type binaryExpr struct {
	left  tree.Expr
	right tree.Expr
	op    string
}

func ParseWhereToFields(whereExpr tree.Expr) []*binaryExpr {
	if whereExpr == nil {
		return nil
	}
	fields := make([]*binaryExpr, 0, 100)

	switch data := whereExpr.(type) {
	case *tree.ComparisonExpr:
		fields = append(fields, &binaryExpr{
			left:  data.Left,
			right: data.Right,
			op:    data.Operator.String(),
		})
	case *tree.AndExpr:
		fields = append(ParseWhereToFields(data.Left), fields...)
		fields = append(fields, ParseWhereToFields(data.Right)...)
	case *tree.OrExpr:
		fields = append(ParseWhereToFields(data.Left), fields...)
		fields = append(fields, ParseWhereToFields(data.Right)...)
	case *tree.ParenExpr:
		fields = append(fields, ParseWhereToFields(data.Expr)...)
	case *tree.NotExpr:
		fields = append(fields, ParseWhereToFields(data.Expr)...)
	case *tree.NumVal:
		// do nothing
	case *tree.StrVal:
		// do nothing
	case *tree.Placeholder:
		// do nothing
	case *tree.Subquery:
		// do nothing
	default:
		panic("parser error | not support where type")
	}
	return fields
}
