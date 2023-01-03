package parser

import (
	"log"

	"github.com/blastrain/vitess-sqlparser/sqlparser"
)

type Parser interface {
}

func New(query string) (parser Parser, err error) {
	stmt, err := sqlparser.Parse(query)
	if err != nil {
		return nil, err
	}

	switch data := stmt.(type) {
	case *sqlparser.Select:
		query := &Select{}
		err = query.parse(data)
		parser = query
	case *sqlparser.Insert:
		query := &Insert{}
		err = query.parse(data)
		parser = query
	case *sqlparser.Update:
		query := &Update{}
		err = query.parse(data)
		parser = query
	case *sqlparser.Delete:
		query := &Delete{}
		err = query.parse(data)
		parser = query
	case *sqlparser.Union,
		*sqlparser.Set,
		*sqlparser.DDL,
		*sqlparser.Show,
		*sqlparser.Use,
		*sqlparser.OtherRead,
		*sqlparser.OtherAdmin,
		*sqlparser.TruncateTable:
		log.Fatalf("parser error | not support query statement %T", data)
	default:
		log.Fatalf("parser error | unknown query statement %T", data)
	}
	return parser, err
}

type Field struct {
	FieldName string
	TableName string
	Val       []byte
}

type FieldAs struct {
	Field string
	As    string
	Table string
}

type TableAs struct {
	Table string
	As    string
}

func (t *TableAs) get() (tableName string) {
	if t.As != "" {
		tableName = t.As
	} else {
		tableName = t.Table
	}

	return tableName
}
