package parser

type QueryType int8

const (
	QueryTypeSelect QueryType = iota + 1
	QueryTypeInsert
	QueryTypeUpdate
	QueryTypeDelete
)

type ParseQuery struct {
	QueryType QueryType
	GroupName string
	QueryName string
	Query     string

	Tpl map[string]string // name:type
	Arg map[string]string // name:type
	Ret map[string]string // name:type

	// options
	SelectSingle     bool
	InsertMulti      bool
	UpdateNullIgnore bool
}

func (t *ParseQuery) Init() {
	t.Tpl = make(map[string]string)
	t.Arg = make(map[string]string)
	t.Ret = make(map[string]string)
}
