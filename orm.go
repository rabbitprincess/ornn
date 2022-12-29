package go_orm_gen

import (
	"fmt"

	"github.com/gokch/go-orm-gen/config"
	"github.com/gokch/go-orm-gen/db"
	"github.com/gokch/go-orm-gen/db/db_mysql"
)

type ORM struct {
	db     *db.DB
	config *config.Config
}

func (t *ORM) Init(db *db.DB) {
	t.db = db
}

func (t *ORM) ConfigLoad(path string) error {
	var err error

	t.config = &config.Config{}
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

func (t *ORM) SchemaLoad(talePrefix string) error {
	vendor := &DbVendor{}
	vendor.Init(db_mysql.NewVendor(t.db))

	schema, err := vendor.SchemaGet()
	if err != nil {
		return err
	}
	// 1. arrs_table_name__prefix 가 존재할 시 해당 prefix 를 가지고 있는 테이블 스키마만 생성
	// 2. custom type 은 업데이트 되지 않음
	err = t.config.Schema.UpdateTable(schema, talePrefix)
	if err != nil {
		return err
	}
	return nil
}

func (t *ORM) GenCode(path string, config map[string]string) (err error) {
	if config == nil {
		return fmt.Errorf("config is emtpy")
	}
	if t.config == nil {
		return fmt.Errorf("json is emtpy")
	}

	gen := &Gen{}
	_, err = gen.Gen(t.db, t.config, config, path)
	if err != nil {
		return err
	}

	return nil
}
