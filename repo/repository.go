package repo

import (
	"dataflow/models"
	"errors"
	"github.com/google/uuid"
	"sync"
	"time"
)

var ErrSaleAlreadyExists = errors.New("sale already exists")

type Repository interface {
	AddSale(sale *models.Sale) error
	GetAllSales() ([]*models.Sale, error)
	GetSalesInRange(startDate time.Time, endDate time.Time, storeId string) ([]*models.Sale, error)
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

func (repo *InMemoryRepository) GetSalesInRange(startDate time.Time, endDate time.Time, storeId string) ([]*models.Sale, error) {
	var sales []*models.Sale
	repo.data.Range(func(k, v interface{}) bool {
		sale := v.(*models.Sale)
		if storeId == sale.StoreId && sale.SaleDate.Before(endDate) && sale.SaleDate.After(startDate) {
			sales = append(sales, sale)
		}
		return true
	})
	return sales, nil
}
