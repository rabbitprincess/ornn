package atlas

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"ariga.io/atlas/sql/mysql"
	"ariga.io/atlas/sql/schema"
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

	bt, err := mysql.MarshalHCL(sch)

	schemaNew := &schema.Schema{}

	err = os.WriteFile("./gen.hcl", bt, 0700)
	if err != nil {
		t.Fatal(err)
	}

	err = mysql.EvalHCLBytes(bt, schemaNew, nil)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(schemaNew.Name)
	fmt.Printf("%T type", schemaNew)
	for _, table := range schemaNew.Tables {
		fmt.Println(table.Name)
	}

}
