package database

import "github.com/jmoiron/sqlx"


type Store struct {
	db *sqlx.DB
}

func NewStoreRepo(db *sqlx.DB) *Store {
	return &Store{
		db: db,
	}
}


