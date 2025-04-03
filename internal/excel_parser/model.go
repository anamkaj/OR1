package excel

import "ord_crm/internal/scraper"

type ExcelRepository interface {
	ExcelParse() ([]Excel, error)
	ExcelImport(data [][]scraper.PivotTable) error
}

type Excel struct {
	ID    string
	Login string
	Count int
}
