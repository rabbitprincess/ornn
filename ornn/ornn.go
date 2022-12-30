package ornn

import (
	"fmt"

	"github.com/gokch/ornn/config"
	"github.com/gokch/ornn/db"
	"github.com/gokch/ornn/db/db_mysql"
)

type ORM struct {
	db     *db.Conn
	vendor *DbVendor
	config *config.Config
}

func (t *ORM) Init(db *db.Conn, config *config.Config) {
	t.db = db
	t.vendor = NewVendor(db_mysql.NewVendor(db))
	t.config = config
}

func (t *ORM) ConfigLoad(path string) error {
	var err error
	if t.config == nil {
		return fmt.Errorf("json is emtpy")
	}

	err = t.config.Load(path)
	if err != nil {
		return err
	}

	return nil
}

func (t *ORM) ConfigSave(path string) error {
	var err error
	if t.config == nil {
		return fmt.Errorf("json is emtpy")
	}

	err = t.config.Save(path)
	if err != nil {
		return err
	}

	return nil
}

func (t *ORM) SchemaLoad(tablePrefix string) error {
	schema, err := t.vendor.SchemaGet()
	if err != nil {
		return err
	}
	// 1. arrs_table_name__prefix 가 존재할 시 해당 prefix 를 가지고 있는 테이블 스키마만 생성
	// 2. custom type 은 업데이트 되지 않음
	err = t.config.Schema.UpdateTable(schema, tablePrefix)
	if err != nil {
		return err
	}
	return nil
}

func (t *ORM) GenCode(path string) (err error) {
	if t.config == nil {
		return fmt.Errorf("config is emtpy")
	}

	gen := &Gen{}
	_, err = gen.Gen(t.db, t.config, path)
	if err != nil {
		return err
	}

	return nil
}
