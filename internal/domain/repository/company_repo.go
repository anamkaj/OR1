package repository

import (
	"ord_crm/internal/models"
)

type CompanyRepository interface {
	GetClientList(user_id string, page int) ([]models.CompanyList, error)
	GetClientLot(company_id int) ([]string, error)
	GetClientBills(company_id int) ([]models.BillsList, error)
}
