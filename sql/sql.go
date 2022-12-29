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
		log.Fatal("bp", "parser error")
	}
	return sql, err
}

type Field struct {
	FldName string
	TblName string
	Val     []byte
}

type FieldAs struct {
	Fld string
	As  string
	Tbl string
}

func (t *FieldAs) get() (fieldName string) {
	if t.Tbl != "" {
		fieldName = t.Tbl + "__"
	}

	if t.As != "" {
		fieldName += t.As
	} else {
		fieldName += t.Fld
	}

	return fieldName
}

type TableAs struct {
	Tbl string
	As  string
}

func (t *TableAs) get() (tblName string) {
	if t.As != "" {
		tblName = t.As
	} else {
		tblName = t.Tbl
	}

	return tblName
}
