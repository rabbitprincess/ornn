// Code generated - DO NOT EDIT.
// This file is a generated and any changes will be lost.

package gen

import (
	"fmt"
	. "github.com/gokch/ornn/db"
)

type Gen struct {
	Newtable Newtable
}

func (t *Gen) Init(
	job *Job,
) {
	t.Newtable.Init(job)
}

func (t *Newtable) Init(
	job *Job,
) {
	t.job = job
}

type Newtable struct {
	job *Job
}

func (t *Newtable) Insert(
	val_a string,
	val_b string,
	val_seq int64,
) (
	lastInsertId int64,
	err error,
) {
	args := []interface{}{
		val_a,
		val_b,
		val_seq,
	}
	
	sql := fmt.Sprintf(
		"INSERT INTO newtable VALUES ($1, $2, $3)",
	)
	
	exec, err := t.job.Exec(
		sql,
		args...,
	)
	if err != nil {
		return 0, err
	}
	
	return exec.LastInsertId()
}

type Newtable_select struct {
	A   string
	B   string
	Seq int64
}

func (t *Newtable) Select(
	where_seq int64,
) (
	selects []*Newtable_select,
	err error,
) {
	args := []interface{}{
		where_seq,
	}
	
	sql := fmt.Sprintf(
		"SELECT * FROM newtable WHERE seq = $1",
	)
	ret, err := t.job.Query(
		sql,
		args...,
	)
	if err != nil {
		return nil, err
	}
	defer ret.Close()
	
	selects = make([]*Newtable_select, 0, 100)
	for ret.Next() {
		scan := &Newtable_select{}
		err := ret.Scan(scan)
		if err != nil {
			return nil, err
		}
		selects = append(selects, scan)
	}
	
	return selects, nil
}

func (t *Newtable) Delete(
	where_seq int64,
) (
	rowAffected int64,
	err error,
) {
	args := []interface{}{
		where_seq,
	}
	
	sql := fmt.Sprintf(
		"DELETE FROM newtable WHERE seq = $1",
	)
			
	exec, err := t.job.Exec(
		sql,
		args...,
	)
	if err != nil {
		return 0, err
	}
	
	return exec.RowsAffected()
}

func (t *Newtable) Update(
	val_a string,
	val_b string,
	val_seq int64,
	where_seq int64,
) (
	rowAffected int64,
	err error,
) {
	sql := fmt.Sprintf(
		"UPDATE newtable SET a = $1, b = $2, seq = $3 WHERE seq = $4",
	)
	args := []interface{}{
		val_a,
		val_b,
		val_seq,
		where_seq,
	}
	
	exec, err := t.job.Exec(
		sql,
		args...,
	)
	if err != nil {
		return 0, err
	}
	
	return exec.RowsAffected()
}

