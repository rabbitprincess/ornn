package config

import (
	"context"
	"encoding/json"
	"fmt"

	"ariga.io/atlas/sql/mysql"
	"ariga.io/atlas/sql/schema"
	"github.com/gokch/ornn/db"
)

type Schema struct {
	*schema.Schema
}

func (t *Schema) InitDefault(db *db.Conn) error {
	var err error

	migrate, err := mysql.Open(db.Db)
	if err != nil {
		return err
	}
	t.Schema, err = migrate.InspectSchema(context.Background(), "", nil)
	if err != nil {
		return err
	}
	return nil
}

func (t *Schema) GetTableFieldMatched(fieldName string, tablesName []string) (matched []string, err error) {
	matched = make([]string, 0, 10)

	for _, tableName := range tablesName {
		table, exist := t.Table(tableName)
		if exist != true {
			return nil, fmt.Errorf("wrong table name in sql query, table name is not exist in schema | table_name : %s", tableName)
		}
		_, exist = table.Column(fieldName)
		if exist == true {
			matched = append(matched, tableName)
		}
	}
	return matched, nil
}

func (t *Schema) UnmarshalJSON(data []byte) error {
	sch := struct {
		Name   string
		Tables []*schema.Table
		Attrs  []schema.Attr // Attrs and options.
	}{}
	err := json.Unmarshal(data, &sch)
	if err != nil {
		return err
	}
	t.Schema = &schema.Schema{
		Name:   sch.Name,
		Tables: sch.Tables,
		Attrs:  sch.Attrs,
	}
	return nil
}

// 임시 - attr 지원 안됨.. 그냥 깔끔하게 hcl 로 갈까...
// MarshalJSON implements json.Marshaler.
func (s *Schema) MarshalJSON() ([]byte, error) {
	type (
		Column struct {
			Name string `json:"name"`
			Type string `json:"type,omitempty"`
			Null bool   `json:"null,omitempty"`
		}
		IndexPart struct {
			Desc   bool   `json:"desc,omitempty"`
			Column string `json:"column,omitempty"`
			Expr   string `json:"expr,omitempty"`
		}
		Index struct {
			Name   string      `json:"name,omitempty"`
			Unique bool        `json:"unique,omitempty"`
			Parts  []IndexPart `json:"parts,omitempty"`
		}
		ForeignKey struct {
			Name       string   `json:"name"`
			Columns    []string `json:"columns,omitempty"`
			References struct {
				Table   string   `json:"table"`
				Columns []string `json:"columns,omitempty"`
			} `json:"references"`
		}
		Table struct {
			Name        string       `json:"name"`
			Columns     []Column     `json:"columns,omitempty"`
			Indexes     []Index      `json:"indexes,omitempty"`
			PrimaryKey  *Index       `json:"primary_key,omitempty"`
			ForeignKeys []ForeignKey `json:"foreign_keys,omitempty"`
		}
		Schema struct {
			Name   string  `json:"name"`
			Tables []Table `json:"tables,omitempty"`
		}
	)

	s2 := Schema{Name: s.Name}
	for _, t1 := range s.Tables {
		t2 := Table{Name: t1.Name}
		for _, c1 := range t1.Columns {
			t2.Columns = append(t2.Columns, Column{
				Name: c1.Name,
				Type: c1.Type.Raw,
				Null: c1.Type.Null,
			})
		}
		idxParts := func(idx *schema.Index) (parts []IndexPart) {
			for _, p1 := range idx.Parts {
				p2 := IndexPart{Desc: p1.Desc}
				switch {
				case p1.C != nil:
					p2.Column = p1.C.Name
				case p1.X != nil:
					switch t := p1.X.(type) {
					case *schema.Literal:
						p2.Expr = t.V
					case *schema.RawExpr:
						p2.Expr = t.X
					}
				}
				parts = append(parts, p2)
			}
			return parts
		}
		for _, idx1 := range t1.Indexes {
			t2.Indexes = append(t2.Indexes, Index{
				Name:   idx1.Name,
				Unique: idx1.Unique,
				Parts:  idxParts(idx1),
			})
		}
		if t1.PrimaryKey != nil {
			t2.PrimaryKey = &Index{Parts: idxParts(t1.PrimaryKey)}
		}
		for _, fk1 := range t1.ForeignKeys {
			fk2 := ForeignKey{Name: fk1.Symbol}
			for _, c1 := range fk1.Columns {
				fk2.Columns = append(fk2.Columns, c1.Name)
			}
			fk2.References.Table = fk1.RefTable.Name
			for _, c1 := range fk1.RefColumns {
				fk2.References.Columns = append(fk2.References.Columns, c1.Name)
			}
			t2.ForeignKeys = append(t2.ForeignKeys, fk2)
		}
		s2.Tables = append(s2.Tables, t2)
	}
	return json.Marshal(s2)
}
