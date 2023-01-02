package db_sqlite

func ConvType(dbType string) (goType string) {
	return convType(dbType, false)
}

func convType(dbType string, unsigned bool) (goType string) {
	switch dbType {
	case "bool", "boolean":
		goType = "bool"
	case "int", "integer", "tinyint", "smallint", "mediumint":
		goType = "int32"
		if unsigned {
			goType = "uint32"
		}
	case "bigint":
		goType = "int64"
		if unsigned {
			goType = "uint64"
		}
	case "numeric", "real", "double", "float", "decimal":
		goType = "float64"
	case "blob":
		goType = "[]byte"
	case "timestamp", "datetime", "date", "timestamp with timezone", "time with timezone", "time without timezone", "timestamp without timezone":
		goType = "Time"
	case "varchar", "character", "varying character", "nchar", "native character", "nvarchar", "text", "clob", "time":
		goType = "string"
	default:
		goType = "interface{]"
	}
	return goType
}
