package db_mysql

import (
	"regexp"

	"github.com/gokch/ornn/db"
)

func ConvType(dbType string) (goType string) {
	parseType := db.ParseType(dbType)
	switch parseType.Type {
	case "bit":
		switch {
		case parseType.Prec == 1 && !parseType.Nullable:
			goType = "bool"
		case parseType.Prec == 1 && parseType.Nullable:
			goType = "sql.NullBool"
		case parseType.Prec <= 8 && !parseType.Nullable:
			goType = "uint8"
		case parseType.Prec <= 16 && !parseType.Nullable:
			goType = "uint16"
		case parseType.Prec <= 32 && !parseType.Nullable:
			goType = "uint32"
		case parseType.Nullable:
			goType = "sql.NullInt64"
		default:
			goType = "uint64"
		}
	case "bool", "boolean":
		goType = "bool"
		if parseType.Nullable {
			goType = "sql.NullBool"
		}
	case "char", "varchar", "tinytext", "text", "mediumtext", "longtext":
		goType = "string"
		if parseType.Nullable {
			goType = "sql.NullString"
		}
	case "tinyint":
		switch {
		case parseType.Prec == 1 && !parseType.Nullable: // force tinyint(1) as bool
			goType = "bool"
		case parseType.Prec == 1 && parseType.Nullable:
			goType = "sql.NullBool"
		case parseType.Nullable:
			goType = "sql.NullInt64"
		default:
			goType = "int8"
			if parseType.Unsigned {
				goType = "uint8"
			}
		}
	case "smallint", "year":
		goType = "int16"
		if parseType.Nullable {
			goType = "sql.NullInt64"
		} else if parseType.Unsigned {
			goType = "uint16"
		}
	case "mediumint", "int", "integer":
		goType = "int32"
		if parseType.Nullable {
			goType = "sql.NullInt64"
		} else if parseType.Unsigned {
			goType = "uint32"
		}
	case "bigint":
		goType = "int64"
		if parseType.Nullable {
			goType = "sql.NullInt64"
		} else if parseType.Unsigned {
			goType = "uint64"
		}
	case "float":
		goType = "float32"
		if parseType.Nullable {
			goType = "sql.NullFloat64"
		}
	case "decimal", "double":
		goType = "float64"
		if parseType.Nullable {
			goType = "sql.NullFloat64"
		}
	case "binary", "blob", "longblob", "mediumblob", "tinyblob", "varbinary":
		goType = "[]byte"
	case "json":
		goType = "json.RawMessage"
	case "timestamp", "datetime", "date":
		goType = "time.Time"
		if parseType.Nullable {
			goType = "sql.NullTime"
		}
	case "time":
		goType = "string"
		if parseType.Nullable {
			goType = "sql.NullString"
		}
	default:
		goType = "interface{}"
	}
	if regexp.MustCompile(`(?i)^set\([^)]*\)$`).MatchString(parseType.Type) {
		goType = "[]byte"
	}

	return goType
}
