package parser_mysql

import (
	"regexp"

	"ariga.io/atlas/sql/schema"
	"github.com/gokch/ornn/parser"
)

func (p *Parser) ConvType(colType *schema.ColumnType) (genType string) {
	parseType := parser.ParseType(colType.Raw)
	parseType.Nullable = colType.Null
	switch parseType.Type {
	case "bit":
		switch {
		case parseType.Prec == 1 && !parseType.Nullable:
			genType = "bool"
		case parseType.Prec == 1 && parseType.Nullable:
			genType = "sql.NullBool"
		case parseType.Prec <= 8 && !parseType.Nullable:
			genType = "uint8"
		case parseType.Prec <= 16 && !parseType.Nullable:
			genType = "uint16"
		case parseType.Prec <= 32 && !parseType.Nullable:
			genType = "uint32"
		case parseType.Nullable:
			genType = "sql.NullInt64"
		default:
			genType = "uint64"
		}
	case "bool", "boolean":
		genType = "bool"
		if parseType.Nullable {
			genType = "sql.NullBool"
		}
	case "char", "varchar", "tinytext", "text", "mediumtext", "longtext":
		genType = "string"
		if parseType.Nullable {
			genType = "sql.NullString"
		}
	case "tinyint":
		switch {
		case parseType.Prec == 1 && !parseType.Nullable: // force tinyint(1) as bool
			genType = "bool"
		case parseType.Prec == 1 && parseType.Nullable:
			genType = "sql.NullBool"
		case parseType.Nullable:
			genType = "sql.NullInt64"
		default:
			genType = "int8"
			if parseType.Unsigned {
				genType = "uint8"
			}
		}
	case "smallint", "year":
		genType = "int16"
		if parseType.Nullable {
			genType = "sql.NullInt64"
		} else if parseType.Unsigned {
			genType = "uint16"
		}
	case "mediumint", "int", "integer":
		genType = "int32"
		if parseType.Nullable {
			genType = "sql.NullInt64"
		} else if parseType.Unsigned {
			genType = "uint32"
		}
	case "bigint":
		genType = "int64"
		if parseType.Nullable {
			genType = "sql.NullInt64"
		} else if parseType.Unsigned {
			genType = "uint64"
		}
	case "float":
		genType = "float32"
		if parseType.Nullable {
			genType = "sql.NullFloat64"
		}
	case "decimal", "double":
		genType = "float64"
		if parseType.Nullable {
			genType = "sql.NullFloat64"
		}
	case "binary", "blob", "longblob", "mediumblob", "tinyblob", "varbinary":
		genType = "[]byte"
	case "json":
		genType = "json.RawMessage"
	case "timestamp", "datetime", "date":
		genType = "time.Time"
		if parseType.Nullable {
			genType = "sql.NullTime"
		}
	case "time":
		genType = "string"
		if parseType.Nullable {
			genType = "sql.NullString"
		}
	default:
		genType = "interface{}"
	}
	if regexp.MustCompile(`(?i)^set\([^)]*\)$`).MatchString(parseType.Type) {
		genType = "[]byte"
	}

	return genType
}
