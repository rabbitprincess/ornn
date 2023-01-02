package atlas

import (
	"context"

	"ariga.io/atlas/schemahcl"
	"ariga.io/atlas/sql/migrate"
	"ariga.io/atlas/sql/mysql"
	"ariga.io/atlas/sql/postgres"
	"ariga.io/atlas/sql/schema"
	"ariga.io/atlas/sql/sqlite"
	"github.com/gokch/ornn/db"
	"github.com/hashicorp/hcl/v2/hclparse"
)

// global
func InspectSchema(dbType DbType, db *db.Conn) (*schema.Schema, error) {
	at := &Atlas{}
	err := at.Init(dbType, db)
	if err != nil {
		return nil, err
	}
	schema, err := at.InspectSchema()
	if err != nil {
		return nil, err
	}
	return schema, nil
}

type Atlas struct {
	Type        DbType
	marshaler   schemahcl.MarshalerFunc
	unmarshaler schemahcl.EvalFunc
	driver      migrate.Driver
}

func (t *Atlas) Init(dbType DbType, db *db.Conn) error {
	var err error

	t.Type = dbType
	switch dbType {
	case DbTypeMySQL, DbTypeMaria, DbTypeTiDB:
		t.marshaler = mysql.MarshalHCL
		t.unmarshaler = mysql.EvalHCL
		t.driver, err = mysql.Open(db.Db)
	case DbTypePostgre, DbTypeCockroachDB:
		t.marshaler = postgres.MarshalHCL
		t.unmarshaler = postgres.EvalHCL
		t.driver, err = postgres.Open(db.Db)
	case DbTypeSQLite:
		t.marshaler = sqlite.MarshalHCL
		t.unmarshaler = sqlite.EvalHCL
		t.driver, err = sqlite.Open(db.Db)
	}
	if err != nil {
		return err
	}
	return nil
}

func (t *Atlas) MarshalHCL(sch *schema.Schema) ([]byte, error) {
	bt, err := t.marshaler.MarshalSpec(sch)
	if err != nil {
		return nil, err
	}
	return bt, nil
}

func (t *Atlas) UnmarshalHCL(bt []byte) (*schema.Schema, error) {
	sch := schema.New("")
	parser := hclparse.NewParser()
	if _, diag := parser.ParseHCL(bt, ""); diag.HasErrors() {
		return nil, diag
	}
	err := t.unmarshaler.Eval(parser, sch, nil)
	if err != nil {
		return nil, err
	}
	return sch, nil
}

func (t *Atlas) InspectSchema() (*schema.Schema, error) {
	sch, err := t.driver.InspectSchema(context.Background(), "", nil)
	if err != nil {
		return nil, err
	}
	return sch, nil
}

func (t *Atlas) MigrateSchema(sch *schema.Schema) error {
	schemaCur, err := t.InspectSchema()
	if err != nil {
		return err
	}
	diffs, err := t.driver.SchemaDiff(schemaCur, sch)
	if err != nil {
		return err
	}
	return t.driver.ApplyChanges(context.Background(), diffs)
}
