package database

import (
	"context"
	"encoding/json"
	"fmt"
	"ord_crm/internal/scraper"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type Store struct {
	db *redis.Client
}

func NewStoreRepo(d *redis.Client) *Store {
	return &Store{
		db: d,
	}
}

func (s *Store) AddNewAct(data []scraper.PivotTable) error {

	for _, r := range data {

		json, err := json.Marshal(r)
		if err != nil {
			fmt.Println("error serialization ..", err)
			return err
		}

		err = s.db.Set(ctx, r.Acts.StatementNumber, json, 0).Err()
		if err != nil {
			fmt.Println("error add strings ..", err)
			return err
		}
	}

	return nil
}

func (s *Store) GetAct(key string) error {

	val, err := s.db.Get(ctx, "name").Result()
	if err != nil {
		fmt.Println("error add strings ..")
		return err
	}

	fmt.Println(val)
	return nil
}
