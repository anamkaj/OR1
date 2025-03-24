package application

import (
	"ord_crm/internal/domain/entity"
	"ord_crm/internal/domain/repository"
)

type ExcelService struct {
	repo repository.ExcelRepository
}

func NewExcelService(e repository.ExcelRepository) *ExcelService {
	return &ExcelService{
		repo: e,
	}
}

func (e *ExcelService) ExcelParse() ([]entity.Excel, error) {
	return e.repo.ExcelParse()
}

func (e *ExcelService) ExcelImport(data [][]entity.PivotTable) error {
	return e.repo.ExcelImport(data)
}
