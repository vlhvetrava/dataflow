package repo

import (
	"dataflow/models"
	"errors"
	"github.com/google/uuid"
	"math/big"
	"sync"
	"time"
)

var ErrSaleAlreadyExists = errors.New("sale already exists")

type Repository interface {
	AddSale(sale *models.Sale) error
	GetAllSales() ([]*models.Sale, error)
	CalculateSales(startDate time.Time, endDate time.Time, storeId string) (*big.Float, error)
}

type InMemoryRepository struct {
	data sync.Map
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{}
}

func (repo *InMemoryRepository) GetAllSales() ([]*models.Sale, error) {
	var sales []*models.Sale
	repo.data.Range(func(k, v interface{}) bool {
		sales = append(sales, v.(*models.Sale))
		return true
	})
	return sales, nil
}

func (repo *InMemoryRepository) AddSale(sale *models.Sale) error {
	sale.ID = uuid.New().String()
	_, loaded := repo.data.LoadOrStore(sale.ID, sale)
	if loaded {
		return ErrSaleAlreadyExists
	}
	return nil
}

func (repo *InMemoryRepository) CalculateSales(startDate time.Time, endDate time.Time, storeId string) (*big.Float, error) {
	totalSales := new(big.Float).SetPrec(20).SetFloat64(0.0)
	repo.data.Range(func(k, v interface{}) bool {
		sale := v.(*models.Sale)
		if storeId == sale.StoreId && sale.SaleDate.Before(endDate) && sale.SaleDate.After(startDate) {
			quantityBigFloat := new(big.Float).SetPrec(64).SetInt64(int64(sale.QuantitySold))
			priceBigFloat := new(big.Float).SetPrec(64).SetFloat64(sale.SalePrice)
			saleAmount := new(big.Float).Mul(quantityBigFloat, priceBigFloat)
			totalSales.Add(totalSales, saleAmount)
		}
		return true
	})
	return totalSales, nil
}
