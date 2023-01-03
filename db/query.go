package db

type ParseQuery struct {
	QueryType string
	ClassName string
	QueryName string
	Query     string

	Tpl map[string]string // name:dbtype
	Arg map[string]string // name:dbtype
	Ret map[string]string // name:dbtype

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

type QueryType int8

const (
	QueryTypeSelect QueryType = iota + 1
	QueryTypeInsert
	QueryTypeUpdate
	QueryTypeDelete
)
