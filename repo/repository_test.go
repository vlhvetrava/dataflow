package repo

import (
	"dataflow/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestInMemoryRepository_GetAllSales(t *testing.T) {
	repo := NewInMemoryRepository()

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

	repo.AddSale(sale1)
	repo.AddSale(sale2)

	sales, err := repo.GetAllSales()

	assert.Nil(t, err)
	assert.Equal(t, 2, len(sales))
	assert.Contains(t, sales, sale1)
	assert.Contains(t, sales, sale2)
}

func TestInMemoryRepository_AddSale(t *testing.T) {
	repo := NewInMemoryRepository()

	sale := models.Sale{
		ProductId:    "12345",
		StoreId:      "6789",
		QuantitySold: 10,
		SalePrice:    19.99,
		SaleDate:     time.Date(2024, 6, 15, 14, 30, 0, 0, time.UTC),
	}

	err := repo.AddSale(&sale)
	assert.Nil(t, err)
}

func TestInMemoryRepository_GetSalesInRange(t *testing.T) {
	repo := NewInMemoryRepository()

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
		SaleDate:     time.Date(2024, 6, 30, 10, 0, 0, 0, time.UTC),
	}

	repo.AddSale(sale1)
	repo.AddSale(sale2)

	startDate := time.Date(2024, 6, 1, 14, 30, 0, 0, time.UTC)
	endDate := time.Date(2024, 6, 16, 14, 30, 0, 0, time.UTC)

	sales, err := repo.GetSalesInRange(startDate, endDate, sale1.StoreId)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(sales))
	assert.Contains(t, sales, sale1)
}

func TestInMemoryRepository_GetSalesInRange_ZeroStartDate(t *testing.T) {
	repo := NewInMemoryRepository()

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
		SaleDate:     time.Date(2024, 6, 30, 10, 0, 0, 0, time.UTC),
	}

	repo.AddSale(sale1)
	repo.AddSale(sale2)

	endDate := time.Date(2024, 6, 16, 14, 30, 0, 0, time.UTC)

	sales, err := repo.GetSalesInRange(time.Time{}, endDate, sale1.StoreId)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(sales))
	assert.Contains(t, sales, sale1)
}

func TestInMemoryRepository_GetSalesInRange_ZeroEndDate(t *testing.T) {
	repo := NewInMemoryRepository()

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
		SaleDate:     time.Date(2024, 6, 30, 10, 0, 0, 0, time.UTC),
	}

	repo.AddSale(sale1)
	repo.AddSale(sale2)

	startDate := time.Date(2024, 6, 1, 14, 30, 0, 0, time.UTC)

	sales, err := repo.GetSalesInRange(startDate, time.Time{}, sale1.StoreId)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(sales))
	assert.Contains(t, sales, sale1)
}

func TestInMemoryRepository_GetSalesInRange_ZeroStartDateAndEndDate(t *testing.T) {
	repo := NewInMemoryRepository()

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
		SaleDate:     time.Date(2024, 6, 30, 10, 0, 0, 0, time.UTC),
	}

	repo.AddSale(sale1)
	repo.AddSale(sale2)

	sales, err := repo.GetSalesInRange(time.Time{}, time.Time{}, sale1.StoreId)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(sales))
	assert.Contains(t, sales, sale1)
}
