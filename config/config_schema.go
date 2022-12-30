package config

import (
	"fmt"
	"strings"
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

func (t *Schema) UpdateTable(schema *Schema, tablePrefix string) error {
	if t.Tables == nil {
		t.Tables = make([]*Table, 0, 10)
	}

	tablesNew := make([]*Table, 0, len(schema.Tables))

	for _, table1 := range schema.Tables {
		var exist bool
		// 1. 중복 테이블은 업데이트
		for _, table2 := range t.Tables {
			if table2.Name == table1.Name {
				// 중복 인덱스 업데이트
				table2.UpdateIndex(table1)
				// 중복 필드 업데이트
				table2.UpdateField(table1)
				tablesNew = append(tablesNew, table2)

				exist = true
				break
			}
		}
		// 2. 새로운 테이블 추가
		if exist == false {
			tablesNew = append(tablesNew, table1)
		}
		// 3. 기존 테이블은 추가하지 않음 ( 삭제 )
	}

	// prefix 가 있을 시 후처리 - 해당 prefix 를 가지고 있는 테이블만 생성
	if tablePrefix != "" {
		tableNewWithPrefix := make([]*Table, 0, len(schema.Tables))
		prefixs := strings.Split(tablePrefix, ",")
		for _, tableNew := range tablesNew {
			for _, prefix := range prefixs {
				// prefix 중 하나랑 매칭될 시
				if strings.HasPrefix(tableNew.Name, prefix) == true {
					tableNewWithPrefix = append(tableNewWithPrefix, tableNew)
					break
				}
			}
		}
		tablesNew = tableNewWithPrefix
	}

	t.Tables = tablesNew
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

func (t *Table) Init(tableName string) {
	t.Name = tableName
	t.Fields = make([]*Field, 0, 10)
	t.Indexs = make([]*Index, 0, 10)
}

func (t *Table) AddField(field *Field) {
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

func (t *Table) UpdateField(table *Table) error {
	fieldNew := make([]*Field, 0, len(table.Fields))
	for _, field1 := range table.Fields {
		var exist bool
		for _, field2 := range t.Fields {
			// 1. 중복 필드는 업데이트
			if field2.Name == field1.Name {
				if field2.TypeDB == field1.TypeDB && field2.TypeGen != field1.TypeGen {
					// db type 은 같은데 bp type 만 다르면 업데이트 하지 않음
					fieldNew = append(fieldNew, field2)
				} else {
					// 나머지 케이스에서는 업데이트
					fieldNew = append(fieldNew, field1)
				}
				exist = true
				break
			}
		}
		// 2. 새로운 필드는 추가
		if exist == false {
			fieldNew = append(fieldNew, field1)
		}
		// 3. 기존 필드는 추가하지 않음
	}
	t.Fields = fieldNew
	return nil
}

func (t *Table) AddIndex(index *Index) {
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
