package repository

import "ord_crm/internal/domain/entity"

type ExcelRepository interface {
	ExcelParse() ([]entity.Excel, error)
	ExcelImport(data [][]entity.PivotTable) error
}
