package config

import (
	"context"

	"ariga.io/atlas/sql/mysql"
	"ariga.io/atlas/sql/schema"
	"github.com/gokch/ornn/db"
)

type Schema struct {
	*schema.Schema `json:"-"`
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
