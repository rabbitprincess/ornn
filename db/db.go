package db

import (
	"context"
	"database/sql"
)

type Conn struct {
	DriverName string
	Dsn        string
	DbName     string

	Db *sql.DB
}

func (t *Conn) Connect(driverName, dsn, dbName string) (err error) {
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

func (t *Conn) SetOpenConns(openConns, idleConns int) {
	if openConns > 0 {
		t.Db.SetMaxOpenConns(openConns)
	}
	if idleConns > 0 {
		t.Db.SetMaxIdleConns(idleConns)
	}
}

func (t *Conn) Begin() (*Job, error) {
	tx, err := t.Db.Begin()
	if err != nil {
		return nil, err
	}
	job := &Job{}
	job.Init(false, t.Db, tx)
	return job, nil
}

func (t *Conn) TxBegin(isoLevel sql.IsolationLevel, readonly bool) (tx *sql.Tx, err error) {
	pt_opt := &sql.TxOptions{
		Isolation: isoLevel,
		ReadOnly:  readonly,
	}

	return t.Db.BeginTx(context.Background(), pt_opt)
}

func (t *Conn) TxBeginCb(isoLevel sql.IsolationLevel, readonly bool, fnCallback func(*sql.Tx) error) (err error) {
	pt_opt := &sql.TxOptions{
		Isolation: isoLevel,
		ReadOnly:  readonly,
	}

	tx, err := t.Db.BeginTx(context.Background(), pt_opt)
	if err != nil {
		return err
	}

	err = fnCallback(tx)
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
