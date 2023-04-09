package db

import "database/sql"

type PGLoader struct {
	defaultConn string
}

func (p PGLoader) Get(connection string) (*sql.DB, error) {
	if connection == "" {
		connection = p.defaultConn
	}
	db, err := sql.Open("postgres", connection)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func NewPGLoader(defaultConn string) *PGLoader {
	return &PGLoader{
		defaultConn: defaultConn,
	}
}
