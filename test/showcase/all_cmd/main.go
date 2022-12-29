package main

import (
	"log"
	"module/solution/bp"
)

//go:generate go run main.go

func main() {
	err := bp.Generate__mysql__golang(
		"127.0.0.1",
		"4001",
		"root",
		"1234",
		"dev_bp_sample",
		"",
		"bp.json",
		"./bp_result/bp_result.go",
		"bp_db",
		"C_DB")
	if err != nil {
		log.Fatal(err)
	}
}
