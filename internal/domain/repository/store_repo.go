package repository

import "ord_crm/internal/domain/entity"

// Интерфейс репозитория для работы с актами
type StoreMet interface {
	AddNewAct(data []entity.PivotTable) error
	GetAct(key string) (*entity.PivotTable, error)
}
