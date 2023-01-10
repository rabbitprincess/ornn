package parser_mysql

import (
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/opcode"
	"github.com/pingcap/tidb/parser/test_driver"
)

// get list recursively left, right
func ParseJoinToTables(join *ast.Join) []*ast.TableSource {
	if join == nil {
		return nil
	}
	nodes := make([]*ast.TableSource, 0, 10)
	if join.Left != nil {
		switch data := join.Left.(type) {
		case *ast.Join:
			nodes = append(nodes, ParseJoinToTables(data)...)
		case *ast.TableSource:
			nodes = append(nodes, data)
		default:
			panic("parser error | not support join type")
		}
	}
	if join.Right != nil {
		switch data := join.Right.(type) {
		case *ast.Join:
			nodes = append(nodes, ParseJoinToTables(data)...)
		case *ast.TableSource:
			nodes = append(nodes, data)
		default:
			panic("parser error | not support join type")
		}
	}
	return nodes
}

func ParseTableName(table *ast.TableSource) string {
	switch data := table.Source.(type) {
	case *ast.TableName:
		return data.Name.String()
	case *ast.SelectStmt:
		return data.Text()
	default:
		panic("parser error | not support table type")
	}
}

type binaryExpr struct {
	left  ast.ExprNode
	op    opcode.Op
	right ast.ExprNode
}

// set column and value to map recursively
func ParseWhereToFields(where ast.ExprNode) []*binaryExpr {
	if where == nil {
		return nil
	}
	fields := make([]*binaryExpr, 0, 10)
	switch data := where.(type) {
	case *ast.BinaryOperationExpr:
		switch data.Op {
		case opcode.LogicAnd, opcode.LogicOr, opcode.LogicXor:
			fields = append(ParseWhereToFields(data.L), fields...)
			fields = append(fields, ParseWhereToFields(data.R)...)
		case opcode.EQ, opcode.GE, opcode.GT, opcode.LE, opcode.LT, opcode.NE:
			fields = append(fields, &binaryExpr{
				left:  data.L,
				op:    data.Op,
				right: data.R,
			})
		case opcode.Like:
			fields = append(fields, &binaryExpr{
				left:  data.L,
				op:    data.Op,
				right: data.R,
			})
		default:
			panic("parser error | not support where type")
		}
	case *ast.ColumnNameExpr:
		fields = append(fields, &binaryExpr{
			left: data,
		})
	default:
		panic("parser error | not support where type")
	}
	return fields
}

func ParseDriverValue(node ast.ExprNode) (*test_driver.ValueExpr, *test_driver.ParamMarkerExpr, bool) {
	switch data := node.(type) {
	case *test_driver.ValueExpr:
		return data, nil, true
	case *test_driver.ParamMarkerExpr:
		return nil, data, true
	default:
		return nil, nil, false
	}
}
