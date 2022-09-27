package bp

import (
	"fmt"
	"testing"
)

type inputData struct {
	PackageName string
	StructName  string
	Fields      []string
}

func Test__generate__go_source(_t *testing.T) {
	data := inputData{
		PackageName: "pkg",
		StructName:  "MyStruct",
		Fields: []string{
			"Field1",
			"Field2",
		},
	}

	w := &Gen_writer{}
	w.Init()

	w.W("package %s\n\n", data.PackageName)
	w.W("func (o %[1]s) ShallowCopy() %[1]s {\n", data.StructName)

	{
		w.Indent__in()
		w.W("return %s {\n", data.StructName)
		{
			w.Indent__in()
			for _, field := range data.Fields {
				w.W("%[1]s: o.%[1]s,\n", field)
			}
			w.Indent__out()
		}

		w.W("}\n")
		w.Indent__out()
	}

	w.W("}\n")
	w.Indent__out()

	fmt.Println(w.String())
}
