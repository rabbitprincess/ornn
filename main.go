package main

import (
	"log"

	"github.com/gokch/ornn/db"
	"github.com/gokch/ornn/db/db_mysql"
	"github.com/gokch/ornn/orm"
)

func main() {
	var err error

	db := &db.Conn{}
	db, err = db_mysql.New("127.0.0.1", "3306", "root", "951753ck", "myTestDatabase")
	if err != nil {
		log.Fatal(err)
	}

	var orm *orm.ORM = &orm.ORM{}
	orm.Init(db)

	// config
	err = orm.ConfigLoad("./output/gen.json")
	if err != nil {
		log.Fatal(err)
	}

	err = orm.SchemaLoad("")
	if err != nil {
		log.Fatal(err)
	}
	err = orm.ConfigSave("./output/gen.json")
	if err != nil {
		log.Fatal(err)
	}

	// code - generate
	err = orm.GenCode("./output/gen.go", map[string]string{})
	if err != nil {
		log.Fatal(err)
	}
}
