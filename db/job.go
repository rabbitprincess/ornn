package db

import "database/sql"

type Job struct { // db or tx
	isTx bool

	db *sql.DB
	tx *sql.Tx
}

func (t *Job) Init(isTx bool, db *sql.DB, tx *sql.Tx) {
	t.isTx = isTx
	t.db = db
	t.tx = tx
}

// args 제작 예정
func (t *Job) Exec(query string, args ...interface{}) (res sql.Result, err error) {
	if t.isTx == false {
		res, err = t.db.Exec(query, args...)
	} else {
		res, err = t.tx.Exec(query, args...)
	}
	return res, err
}

func (t *Job) Query(query string, _args ...interface{}) (rows *sql.Rows, err error) {
	if t.isTx == false {
		rows, err = t.db.Query(query, _args...)
	} else {
		rows, err = t.tx.Query(query, _args...)
	}
	return rows, err
}
