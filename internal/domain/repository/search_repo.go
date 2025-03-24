package repository

import "ord_crm/internal/domain/entity"

type Search interface {
	SearchID(id string) (int, error)
	ActSearch(id string, bill []string) ([]entity.PivotTable, error)
	Bills(link string) (entity.Bills, error)
}
