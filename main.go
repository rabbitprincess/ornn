package main

import (
	"log"

	"github.com/gokch/ornn/atlas"
	"github.com/gokch/ornn/config"
	"github.com/gokch/ornn/db/db_postgres"
	"github.com/gokch/ornn/ornn"
)

func main() {
	// connect db
	dbType := atlas.DbTypePostgre
	db, err := db_postgres.New("127.0.0.1", "5432", "postgres", "", "postgres")
	if err != nil {
		log.Fatal(err)
	}

	// dbType := atlas.DbTypeMySQL
	// db, err := db_mysql.New("127.0.0.1", "3306", "root", "1234", "mysql")
	// if err != nil {
	// log.Fatal(err)
	// }

	// set conf
	conf := &config.Config{}
	{
		// load
		err = conf.Load("./output/gen.json")
		if err != nil {
			log.Fatal(err)
		}
		// init
		conf.Global.InitDefault()
		if err != nil {
			log.Fatal(err)
		}
		schema, err := atlas.InspectSchema(dbType, db)
		if err != nil {
			log.Fatal(err)
		}
		err = conf.InitSchema(dbType, schema)
		if err != nil {
			log.Fatal(err)
		}
		// save
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
