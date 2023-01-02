package main

import (
	"log"

	"github.com/gokch/ornn/atlas"
	"github.com/gokch/ornn/config"
	"github.com/gokch/ornn/db/db_mysql"
	"github.com/gokch/ornn/ornn"
)

func main() {
	db, err := db_mysql.New("127.0.0.1", "3306", "root", "951753ck", "test")
	if err != nil {
		log.Fatal(err)
	}

	config := &config.Config{}
	config.Global.InitDefault()
	if err != nil {
		log.Fatal(err)
	}

	at := &atlas.Atlas{}
	at.Init(atlas.DbTypeMySQL, db)
	schema, err := at.InspectSchema()
	if err != nil {
		log.Fatal(err)
	}

	ornn := &ornn.ORNN{}
	ornn.Init(db, config)

	err = ornn.ConfigLoad("./output/gen.json")
	if err != nil {
		log.Fatal(err)
	}

	config.InitSchema(schema)
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
