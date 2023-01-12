package codegen

type CodeGen struct {
	Global
}

func (t *CodeGen) Code() (code string) { // for I_Gen interface
	writer := &Writer{}
	writer.Init()
	t.Global.Code(writer)

	return writer.String()
}
