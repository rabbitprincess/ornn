package sql

import (
	"log"

	parser "github.com/blastrain/vitess-sqlparser/sqlparser"
)

type SQL interface {
}

func New(query string) (sql SQL, err error) {
	stmt, err := parser.Parse(query)
	if err != nil {
		return nil, err
	}

	switch data := stmt.(type) {
	case *parser.Select:
		query := &Select{}
		err = query.parse(data)
		sql = query
	case *parser.Insert:
		query := &Insert{}
		err = query.parse(data)
		sql = query
	case *parser.Update:
		query := &Update{}
		err = query.parse(data)
		sql = query
	case *parser.Delete:
		query := &Delete{}
		err = query.parse(data)
		sql = query
	default:
		log.Fatal("parser error")
	}
	return sql, err
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
