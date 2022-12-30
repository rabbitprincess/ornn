package config

import (
	"fmt"
)

//------------------------------------------------------------------------------------------------//
// schema

type Schema struct {
	Tables []*Table `json:"tables"`
}

func (t *Schema) AddTable(table *Table) {
	if t.Tables == nil {
		t.Tables = make([]*Table, 0, 10)
	}
	t.Tables = append(t.Tables, table)
}

func (t *Schema) GetTable(tableName string) *Table {
	if t.Tables == nil {
		t.Tables = make([]*Table, 0, 10)
	}

	for _, pt := range t.Tables {
		if tableName == pt.Name {
			return pt
		}
	}
	return nil
}

func (t *Schema) GetTableFieldMatched(fieldName string, tablesName []string) (matched []string, err error) {
	matched = make([]string, 0, 10)

	for _, tableName := range tablesName {
		table := t.GetTable(tableName)
		if table == nil {
			return nil, fmt.Errorf("wrong table name in sql query, table name is not exist in schema | table_name : %s", tableName)
		}
		if table.GetField(fieldName) != nil {
			matched = append(matched, tableName)
		}
	}
	return matched, nil
}

//------------------------------------------------------------------------------------------------//
// table

type Table struct {
	Name   string   `json:"table_name"`
	Indexs []*Index `json:"indexs"`
	Fields []*Field `json:"fields"`
}

func (t *Table) AddField(field *Field) {
	if t.Fields == nil {
		t.Fields = make([]*Field, 0, 10)
	}
	t.Fields = append(t.Fields, field)
}

func (t *Table) GetField(fieldName string) *Field {
	for _, field := range t.Fields {
		if fieldName == field.Name {
			return field
		}
	}
	return nil
}

func (t *Table) AddIndex(index *Index) {
	if t.Indexs == nil {
		t.Indexs = make([]*Index, 0, 10)
	}
	t.Indexs = append(t.Indexs, index)
}

func (t *Table) UpdateIndex(table *Table) error {
	t.Indexs = table.Indexs
	return nil
}

//------------------------------------------------------------------------------------------------//
// index

type Index struct {
	Name string   `json:"name"`
	Type string   `json:"type"`
	Keys []string `json:"keys"`

	Comment string `json:"comment,omitempty"`
}

func (t *Index) Set(name string, idxType string, keys []string) {
	t.Name = name
	t.Type = idxType
	t.Keys = keys
}

//------------------------------------------------------------------------------------------------//
// field

type Field struct {
	Name    string `json:"name"`
	TypeDB  string `json:"type_db"`
	TypeGen string `json:"type_gen"`

	Comment string `json:"comment,omitempty"`
}

func (t *Field) Set(name string, typeDB string, typeGen string) {
	t.Name = name
	t.TypeDB = typeDB
	t.TypeGen = typeGen
}
