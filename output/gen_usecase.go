package gen

import (
	"database/sql"

	"github.com/gokch/ornn/db"
)

func NewLogic(db *db.Conn) *Logic {
	return &Logic{
		db: db,
	}
}

type Logic struct {
	db *db.Conn
}

func (t *Logic) WithoutTx() error {
	gen := &Schema{}
	gen.Init(t.db.Job())

	// write func
	_, err := gen.Test.Select()
	if err != nil {
		return err
	}
	return nil
}

func (t *Logic) WithTx() error {
	gen := &Schema{}
	job, err := t.db.TxJob(sql.LevelSerializable, false)
	if err != nil {
		return err
	}
	gen.Init(job)

	// write func
	_, err = gen.Test.Select()

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

func (t *Logic) WithTxFunc() error {
	return t.db.TxJobFunc(sql.LevelSerializable, false, func(job *db.Job) error {
		gen := &Schema{}
		gen.Init(job)
		// write func
		_, err := gen.Test.Select()
		if err != nil {
			return err
		}
		return nil
	})
}
