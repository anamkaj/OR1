package application

import (
	"ord_crm/internal/domain/repository"
	"ord_crm/internal/models"
)

type CompanyService struct {
	repo repository.CompanyRepository
}

func NewCompanyService(r repository.CompanyRepository) *CompanyService {
	return &CompanyService{
		repo: r,
	}
}

// GetLot получает список компаний по ID
func (s *CompanyService) GetClientList(user_id string, page int) ([]models.CompanyList, error) {
	return s.repo.GetClientList(user_id, page)
}

func (s *CompanyService) GetClientLot(id int) ([]string, error) {
	return s.repo.GetClientLot(id)
}

func (s *CompanyService) GetClientBills(id int) ([]models.BillsList, error) {
	return s.repo.GetClientBills(id)
}

