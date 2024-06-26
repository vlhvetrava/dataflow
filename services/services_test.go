package services

import (
	"dataflow/models"
	"dataflow/repo"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
	"time"
)

func TestDataService_GetAllSales(t *testing.T) {
	mockRepo := new(repo.MockRepository)
	service := NewDataService(mockRepo)

	sale1 := &models.Sale{
		ProductId:    "12345",
		StoreId:      "6789",
		QuantitySold: 10,
		SalePrice:    19.99,
		SaleDate:     time.Date(2024, 6, 15, 14, 30, 0, 0, time.UTC),
	}
	sale2 := &models.Sale{
		ProductId:    "54321",
		StoreId:      "9876",
		QuantitySold: 5,
		SalePrice:    9.99,
		SaleDate:     time.Date(2024, 6, 16, 10, 0, 0, 0, time.UTC),
	}

	mockRepo.On("GetAllSales").Return([]*models.Sale{sale1, sale2}, nil)

	sales, err := service.GetAllSales()

	assert.Nil(t, err)
	assert.Equal(t, 2, len(sales))
	assert.Contains(t, sales, sale1)
	assert.Contains(t, sales, sale2)
}

func TestDataService_AddSale(t *testing.T) {
	mockRepo := new(repo.MockRepository)
	service := NewDataService(mockRepo)

	sale := &models.Sale{
		ProductId:    "12345",
		StoreId:      "6789",
		QuantitySold: 10,
		SalePrice:    19.99,
		SaleDate:     time.Date(2024, 6, 15, 14, 30, 0, 0, time.UTC),
	}

	mockRepo.On("AddSale", sale).Return(nil)

	err := service.AddSale(sale)
	assert.Nil(t, err)
}

func TestDataService_CalculateSales(t *testing.T) {
	mockRepo := new(repo.MockRepository)
	service := NewDataService(mockRepo)

	sale1 := &models.Sale{
		ProductId:    "12345",
		StoreId:      "6789",
		QuantitySold: 10,
		SalePrice:    19.99,
		SaleDate:     time.Date(2024, 6, 15, 14, 30, 0, 0, time.UTC),
	}

	startDate := time.Date(2024, 6, 1, 14, 30, 0, 0, time.UTC)
	endDate := time.Date(2024, 6, 16, 14, 30, 0, 0, time.UTC)
	storeId := "6789"

	mockRepo.On("GetSalesInRange", startDate, endDate, storeId).Return([]*models.Sale{sale1}, nil)

	expectedTotal := new(big.Float).SetPrec(20).SetFloat64(199.90)

	totalSales, err := service.CalculateSales(startDate, endDate, sale1.StoreId)
	assert.Nil(t, err)
	assert.Equal(t, expectedTotal, totalSales)
}

func TestDataService_GetAllSales_Empty(t *testing.T) {
	mockRepo := new(repo.MockRepository)
	service := NewDataService(mockRepo)

	mockRepo.On("GetAllSales").Return([]*models.Sale{}, nil)

	sales, err := service.GetAllSales()

	assert.Nil(t, err)
	assert.Equal(t, 0, len(sales))
}

func TestDataService_CalculateSales_NoSalesInDateRange(t *testing.T) {
	mockRepo := new(repo.MockRepository)
	service := NewDataService(mockRepo)

	sale1 := &models.Sale{
		ProductId:    "12345",
		StoreId:      "6789",
		QuantitySold: 10,
		SalePrice:    19.99,
		SaleDate:     time.Date(2024, 6, 15, 14, 30, 0, 0, time.UTC),
	}

	startDate := time.Date(2024, 6, 1, 14, 30, 0, 0, time.UTC)
	endDate := time.Date(2024, 6, 10, 14, 30, 0, 0, time.UTC)
	storeId := "6789"

	mockRepo.On("GetSalesInRange", startDate, endDate, storeId).Return([]*models.Sale{}, nil)

	expectedTotal := new(big.Float).SetPrec(20).SetFloat64(0.0)

	totalSales, err := service.CalculateSales(startDate, endDate, sale1.StoreId)
	assert.Nil(t, err)
	assert.Equal(t, expectedTotal, totalSales)
}

func TestDataService_CalculateSales_NoSalesForStore(t *testing.T) {
	mockRepo := new(repo.MockRepository)
	service := NewDataService(mockRepo)

	startDate := time.Date(2024, 6, 1, 14, 30, 0, 0, time.UTC)
	endDate := time.Date(2024, 6, 16, 14, 30, 0, 0, time.UTC)
	storeId := "9876"

	mockRepo.On("GetSalesInRange", startDate, endDate, storeId).Return([]*models.Sale{}, nil)

	expectedTotal := new(big.Float).SetPrec(20).SetFloat64(0.0)

	totalSales, err := service.CalculateSales(startDate, endDate, storeId)
	assert.Nil(t, err)
	assert.Equal(t, expectedTotal, totalSales)
}

func TestDataService_CalculateSales_ZeroStartDate(t *testing.T) {
	mockRepo := new(repo.MockRepository)
	service := NewDataService(mockRepo)

	endDate := time.Date(2024, 6, 16, 14, 30, 0, 0, time.UTC)
	storeId := "6789"

	mockRepo.On("GetSalesInRange", time.Time{}, endDate, storeId).Return([]*models.Sale{}, nil)

	expectedTotal := new(big.Float).SetPrec(20).SetFloat64(0.0)

	totalSales, err := service.CalculateSales(time.Time{}, endDate, storeId)
	assert.Nil(t, err)
	assert.Equal(t, expectedTotal, totalSales)
}

func TestDataService_CalculateSales_ZeroEndDate(t *testing.T) {
	mockRepo := new(repo.MockRepository)
	service := NewDataService(mockRepo)

	startDate := time.Date(2024, 6, 1, 14, 30, 0, 0, time.UTC)
	storeId := "6789"

	mockRepo.On("GetSalesInRange", startDate, time.Time{}, storeId).Return([]*models.Sale{}, nil)

	expectedTotal := new(big.Float).SetPrec(20).SetFloat64(0.0)

	totalSales, err := service.CalculateSales(startDate, time.Time{}, storeId)
	assert.Nil(t, err)
	assert.Equal(t, expectedTotal, totalSales)
}

func TestDataService_CalculateSales_ZeroStartDateAndEndDate(t *testing.T) {
	mockRepo := new(repo.MockRepository)
	service := NewDataService(mockRepo)

	storeId := "6789"

	mockRepo.On("GetSalesInRange", time.Time{}, time.Time{}, storeId).Return([]*models.Sale{}, nil)

	expectedTotal := new(big.Float).SetPrec(20).SetFloat64(0.0)

	totalSales, err := service.CalculateSales(time.Time{}, time.Time{}, storeId)
	assert.Nil(t, err)
	assert.Equal(t, expectedTotal, totalSales)
}
