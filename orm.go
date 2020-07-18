package orm

import "database/sql"

type DB struct {
	db     *sql.DB
	models []*Model
}

func Open(driver, datasource string) (*DB, error) {
	db, err := sql.Open(driver, datasource)
	if err != nil {
		return nil, err
	}
	return &DB{db: db}, nil
}
