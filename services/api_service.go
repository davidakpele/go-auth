package services

import (
	"api-service/repositories"
)

type APIService struct {
	apiRepo repositories.APIRepository
}

func NewAPIService(apiRepo repositories.APIRepository) *APIService {
	return &APIService{apiRepo: apiRepo}
}

func (s *APIService) GetProductById(id int) (interface{}, error) {
	return s.apiRepo.GetProductById(id)
}

