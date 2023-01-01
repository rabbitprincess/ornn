package main

import (
	"log"

	"github.com/gokch/ornn/config"
	"github.com/gokch/ornn/db/db_mysql"
	"github.com/gokch/ornn/ornn"
)

func main() {
	db, err := db_mysql.New("127.0.0.1", "3306", "root", "1234", "test")
	if err != nil {
		log.Fatal(err)
	}

	config := &config.Config{}
	config.Global.InitDefault()
	if err != nil {
		log.Fatal(err)
	}

	ornn := &ornn.ORNN{}
	ornn.Init(db, config)
	/*
		err = ornn.ConfigLoad("./output/gen.json")
		if err != nil {
			log.Fatal(err)
		}
	*/
	config.InitSchema(db)
	err = config.VendorBySchema()
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
