package db_mysql

import (
	"strings"
)

func ConvType(dbType string) (genType string) {
	var unsigned bool

	opts := strings.Split(string(dbType), " ")
	for i, opt := range opts {
		opt = strings.ToLower(opt)
		if opt == "unsigned" {
			unsigned = true

			if i < len(opts)-1 {
				opts = append(opts[:i], opts[i+1:]...)
			} else {
				opts = opts[:i]
			}
		}
	}

	if len(opts) == 0 {
		return ""
	}

	fieldTypeWithLen := opts[0]
	pos := strings.Index(fieldTypeWithLen, "(")
	if pos != -1 {
		fieldTypeWithLen = fieldTypeWithLen[0:pos]
	}

	return convType(fieldTypeWithLen, unsigned)
}

func convType(dbType string, unsigned bool) string {
	switch strings.ToLower(dbType) {
	case "char", "varchar", "tinytext", "text", "mediumtext", "longtext", "json":
		return "string"
	case "binary", "varbinary", "tinyblob", "blob", "mediumblob", "longblob":
		return "[]byte"
	case "tinyint":
		if unsigned == true {
			return "uint8"
		}
		return "int8"
	case "smallint":
		if unsigned == true {
			return "uint16"
		}
		return "int16"
	case "int":
		if unsigned == true {
			return "uint32"
		}
		return "int32"
	case "bigint":
		if unsigned == true {
			return "uint64"
		}
		return "int64"
	case "float":
		return "float32"
	case "double", "real":
		return "float64"
	default:
		return "interface{}"
	}
}
