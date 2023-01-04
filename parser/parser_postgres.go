package parser

import (
	"strings"

	"github.com/auxten/postgresql-parser/pkg/sql/parser"
	"github.com/auxten/postgresql-parser/pkg/walk"
	"github.com/gokch/ornn/config"
)

// TODO
type ParserPostgres struct {
	sch *config.Schema
}

func (p *ParserPostgres) Init(sch *config.Schema) {
	p.sch = sch
}

func (p *ParserPostgres) Parse(sql string) (*ParseQuery, error) {
	w := &walk.AstWalker{
		// Fn: p.walker,
	}

	stmts, err := parser.Parse(sql)
	if err != nil {
		return nil, err
	}

	_, err = w.Walk(stmts, nil)
	if err != nil {
		return nil, err
	}

	parseQuery := &ParseQuery{}
	return parseQuery, nil
}

func (p *ParserPostgres) ConvType(dbType string) (genType string) {
	parseType := ParseType(dbType)

	if strings.HasPrefix(parseType.Type, "SETOF ") {
		genType = p.ConvType(parseType.Type[len("SETOF "):])
		return "[]" + genType
	}
	typNullable := parseType.Nullable && !parseType.IsArray
	typ := parseType.Type
	switch {
	case typ == `"char"`:
		typ = "char"
	case strings.HasPrefix(typ, "information_schema."):
		switch strings.TrimPrefix(typ, "information_schema.") {
		case "cardinal_number":
			typ = "integer"
		case "character_data", "sql_identifier", "yes_or_no":
			typ = "character varying"
		case "time_stamp":
			typ = "timestamp with time zone"
		}
	}
	switch typ {
	case "boolean":
		genType = "bool"
		if typNullable {
			genType = "sql.NullBool"
		}
	case "bpchar", "character varying", "character", "inet", "money", "text", "name":
		genType = "string"
		if typNullable {
			genType = "sql.NullString"
		}
	case "smallint":
		genType = "int16"
		if typNullable {
			genType = "sql.NullInt64"
		} else if parseType.Unsigned {
			genType = "uint16"
		}
	case "integer":
		genType = "int32"
		if typNullable {
			genType = "sql.NullInt64"
		} else if parseType.Unsigned {
			genType = "uint32"
		}
	case "bigint":
		genType = "int64"
		if typNullable {
			genType = "sql.NullInt64"
		} else if parseType.Unsigned {
			genType = "uint64"
		}
	case "real":
		genType = "float32"
		if typNullable {
			genType = "sql.NullFloat64"
		}
	case "double precision", "numeric":
		genType = "float64"
		if typNullable {
			genType = "sql.NullFloat64"
		}
	case "date", "timestamp with time zone", "time with time zone", "time without time zone", "timestamp without time zone":
		genType = "time.Time"
		if typNullable {
			genType = "sql.NullTime"
		}
	case "bit":
		genType = "uint8"
		if typNullable {
			genType = "*uint8"
		}
	case "any", "bit varying", "bytea", "interval", "json", "jsonb", "xml":
		// TODO: write custom type for interval marshaling
		// TODO: marshalling for json types
		genType = "[]byte"
	case "hstore":
		genType = "hstore.Hstore"
	case "uuid":
		genType = "uuidbType.UUID"
		if typNullable {
			genType = "uuidbType.NullUUID"
		}
	default:
		genType = "interface{}"
	}
	return genType
}
