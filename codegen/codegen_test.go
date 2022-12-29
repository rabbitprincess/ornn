package codegen

import (
	"fmt"
	"testing"
)

type inputData struct {
	PackageName string
	StructName  string
	Fields      []string
}

func TestCodeGen(_t *testing.T) {
	data := inputData{
		PackageName: "pkg",
		StructName:  "MyStruct",
		Fields: []string{
			"Field1",
			"Field2",
		},
	}

	w := &Writer{}
	w.Init()

	w.W("package %s\n\n", data.PackageName)
	w.W("func (o %[1]s) ShallowCopy() %[1]s {\n", data.StructName)

	{
		w.IndentIn()
		w.W("return %s {\n", data.StructName)
		{
			w.IndentIn()
			for _, field := range data.Fields {
				w.W("%[1]s: o.%[1]s,\n", field)
			}
			w.IndentOut()
		}

		w.W("}\n")
		w.IndentOut()
	}

	w.W("}\n")
	w.IndentOut()

	fmt.Println(w.String())
}
