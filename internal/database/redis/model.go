package database

import "ord_crm/internal/scraper"

type StoreMet interface {
	AddNewAct(data []scraper.PivotTable) error
	GetAct(key string) (*scraper.PivotTable, error)
}
