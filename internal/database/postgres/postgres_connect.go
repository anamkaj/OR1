package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
	"log"
	"ord_crm/internal/config"
)

func PGconnect() (*sqlx.DB, error) {
	token, err := config.GetToken()
	if err != nil {
		fmt.Printf("Error: %s", err)
		return nil, err
	}

	db, err := sqlx.Connect("postgres", token.PostgreSQL)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		fmt.Printf("error connect PG %s", err)
		return nil, err

	}

	fmt.Println("PG connected")

	return db, nil

}
