package main

import (
	"log"

	"github.com/gokch/ornn/config"
	"github.com/gokch/ornn/db"
	"github.com/gokch/ornn/db/db_mysql"
	"github.com/gokch/ornn/ornn"
)

func main() {
	var err error

	db := &db.Conn{}
	db, err = db_mysql.New("127.0.0.1", "3306", "root", "1234", "myTestDatabase")
	if err != nil {
		log.Fatal(err)
	}

	var ornn *ornn.ORM = &ornn.ORM{}
	config := &config.Config{}
	config.Global.InitDefault()

	ornn.Init(db, config)

	// config
	err = ornn.ConfigLoad("./output/gen.json")
	if err != nil {
		log.Fatal(err)
	}

	err = ornn.SchemaLoad("")
	if err != nil {
		log.Fatal(err)
	}
	err = ornn.ConfigSave("./output/gen.json")
	if err != nil {
		log.Fatal(err)
	}

	// code - generate
	err = ornn.GenCode("./output/gen.go")
	if err != nil {
		log.Fatal(err)
	}
}
