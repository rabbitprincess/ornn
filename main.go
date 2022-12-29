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

func main(dbAddr, dbPort, dbId, dbPw, dbName string, tableName string, configPath, filePath, packageName string, className string) {

	db := &db.DB{}
	db.Connect("mysql", db_mysql.NewDsn(dbId, dbPw, dbAddr, dbPort, dbName), dbName)

	var orm *ORM = &ORM{}
	orm.Init(db)

	// config
	err := orm.ConfigLoad(configPath)
	if err != nil {
		log.Fatal(err)
	}
	err = orm.SchemaLoad(tableName)
	if err != nil {
		log.Fatal(err)
	}
	err = orm.ConfigSave(configPath)
	if err != nil {
		log.Fatal(err)
	}

	// code - generate
	err = orm.GenCode(filePath, map[string]string{
		DEF_s_gen_config__go__db__package_name: packageName,
		DEF_s_gen_config__go__db__class__name:  className,
	})
	if err != nil {
		log.Fatal(err)
	}
}
