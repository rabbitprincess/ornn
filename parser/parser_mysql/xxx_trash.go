package parser_mysql

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

type Delete struct {
	TableAs []*TableAs
}

func (t *Delete) parse(psr *sqlparser.Delete) error {
	for _, tableExprs := range psr.TableExprs {
		tableAs := &TableAs{}

		switch data := tableExprs.(type) {
		case *sqlparser.AliasedTableExpr: // 단순 테이블
			tableAs.Table = data.Expr.(sqlparser.TableName).Name.String()
			tableAs.As = data.As.String()
		case *sqlparser.ParenTableExpr:
			// 임시 - 작업필요
			// -> 반드시 sub query 를 재귀호출로 해체하여 제일 외부 에 있는 () 에 대해 서만 table list 에 남긴다. = *  타입 지정 문제
			log.Fatal("need more programming")
		case *sqlparser.JoinTableExpr:
			// 임시 - 작업필요
			// -> 반드시 sub query 를 재귀호출로 해체하여 제일 외부 에 있는 () 에 대해 서만 table list 에 남긴다. = *  타입 지정 문제
			log.Fatal("need more programming")
		}
		t.AddTable(tableAs)
	}


		// where
		// 재귀를 통해 모든 where 문을 배열로 전처리
		// select 일 경우 select 문 재귀 처리
		// in 또는 not in 이면 multi arg 처리
		switch data := sqlparser.Where.Expr.(type) {
		case *AndExpr:
			{
				// 임시 - 작업 예정 - multi where 일 경우 재귀를 통해 타입을 하나씩 얻어낸다
			}
		case *ComparisonExpr:
			{
				if data.Operator == InStr || data.Operator == NotInStr {

				}
			}
		}


	return nil
}

