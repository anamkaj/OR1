package application

import (
	"ord_crm/internal/domain/entity"
	"ord_crm/internal/domain/repository"
)

type SearchService struct {
	repo repository.Search
}

func NewSearchService(r repository.Search) *SearchService {
	return &SearchService{
		repo: r,
	}
}

func (s *SearchService) ActSearch(id string, bill []string) ([]entity.PivotTable, error) {
	return s.repo.ActSearch(id, bill)
}

func (s *SearchService) Bills(link string) (entity.Bills, error) {
	return s.repo.Bills(link)
}

func (s *SearchService) SearchID(id string) (int, error) {
	return s.repo.SearchID(id)
}
