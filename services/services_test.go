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
	repository := repo.NewInMemoryRepository()
	service := NewDataService(repository)

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

	service.AddSale(sale1)
	service.AddSale(sale2)

	sales, err := service.GetAllSales()

	assert.Nil(t, err)
	assert.Equal(t, 2, len(sales))
	assert.Contains(t, sales, sale1)
	assert.Contains(t, sales, sale2)
}

func TestDataService_AddSale(t *testing.T) {
	repository := repo.NewInMemoryRepository()
	service := NewDataService(repository)

	sale := &models.Sale{
		ProductId:    "12345",
		StoreId:      "6789",
		QuantitySold: 10,
		SalePrice:    19.99,
		SaleDate:     time.Date(2024, 6, 15, 14, 30, 0, 0, time.UTC),
	}

	err := service.AddSale(sale)
	assert.Nil(t, err)
}

func TestDataService_CalculateSales(t *testing.T) {
	repository := repo.NewInMemoryRepository()
	service := NewDataService(repository)

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

	service.AddSale(sale1)
	service.AddSale(sale2)

	startDate := time.Date(2024, 6, 1, 14, 30, 0, 0, time.UTC)
	endDate := time.Date(2024, 6, 16, 14, 30, 0, 0, time.UTC)
	expectedTotal := new(big.Float).SetPrec(20).SetFloat64(199.90)

	totalSales, err := service.CalculateSales(startDate, endDate, sale1.StoreId)
	assert.Nil(t, err)
	assert.Equal(t, expectedTotal, totalSales)
}

func TestDataService_GetAllSales_Empty(t *testing.T) {
	repository := repo.NewInMemoryRepository()
	service := NewDataService(repository)

	sales, err := service.GetAllSales()

	assert.Nil(t, err)
	assert.Equal(t, 0, len(sales))
}

func TestDataService_CalculateSales_NoSalesInDateRange(t *testing.T) {
	repository := repo.NewInMemoryRepository()
	service := NewDataService(repository)

	sale1 := &models.Sale{
		ProductId:    "12345",
		StoreId:      "6789",
		QuantitySold: 10,
		SalePrice:    19.99,
		SaleDate:     time.Date(2024, 6, 15, 14, 30, 0, 0, time.UTC),
	}

	service.AddSale(sale1)

	startDate := time.Date(2024, 6, 1, 14, 30, 0, 0, time.UTC)
	endDate := time.Date(2024, 6, 10, 14, 30, 0, 0, time.UTC)
	expectedTotal := new(big.Float).SetPrec(20).SetFloat64(0.0)

	totalSales, err := service.CalculateSales(startDate, endDate, sale1.StoreId)
	assert.Nil(t, err)
	assert.Equal(t, expectedTotal, totalSales)
}

func TestDataService_CalculateSales_NoSalesForStore(t *testing.T) {
	repository := repo.NewInMemoryRepository()
	service := NewDataService(repository)

	sale1 := &models.Sale{
		ProductId:    "12345",
		StoreId:      "6789",
		QuantitySold: 10,
		SalePrice:    19.99,
		SaleDate:     time.Date(2024, 6, 15, 14, 30, 0, 0, time.UTC),
	}

	service.AddSale(sale1)

	startDate := time.Date(2024, 6, 1, 14, 30, 0, 0, time.UTC)
	endDate := time.Date(2024, 6, 16, 14, 30, 0, 0, time.UTC)
	storeId := "9876"
	expectedTotal := new(big.Float).SetPrec(20).SetFloat64(0.0)

	totalSales, err := service.CalculateSales(startDate, endDate, storeId)
	assert.Nil(t, err)
	assert.Equal(t, expectedTotal, totalSales)
}
