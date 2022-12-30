package ornn

import (
	"fmt"
	"os"

	"github.com/gokch/ornn/config"
	"github.com/gokch/ornn/db"
	"github.com/gokch/ornn/db/db_mysql"
)

type ORNN struct {
	db     *db.Conn
	vendor *Vendor
	config *config.Config
}

func (t *ORNN) Init(db *db.Conn, config *config.Config) {
	t.db = db
	t.vendor = NewVendor(db_mysql.NewVendor(db))
	t.config = config
}

func (t *ORNN) ConfigLoad(path string) error {
	if t.config == nil {
		return fmt.Errorf("json is emtpy")
	}

	t.config.Load(path)

	return nil
}

func (t *ORNN) ConfigSave(path string) error {
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

func (t *ORNN) InitConfigBySchema(tablePrefix string) error {
	// TODO : prefix 가 있을 시 해당 prefix 를 가지고 있는 테이블 스키마만 생성
	// TODO : custom type 은 업데이트 되지 않음
	// init config by schema
	err := t.vendor.VendorBySchema(t.config)
	if err != nil {
		return err
	}
	return nil
}

func (t *ORNN) GenCode(path string) (err error) {
	if t.config == nil {
		return fmt.Errorf("config is emtpy")
	}

	// gen code
	gen := &Gen{}
	code, err := gen.Gen(t.db, t.config, path)
	if err != nil {
		return err
	}

	// write code to file
	err = os.WriteFile(path, []byte(code), 0700)
	if err != nil {
		return err
	}
	return nil
}
