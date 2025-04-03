package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AuthToken struct {
	Token        string
	Session      string
	RedisPass    string
	PostgreSQL   string
	BaseURL      string
	SearchURL    string
	Specialist_1 string
	Specialist_2 string
	BillsURL        string
}

func GetToken() (AuthToken, error) {
	err := godotenv.Load("./../.env")
	if err != nil {
		log.Fatalln("Error loading .env file")
		return AuthToken{}, err
	}

	access_token := os.Getenv("ACCESS_TOKEN")
	id_session := os.Getenv("ID_SESSION")
	rd_pass := os.Getenv("REDIS_PASS")
	db_postgres := os.Getenv("CLIENT_TABLE")
	base_url := os.Getenv("BASEURL")
	search_url := os.Getenv("SEARCHURL")
	specialist_1 := os.Getenv("SPECIALIST_1")
	specialist_2 := os.Getenv("SPECIALIST_2")
	bills := os.Getenv("BILLSURL")

	token := AuthToken{
		Token:        access_token,
		Session:      id_session,
		RedisPass:    rd_pass,
		PostgreSQL:   db_postgres,
		BaseURL:      base_url,
		SearchURL:    search_url,
		Specialist_1: specialist_1,
		Specialist_2: specialist_2,
		BillsURL:        bills,
	}

	return token, nil
}
