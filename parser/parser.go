package parser

type Parser interface {
	Parse(sql string) (*ParsedQuery, error)
	ConvType(dbType string) (genType string)
}

type QueryType int8

const (
	QueryTypeSelect QueryType = iota + 1
	QueryTypeInsert
	QueryTypeUpdate
	QueryTypeDelete
)

type ParsedQuery struct {
	QueryType QueryType
	Query     string

	Tpl map[string]string // name:type
	Arg map[string]string // name:type
	Ret map[string]string // name:type

	// options
	SelectSingle     bool
	InsertMulti      bool
	UpdateNullIgnore bool
}

func (t *ParsedQuery) Init(query string) {
	t.Query = query
	t.Tpl = make(map[string]string)
	t.Arg = make(map[string]string)
	t.Ret = make(map[string]string)
}