func (t *Delete) AddTable(as *TableAs) {
	if t.TableAs == nil {
		t.TableAs = make([]*TableAs, 0, 10)
	}
	t.TableAs = append(t.TableAs, as)
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


type Insert struct {
	TableName string
	Fields    []*Field
}

func (t *Insert) parse(psr *sqlparser.Insert) error {
	// table name
	t.TableName = psr.Table.Name.String()

	switch row := psr.Rows.(type) {
	case sqlparser.Values:
		if len(row) != 1 {
			// bp config 에서 multi insert 쿼리 직접 입력 금지
			// multi insert 는 bp config 의 mutli insert option 을 true 로 설정하고, 단일 쿼리를 그대로 사용
			log.Fatal("parser error - multi insert query is forbidden. use multi insert option")
		}

		lenField := len(psr.Columns)
		lenVal := len(row[0])
		if lenField != lenVal && lenField != 0 { // 0 일때는 필드명이 모두 비어있는 상태
			return fmt.Errorf("field name length is not same as value length")
		}

		for i, val := range row[0] {
			field := &Field{}
			if psr.Columns != nil {
				field.FieldName = psr.Columns[i].String()
			}

			// value 값이 NULL 이면 처리하지 않음
			if _, ok := val.(*sqlparser.NullVal); ok == true {
				field.Val = []byte("")
				t.addField(field)
				continue
			}

			if sqlVal, ok := val.(*sqlparser.SQLVal); ok == true {
				field.Val = sqlVal.Val
				t.addField(field)
				continue
			}

			return fmt.Errorf("parser error - unexpacted type of field value")
		}
	case *sqlparser.Select:
		// 임시 - 작업필요
		// -> insert 에 입력값으로 select query 를 사용하는 경우로
		// -> 추후 작업할때 where 절에 ? 를 처리해주면 됨
		log.Fatal("need more programming")
	default:
		log.Fatal("need more programming")
	}

	return nil
}

func (t *Insert) addField(field *Field) {
	if t.Fields == nil {
		t.Fields = make([]*Field, 0, 10)
	}

	t.Fields = append(t.Fields, field)
}


type Select struct {
	FieldAs []*FieldAs
	TableAs []*TableAs
	Offset  *int64
	Limit   *int64
}

func (t *Select) parse(psr *sqlparser.Select) error {
	// from
	for _, from := range psr.From {
		tableAs := &TableAs{}

		switch data := from.(type) {
		case *sqlparser.AliasedTableExpr: // 단순 테이블
			tableAs.Table = data.Expr.(sqlparser.TableName).Name.String()
			tableAs.As = data.As.String()
		case *sqlparser.ParenTableExpr:
			// 임시 - 작업필요
			// -> 반드시 sub query 를 재귀호출로 해체하여 제일 외부 에 있는 () 에 대해 서만 table list 에 남긴다. = *  타입 지정 문제
			log.Fatal("need more programming")
		case *sqlparser.JoinTableExpr:
			// 임시 - 작업필요
			// -> 반드시 sub query 를 재귀호출로 해체하여 제일 외부 에 있는 () 에 대해 서만 table list 에 남긴다. = *  타입 지정 문제
			log.Fatal("need more programming")
		}

		t.addTable(tableAs)
	}

	// select
	for _, selectExpr := range psr.SelectExprs {
		fieldAs := &FieldAs{}

		switch data := selectExpr.(type) {
		case *sqlparser.StarExpr:
			fieldAs.Field = "*"
			fieldAs.Table = data.TableName.Name.String()
			if fieldAs.Table == "" {
				// 임시 - 작업중
			}
		case *sqlparser.AliasedExpr:
			fieldAs.Field = data.Expr.(*sqlparser.ColName).Name.String()
			fieldAs.As = data.As.String()
			fieldAs.Table = data.Expr.(*sqlparser.ColName).Qualifier.Name.String()

		case sqlparser.Nextval:
			log.Fatal("need more programming")
		}

		t.addField(fieldAs)
	}


		// where
		// 재귀를 통해 모든 where 를 배열로 얻기
		// where in or not in 이면 multi arg 처리
		switch data := psr.Where.Expr.(type) {
		case *AndExpr:
			{
				// 임시 - 작업 예정 - multi where 일 경우 재귀를 통해 타입을 하나씩 얻어낸다
			}
		case *ComparisonExpr:
			{
				if data.Operator == InStr || data.Operator == NotInStr {

				}
			}
		}


	// limit, offset
	if psr.Limit != nil {
		if offset, ok := psr.Limit.Offset.(*sqlparser.SQLVal); ok == true {
			if offset.Type == sqlparser.IntVal {
				n, err := strconv.ParseInt(string(offset.Val), 10, 64)
				if err == nil {
					t.Offset = &n
				}
			}
		}
		if limit, ok := psr.Limit.Rowcount.(*sqlparser.SQLVal); ok == true {
			if limit.Type == sqlparser.IntVal {
				n, err := strconv.ParseInt(string(limit.Val), 10, 64)
				if err == nil {
					t.Limit = &n
				}
			}
		}
	}

	return nil
}

func (t *Select) GetTableNames() []string {
	tables := make([]string, len(t.TableAs))

	for i, table := range t.TableAs {
		tables[i] = table.Table
	}
	return tables
}

func (t *Select) getTableName(tableNameSql string) (tableName string, err error) {
	// 지정된 table name 이 없는경우 해석하지 않는다.
	if tableNameSql == "" {
		return "", nil
	}

	// 지정된 table name 이 있는 경우 as 인지 schema 인지 구분하여 리턴
	for _, table := range t.TableAs {
		tableNameSchema := table.get()
		if tableNameSchema == tableNameSql {
			return table.Table, nil // schema table 리턴
		}
	}

	// 테이블 이름이 매칭이 안되는 케이스 - 입력값(json) 오류
	return "", fmt.Errorf("not exist table name in select expr - table name : %s", tableNameSql)
}

func (t *Select) addField(as *FieldAs) {
	if t.FieldAs == nil {
		t.FieldAs = make([]*FieldAs, 0, 10)
	}
	t.FieldAs = append(t.FieldAs, as)
}

func (t *Select) addTable(as *TableAs) {
	if t.TableAs == nil {
		t.TableAs = make([]*TableAs, 0, 10)
	}
	t.TableAs = append(t.TableAs, as)
}


type Update struct {
	TableAs []*TableAs
	Field   []*Field
}

func (t *Update) parse(psr *sqlparser.Update) error {
	for _, tableExpr := range psr.TableExprs {
		tableAs := &TableAs{}

		switch data := tableExpr.(type) {
		case *sqlparser.AliasedTableExpr: // 단순 테이블
			tableAs.Table = data.Expr.(sqlparser.TableName).Name.String()
			tableAs.As = data.As.String()
		case *sqlparser.ParenTableExpr:
			// 임시 - 작업필요
			// -> 반드시 sub query 를 재귀호출로 해체하여 제일 외부 에 있는 () 에 대해 서만 table list 에 남긴다. = *  타입 지정 문제
			log.Fatal("need more programming")
		case *sqlparser.JoinTableExpr:
			// 임시 - 작업필요
			// -> 반드시 sub query 를 재귀호출로 해체하여 제일 외부 에 있는 () 에 대해 서만 table list 에 남긴다. = *  타입 지정 문제
			log.Fatal("need more programming")
		}
		t.AddTable(tableAs)
	}

	// field value
	for _, pt_expr := range psr.Exprs {

		pt_field := &Field{}
		pt_field.TableName = pt_expr.Name.Qualifier.Name.String()
		pt_field.FieldName = pt_expr.Name.Name.String()
		switch data := pt_expr.Expr.(type) {
		case *sqlparser.SQLVal:
			{
				pt_field.Val = data.Val
			}
		case *sqlparser.BinaryExpr:
			{
				data_left, is_ok_left := data.Left.(*sqlparser.SQLVal)
				data_right, is_ok_right := data.Right.(*sqlparser.SQLVal)
				if is_ok_left == true && is_ok_right == true {
					// 양쪽 다 val 일 경우 - 에러 ( set u8_num = %u8_num% + %u8_num% )
					log.Fatal("need more programming")
				} else if is_ok_left == true {
					// 왼쪽이 val 일 경우 ( set u8_num = %u8_num% + u8_num )
					pt_field.Val = data_left.Val
				} else if is_ok_right == true {
					// 오른쪽이 val 일 경우 ( set u8_num = u8_num + %u8_num% )
					pt_field.Val = data_right.Val
				}
			}
		case *sqlparser.FuncExpr:
			{
				s_func_name := strings.ToLower(data.Name.String())
				switch s_func_name {
				case "ifnull": // 임시 - 개선 필요 - 현재 ifnull 만 적용
					{
						// update 시 nil 값을 입력하면 업데이트 하지 않도록 하기 위함
						// sql - "seq = ifnull(%seq%, seq)" 식으로 작성
						pt_alised_expr, is_ok := data.Exprs[0].(*sqlparser.AliasedExpr)
						if is_ok == false {
							log.Fatal("need more programming")
						}
						pt_val, is_ok := pt_alised_expr.Expr.(*sqlparser.SQLVal)
						pt_field.Val = pt_val.Val
						// is_pointer 플래그 추가
						// nil 을 입력받을 수 있어야 함으로 포인터로 세팅
					}
				default:
					{
						log.Fatal("need more programming")
					}
				}
			}
		default:
			log.Fatal("need more programming")
		}

		t.AddField(pt_field)
	}


		// where
		// 재귀를 통해 모든 where 문을 배열로 전처리
		// select 일 경우 select 문 재귀 처리
		// in 또는 not in 이면 multi arg 처리
		switch data := _pt_parser.Where.Expr.(type) {
		case *AndExpr:
			{
				// 임시 - 작업 예정 - multi where 일 경우 재귀를 통해 타입을 하나씩 얻어낸다
			}
		case *ComparisonExpr:
			{
				if data.Operator == InStr || data.Operator == NotInStr {

				}
			}
		}


	return nil
}

func (t *Update) GetTableNames() []string {
	tableName := make([]string, len(t.TableAs))

	for i, pt_table := range t.TableAs {
		tableName[i] = pt_table.Table
	}
	return tableName
}

func (t *Update) AddField(field *Field) {
	if t.Field == nil {
		t.Field = make([]*Field, 0, 10)
	}

	t.Field = append(t.Field, field)
}

func (t *Update) AddTable(tableAs *TableAs) {
	if t.TableAs == nil {
		t.TableAs = make([]*TableAs, 0, 10)
	}
	t.TableAs = append(t.TableAs, tableAs)
}

*/
