package services

import (
	"dataflow/models"
	"dataflow/repo"
	"errors"
	"math/big"
	"time"
)

var ErrWrongDate = errors.New("start date must be before end date")

type DataService interface {
	GetAllSales() ([]*models.Sale, error)
	AddSale(sale *models.Sale) error
	CalculateSales(startDate time.Time, endDate time.Time, storeId string) (*big.Float, error)
}

type dataService struct {
	repo repo.Repository
}

func NewDataService(repo repo.Repository) DataService {
	return &dataService{repo}
}

func (ds *dataService) GetAllSales() ([]*models.Sale, error) {
	return ds.repo.GetAllSales()
}

func (ds *dataService) AddSale(sale *models.Sale) error {
	err := ds.repo.AddSale(sale)
	if err != nil {
		return err
	}
	return nil
}

func (ds *dataService) CalculateSales(startDate time.Time, endDate time.Time, storeId string) (*big.Float, error) {
	if startDate.After(endDate) {
		return new(big.Float).SetFloat64(0.0), ErrWrongDate
	}
	sales, err := ds.repo.GetSalesInRange(startDate, endDate, storeId)
	if err != nil {
		return nil, err
	}

	totalSales := new(big.Float).SetPrec(20).SetFloat64(0.0)
	for _, sale := range sales {
		quantityBigFloat := new(big.Float).SetPrec(64).SetInt64(int64(sale.QuantitySold))
		priceBigFloat := new(big.Float).SetPrec(64).SetFloat64(sale.SalePrice)
		saleAmount := new(big.Float).Mul(quantityBigFloat, priceBigFloat)
		totalSales.Add(totalSales, saleAmount)
	}
	return totalSales, err
}
