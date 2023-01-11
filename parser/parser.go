package parser

type Parser interface {
	Parse(sql string) (*ParsedQuery, error)
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

	Tpl []*ParsedQueryField // TODO
	Arg []*ParsedQueryField
	Ret []*ParsedQueryField

	// options
	SelectSingle bool
	InsertMulti  bool
}

func (t *ParsedQuery) Init(query string) {
	t.Query = query
	t.Tpl = make([]*ParsedQueryField, 0, 10)
	t.Arg = make([]*ParsedQueryField, 0, 10)
	t.Ret = make([]*ParsedQueryField, 0, 10)
}

func NewField(name, goType string) *ParsedQueryField {
	return &ParsedQueryField{
		Name:   name,
		GoType: goType,
	}
}

type ParsedQueryField struct {
	Name   string
	GoType string
}
