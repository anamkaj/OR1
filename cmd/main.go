package main

import (
	"fmt"
	"ord_crm/cmd/app"
	pgdb "ord_crm/internal/infrastructure/database/postgres"
	redisdb "ord_crm/internal/infrastructure/database/redis"
)

func main() {

	db, err := redisdb.RedisConnect()
	if err != nil {
		fmt.Println("Error connecting redis", err)
	}

	pg, err := pgdb.PGconnect()
	if err != nil {
		fmt.Println("Error connecting Postgres", err)
	}
	fmt.Println(pg)

	n := app.NewApp(db)

	err = n.App()
	if err != nil {
		fmt.Println("error start New app ....")

	}

}
