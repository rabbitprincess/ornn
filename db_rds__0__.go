package go_orm_gen

import (
	"log"
	"module/db"

	"github.com/blastrain/vitess-sqlparser/sqlparser"
)

// RDS 종류
type TD_N1_db_rds_type int8

const (
	TD_N1_db_rds_type__unknown TD_N1_db_rds_type = iota + 0
	TD_N1_db_rds_type__mysql
)

// 기본 파일명
const (
	DEF_S_default_filepath__json   = "bp.json"
	DEF_S_default_filepath__go__db = "bp_db.go"
)

type T_DB_RDS struct {
	pc_db             *db.C_DB_conn
	td_n1_db_rds_type TD_N1_db_rds_type
	i_dbms            I_DB_RDS__vendor
}

//------------------------------------------------------------------------------------------------//
// Schema

func (t *T_DB_RDS) Init(_pc_db *db.C_DB_conn, _td_n1_db_rds_type TD_N1_db_rds_type) {
	t.pc_db = _pc_db

	t.td_n1_db_rds_type = _td_n1_db_rds_type

	// set dbms
	switch _td_n1_db_rds_type {
	case TD_N1_db_rds_type__mysql:
		t.i_dbms = &T_DB_RDS__vendor__mysql{}
	default:
		log.Fatal("DBMS Not support yet. | type iota - %d", _td_n1_db_rds_type)
	}
}

func (t *T_DB_RDS) Conv_field_type__to_bp(_s_field_type__db string) string {
	return t.i_dbms.conv_field_type__to_bp(_s_field_type__db)
}

func (t *T_DB_RDS) get_schema() (*T_BP__config__db__schema, error) {
	// sql(create table) 추출
	arrs_sql__create_table, err := t.i_dbms.get_sql__create_table(t.pc_db)
	if err != nil {
		return nil, err
	}

	pt_json__db__schema := &T_BP__config__db__schema{}
	pt_json__db__schema.init()

	// 뽑아낸 쿼리를 이용해서 table 을 제작
	for _, s_sql__create_table := range arrs_sql__create_table {
		pt_table, err := t.get_schema____create_table(s_sql__create_table)
		if err != nil {
			return nil, err
		}
		pt_json__db__schema.add_table(pt_table)
	}
	return pt_json__db__schema, nil
}

func (t *T_DB_RDS) get_schema____create_table(_s_sql__create_table string) (*T_BP__config__db__schema__table, error) {
	stmt, err := sqlparser.Parse(_s_sql__create_table)
	if err != nil {
		panic(err)
	}
	pt_parser := stmt.(*sqlparser.CreateTable)
	s_table_name := pt_parser.NewName.Name.String()

	pt_table := &T_BP__config__db__schema__table{}
	pt_table.init(s_table_name)

	for _, pt_index__parser := range pt_parser.Constraints {
		pt_index := &T_BP__config__db__schema__index{}

		arrs_keys := make([]string, 0, len(pt_index__parser.Keys))
		for _, pt_key := range pt_index__parser.Keys {
			arrs_keys = append(arrs_keys, pt_key.String())
		}
		pt_index.set(
			pt_index__parser.Name,
			pt_index__parser.Type.String(),
			arrs_keys,
		)
		pt_table.add_index(pt_index)
	}

	for _, pt_field__parser := range pt_parser.Columns {
		pt_field := &T_BP__config__db__schema__field{}
		pt_field.set(
			pt_field__parser.Name,
			pt_field__parser.Type,
			t.i_dbms.conv_field_type__to_bp(pt_field__parser.Type),
		)

		pt_table.add_field(pt_field)
	}
	return pt_table, nil
}
