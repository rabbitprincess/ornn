package atlas

import (
	"context"
	"fmt"
	"log"
	"testing"

	"ariga.io/atlas/schemahcl"
	"ariga.io/atlas/sql/mysql"
	"github.com/gokch/ornn/db/db_mysql"
)

func TestAtlas(t *testing.T) {
	db, err := db_mysql.New("127.0.0.1", "3306", "root", "1234", "test")
	if err != nil {
		log.Fatal(err)
	}

	atlas, err := mysql.Open(db.Db)
	if err != nil {
		t.Fatal(err)
	}
	sch, err := atlas.InspectSchema(context.Background(), "", nil)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(sch.Name)

	bt, err := schemahcl.Marshal(sch)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(bt))

}
