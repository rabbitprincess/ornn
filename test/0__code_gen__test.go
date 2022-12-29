package go_orm_gen_test

import (
	"module/solution/bp"
	"testing"
)

//------------------------------------------------------------------------------//
// generate go code

func Test__Generate__mysql__golang(_t *testing.T) {
	var err error
	err = bp.Generate__mysql__golang(
		"127.0.0.1",
		"4001",
		"root",
		"1234",
		"dev_bp_test",
		"",
		"bp.json",
		"./test_output/bp_result.go",
		"bp_db",
		"C_DB")
	if err != nil {
		_t.Fatal(err)
	}
}
