package application

import (
	"ord_crm/internal/domain/entity"
	database "ord_crm/internal/infrastructure/database/redis"
)

type StoreService struct {
	store *database.Store
}

// Конструктор для создания сервиса
func NewStoreService(store *database.Store) *StoreService {
	return &StoreService{
		store: store,
	}
}

// Пример использования метода репозитория
func (s *StoreService) AddNewAct(data []entity.PivotTable) error {
	return s.store.AddNewAct(data)
}

func (s *StoreService) GetAct(key string) error {
	return s.store.GetAct(key)
}
