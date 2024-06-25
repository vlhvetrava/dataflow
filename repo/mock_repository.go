package repo

import (
	"dataflow/models"
	"github.com/stretchr/testify/mock"
	"time"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) AddSale(sale *models.Sale) error {
	args := m.Called(sale)
	return args.Error(0)
}

func (m *MockRepository) GetAllSales() ([]*models.Sale, error) {
	args := m.Called()
	return args.Get(0).([]*models.Sale), args.Error(1)
}

func (m *MockRepository) GetSalesInRange(startDate time.Time, endDate time.Time, storeId string) ([]*models.Sale, error) {
	args := m.Called(startDate, endDate, storeId)
	return args.Get(0).([]*models.Sale), args.Error(1)
}
