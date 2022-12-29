package codegen

import (
	"fmt"
	"testing"
)

func TestCodeGen(_t *testing.T) {
	input := struct {
		PackageName string
		StructName  string
		Fields      []string
	}{
		PackageName: "pkg",
		StructName:  "MyStruct",
		Fields: []string{
			"Field1",
			"Field2",
		},
	}

	w := &Writer{}
	w.Init()

	w.W("package %s\n\n", input.PackageName)
	w.W("func (o %[1]s) ShallowCopy() %[1]s {\n", input.StructName)

	{
		w.IndentIn()
		w.W("return %s {\n", input.StructName)
		{
			w.IndentIn()
			for _, field := range input.Fields {
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
