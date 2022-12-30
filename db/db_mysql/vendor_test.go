package db_mysql

import (
	"fmt"
	"testing"
)

func TestConn(t *testing.T) {
	db, err := New("127.0.0.1", "3306", "root", "951753ck", "aergo_indexer")
	if err != nil {
		t.Fatal(err)
	}
	job, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	_ = job
}

func TestVendor(t *testing.T) {
	db, err := New("127.0.0.1", "3306", "root", "951753ck", "aergo_indexer")
	if err != nil {
		t.Fatal(err)
	}
	vendor := NewVendor(db)
	tables, err := vendor.allTable()
	if err != nil {
		t.Fatal(err)
	}
	for _, tbl := range tables {
		fmt.Println(tbl)
	}
}
