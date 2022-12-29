package go_orm_gen

import (
	"log"

	"github.com/gokch/go-orm-gen/db"
	"github.com/gokch/go-orm-gen/db/db_mysql"
)

// 기본 파일명
const (
	DEF_S_default_filepath__json   = "bp.json"
	DEF_S_default_filepath__go__db = "bp_db.go"
)

func main() {
	db := &db.DB{}
	db.Connect("mysql", db_mysql.NewDsn("user", "pw", "127.0.0.1", "4001", "test_db"), "test_db")

	var orm *ORM = &ORM{}
	orm.Init(db)

	// config
	err := orm.ConfigLoad("bp.json")
	if err != nil {
		log.Fatal(err)
	}
	err = orm.SchemaLoad("")
	if err != nil {
		log.Fatal(err)
	}
	err = orm.ConfigSave("bp.json")
	if err != nil {
		log.Fatal(err)
	}

	// code - generate
	err = orm.GenCode("bp.json", map[string]string{
		DEF_s_gen_config__go__db__package_name: "gen",
		DEF_s_gen_config__go__db__class__name:  "gen",
	})
	if err != nil {
		log.Fatal(err)
	}
}
