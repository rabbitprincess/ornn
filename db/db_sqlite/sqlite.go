package db_sqlite

import "github.com/gokch/ornn/db"

func New(path string) (*db.Conn, error) {
	conn := &db.Conn{}
	err := conn.Connect("sqlite3", path, "")
	if err != nil {
		return nil, err
	}

	return conn, nil
}
