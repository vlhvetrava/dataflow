package repo

import (
	"dataflow/models"
	"github.com/stretchr/testify/assert"
	"math/big"
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

func TestInMemoryRepository_CalculateSales(t *testing.T) {
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

	startDate := time.Date(2024, 6, 1, 14, 30, 0, 0, time.UTC)
	endDate := time.Date(2024, 6, 16, 14, 30, 0, 0, time.UTC)
	expectedTotal := new(big.Float).SetPrec(20).SetFloat64(199.90)

	totalSales, err := repo.CalculateSales(startDate, endDate, sale1.StoreId)
	assert.Nil(t, err)
	assert.Equal(t, expectedTotal, totalSales)
}