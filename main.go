package main

import (
	"log"

	"github.com/gokch/ornn/atlas"
	"github.com/gokch/ornn/config"
	"github.com/gokch/ornn/db/db_mysql"
	"github.com/gokch/ornn/ornn"
)

func main() {
	// connect db ( mysql only... )
	db, err := db_mysql.New("127.0.0.1", "3306", "root", "951753ck", "test")
	if err != nil {
		log.Fatal(err)
	}

	// set conf
	conf := &config.Config{}
	{
		conf.Global.InitDefault()
		if err != nil {
			log.Fatal(err)
		}
		err = conf.Load("./output/gen.json")
		if err != nil {
			log.Fatal(err)
		}
		at := &atlas.Atlas{}
		at.Init(atlas.DbTypeMySQL, db)
		schema, err := at.InspectSchema()
		if err != nil {
			log.Fatal(err)
		}
		err = conf.InitSchema(schema)
		if err != nil {
			log.Fatal(err)
		}
		err = conf.Save("./output/gen.json")
		if err != nil {
			log.Fatal(err)
		}
	}

	// init ornn
	ornn := &ornn.ORNN{}
	ornn.Init(db, conf)

	// code generate
	err = ornn.GenCode("./output/gen.go")
	if err != nil {
		log.Fatal(err)
	}
}
