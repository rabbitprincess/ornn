package parser

import (
	"fmt"
	"regexp"

	"github.com/gokch/ornn/config"
	tiparser "github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	_ "github.com/pingcap/tidb/parser/test_driver"
)

// TODO
type ParserMysql struct {
	conf *config.Schema
}

func (p *ParserMysql) Parse(sql string) (*ParseQuery, error) {
	ps := tiparser.New()
	stmtNodes, _, err := ps.Parse(sql, "", "")
	if err != nil {
		return nil, err
	}

	parseQuery := &ParseQuery{}
	parseQuery.Init()

	for _, stmtNode := range stmtNodes {
		switch stmt := stmtNode.(type) {
		case *ast.SelectStmt:
			p.parseSelect(stmt, parseQuery)
		case *ast.InsertStmt:

		case *ast.UpdateStmt:

		case *ast.DeleteStmt:

		default:
			return nil, fmt.Errorf("parser error | not support query statement %T", stmt)
		}
	}

	return parseQuery, nil
}

func (p *ParserMysql) parseSelect(stmt *ast.SelectStmt, parseQuery *ParseQuery) {
	// select
	for _, field := range stmt.Fields.Fields {
		var fieldName, fieldType string
		fieldName = field.AsName.O
		fmt.Println(field.AsName.O)
		fmt.Println(field.Expr.GetType().String())
		/*
			if fieldType == "" {
				colType, _ := p.conf.GetFieldType("", fieldName)
				fieldType = p.ConvType(colType)
			}
		*/
		parseQuery.Ret[fieldName] = fieldType
	}

	// where

	// single select 처리
	// 코드 생성 시 단일 구조체 반환 목적
	if stmt.Limit.Count.Text() == "1" {
		parseQuery.SelectSingle = true
	}
	return
}

func (p *ParserMysql) ConvType(dbType string) (genType string) {
	parseType := ParseType(dbType)
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

/*
func Select(db *db.Conn, conf *config.Config, query *config.Query, parseQuery *parser.ParseQuery, sqlSelect *parser.Select) error {
	// 필드 정보를 얻어온다.
	sqlWithoutWhere, _ := sql.Util_SplitByDelimiter(query.Sql, "where")
	sqlAfterArg := sql.Util_ReplaceBetweenDelimiter(sqlWithoutWhere, sql.PrepareStatementDelimeter, sql.PrepareStatementAfter)
	sqlAfterArgClearTpl := sql.Util_ReplaceInDelimiter(sqlAfterArg, sql.TplDelimiter, sql.TplSplit)

	rows, err := db.Job().Query(sqlAfterArgClearTpl)
	if err != nil {
		return err
	}

	cols, err := rows.ColumnTypes()
	if err != nil {
		return err
	}

	for _, col := range cols {
		var fieldName, fieldType string
		fieldName = col.Name()
		fieldType = query.GetCustomType(fieldName)

		// if custom type is not defined, get database type
		if fieldType == "" {
			colType, _ := conf.Schema.GetFieldType("", fieldName)
			fieldType = conf.Schema.ConvType(colType)
		}
		parseQuery.Ret[fieldName] = fieldType
	}

	// single select 처리
	// 코드 생성 시 단일 구조체 반환 목적
	if sqlSelect.Limit != nil && *(sqlSelect.Limit) == 1 {
		parseQuery.SelectSingle = true
	}
	return nil
}

func Insert(conf *config.Config, query *config.Query, parseQuery *parser.ParseQuery, sqlInsert *parser.Insert) error {
	// 필드 정보를 얻어온다.
	schemaTable, exist := query.Schema.Table(sqlInsert.TableName)
	if exist != true {
		return fmt.Errorf("table name is not exist | table name - %s", sqlInsert.TableName)
	}

	// 스키마와 파서의 전체 필드 숫자가 다르면 -> 파서에서 모든 필드 이름이 제공되어야 함 -> 하나라도 없으면 에러
	if len(sqlInsert.Fields) != len(schemaTable.Columns) {
		for _, field := range sqlInsert.Fields {
			if field.FieldName == "" {
				return fmt.Errorf("field name is empty")
			}
		}
	} else {
		// 스키마와 파서의 전체 필드수가 같으면 -> 파서에서 모든 필드 이름이 없어도 가능 -> 스키마에서 추출하여 모든 필드명을 채움
		for i, field := range sqlInsert.Fields {
			field.FieldName = schemaTable.Columns[i].Name
		}
	}

	// 필드 이름을 모두 채운 상태에서 처리 시작
	for _, field := range sqlInsert.Fields {
		// 입력값이 ? (arg) 형식이 아니면 func arg 를 만들 필요가 없음으로 continue
		if sql.Util_IsParserValArg(field.Val) == false {
			continue
		}

		// 입력값이 ? (arg) 일 때만 필드이름 조사 = func arg 의 name 으로 활용
		schemaField, exist := schemaTable.Column(field.FieldName)
		if exist != true {
			return fmt.Errorf("not exist field in schema | field name : %s", field.FieldName)
		}

		parseQuery.Arg[field.FieldName] = conf.Schema.ConvType(schemaField.Type.Raw)
	}

	// multi insert 처리
	parseQuery.InsertMulti = query.InsertMulti

	return nil
}

func Update(conf *config.Config, query *config.Query, parseQuery *parser.ParseQuery, sqlUpdate *parser.Update) error {
	// set
	for _, field := range sqlUpdate.Field {
		// 입력값이 ? (arg) 형식이 아니면 func arg 를 만들 필요가 없음으로 continue
		if sql.Util_IsParserValArg(field.Val) == false {
			continue
		}

		fieldName := field.FieldName
		tableName := field.TableName

		// 정의된 table name 이 없으면 update 대상 테이블 중 매칭되는 테이블을 찾는다
		if tableName == "" {
			tables := sqlUpdate.GetTableNames()
			tablesMatch, err := query.Schema.GetTableFieldMatched(fieldName, tables)
			if err != nil {
				return err
			}

			// parse 에러 처리
			// 두개 이상의 테이블이 매칭됨
			if len(tablesMatch) > 1 {
				var dup string
				for _, table := range tablesMatch {
					dup += fmt.Sprintf("%s, ", table)
				}
				dup = dup[:len(dup)-2]
				return fmt.Errorf("duplicated field name in multiple table | field name - %s | tables name - %s", fieldName, dup)
			}
			// 매칭되는 테이블이 한개도 없음
			if len(tablesMatch) == 0 {
				return fmt.Errorf("no tables match the field | field name - %s", fieldName)
			}

			// 테이블 이름 설정 ( 임시 - 현재는 0번 테이블 )
			tableName = tablesMatch[0]
		}

		// 테이블과 필드 이름을 이용해 필드 타입을 찾아낸다
		var genType string
		{
			schemaTable, exist := query.Schema.Table(tableName)
			if exist != true {
				return fmt.Errorf("not exist table | table name - %s", tableName)
			}
			schemaField, exist := schemaTable.Column(fieldName)
			if exist != true {
				return fmt.Errorf("not exist field | field name - %s", field.FieldName)
			}
			genType = conf.Schema.ConvType(schemaField.Type.Raw)
		}

		parseQuery.Arg[field.FieldName] = genType
	}
	// update 시 null 값 ignore 처리
	parseQuery.UpdateNullIgnore = query.UpdateNullIgnore

	return nil
}

func Delete(conf *config.Config, query *config.Query, parseQuery *parser.ParseQuery, sqlDelete *parser.Delete) error {
	return nil
}
*/
