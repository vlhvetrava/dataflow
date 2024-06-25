package services

import (
	"dataflow/models"
	"github.com/stretchr/testify/mock"
	"math/big"
	"time"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) GetAllSales() ([]*models.Sale, error) {
	args := m.Called()
	return args.Get(0).([]*models.Sale), args.Error(1)
}

func (m *MockService) AddSale(sale *models.Sale) error {
	args := m.Called(sale)
	return args.Error(0)
}

func (m *MockService) CalculateSales(startDate time.Time, endDate time.Time, storeId string) (*big.Float, error) {
	args := m.Called(startDate, endDate, storeId)
	return args.Get(0).(*big.Float), args.Error(1)
}
