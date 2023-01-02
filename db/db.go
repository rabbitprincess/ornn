package db

import (
	"database/sql"
)

type Conn struct {
	DriverName string
	Dsn        string
	DbName     string

	db *sql.DB
}

func (t *Conn) Connect(driverName, dsn, dbName string) (err error) {
	t.DriverName = driverName
	t.Dsn = dsn
	t.DbName = dbName

	t.db, err = sql.Open(t.DriverName, t.Dsn)
	if err != nil {
		return err
	}

	err = t.db.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (t *Conn) Raw() *sql.DB {
	return t.db
}

func (t *Conn) SetOpenConns(openConns, idleConns int) {
	if openConns > 0 {
		t.db.SetMaxOpenConns(openConns)
	}
	if idleConns > 0 {
		t.db.SetMaxIdleConns(idleConns)
	}
}

func (t *Conn) Job() *Job {
	job := NewJob(t.db)
	return job
}

func (t *Conn) TxJob(isoLevel sql.IsolationLevel, readonly bool) (job *Job, err error) {
	job = t.Job()
	err = job.BeginTx(isoLevel, readonly)
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (t *Conn) TxJobFunc(isoLevel sql.IsolationLevel, readonly bool, fn func(*Job) error) (err error) {
	job := NewJob(t.db)
	err = job.BeginTx(isoLevel, readonly)
	if err != nil {
		return err
	}

	err = fn(job)
	if err != nil {
		job.Rollback()
		return err
	}
	err = job.Commit()
	if err != nil {
		return err
	}

	return nil
}
