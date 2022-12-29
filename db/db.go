package db

import (
	"context"
	"database/sql"
)

type DB struct {
	Job

	DriverName string
	Dsn        string
	DbName     string

	Db *sql.DB
}

func (t *DB) Connect(driverName, dsn, dbName string) (err error) {
	t.DriverName = driverName
	t.Dsn = dsn
	t.DbName = dbName

	t.Db, err = sql.Open(t.DriverName, t.Dsn)
	if err != nil {
		return err
	}

	err = t.Db.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (t *DB) SetOpenConns(_openConns, _idleConns int) {
	if _openConns > 0 {
		t.Db.SetMaxOpenConns(_openConns)
	}
	if _idleConns > 0 {
		t.Db.SetMaxIdleConns(_idleConns)
	}
}

func (t *DB) TxBegin(_isoLevel sql.IsolationLevel, _readonly bool) (tx *sql.Tx, err error) {
	pt_opt := &sql.TxOptions{
		Isolation: _isoLevel,
		ReadOnly:  _readonly,
	}

	return t.Db.BeginTx(context.Background(), pt_opt)
}

func (t *DB) TxBegin__callback(_isoLevel sql.IsolationLevel, _readonly bool, _fnCallback func(*sql.Tx) error) (err error) {
	pt_opt := &sql.TxOptions{
		Isolation: _isoLevel,
		ReadOnly:  _readonly,
	}

	tx, err := t.Db.BeginTx(context.Background(), pt_opt)
	if err != nil {
		return err
	}

	err = _fnCallback(tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
