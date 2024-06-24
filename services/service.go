package services

import (
	"dataflow/models"
	"dataflow/repo"
	"errors"
	"math/big"
	"time"
)

var ErrWrongDate = errors.New("start date must be before end date")

type DataService struct {
	repo repo.Repository
}

func NewDataService(repo *repo.InMemoryRepository) *DataService {
	return &DataService{repo}
}

func (ds *DataService) GetAllSales() ([]*models.Sale, error) {
	return ds.repo.GetAllSales()
}

func (ds *DataService) AddSale(sale *models.Sale) error {
	err := ds.repo.AddSale(sale)
	if err != nil {
		return err
	}
	return nil
}

func (ds *DataService) CalculateSales(startDate time.Time, endDate time.Time, storeId string) (*big.Float, error) {
	if startDate.After(endDate) {
		return new(big.Float).SetFloat64(0.0), ErrWrongDate
	}
	sales, err := ds.repo.CalculateSales(startDate, endDate, storeId)
	return sales, err
}
