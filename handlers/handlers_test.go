package handlers

import (
	"bytes"
	"dataflow/models"
	"dataflow/repo"
	"dataflow/services"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func setupHandler() *DataHandler {
	repository := repo.NewInMemoryRepository()
	service := services.NewDataService(repository)
	handler := NewDataHandler(service)
	return handler
}

func TestDataHandler_GetData(t *testing.T) {
	handler := setupHandler()

	sale1 := &models.Sale{
		ID:           "1",
		ProductId:    "12345",
		StoreId:      "6789",
		QuantitySold: 10,
		SalePrice:    19.99,
		SaleDate:     time.Date(2024, 6, 15, 14, 30, 0, 0, time.UTC),
	}
	sale2 := &models.Sale{
		ID:           "2",
		ProductId:    "54321",
		StoreId:      "9876",
		QuantitySold: 5,
		SalePrice:    9.99,
		SaleDate:     time.Date(2024, 6, 16, 10, 0, 0, 0, time.UTC),
	}

	handler.service.AddSale(sale1)
	handler.service.AddSale(sale2)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler.GetData(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var sales []*models.Sale
	err := json.Unmarshal(w.Body.Bytes(), &sales)
	assert.NoError(t, err)
	assert.Len(t, sales, 2)
	assert.Contains(t, sales, sale1)
	assert.Contains(t, sales, sale2)
}

func TestDataHandler_AddData(t *testing.T) {
	handler := setupHandler()

	sale := &models.Sale{
		ProductId:    "12345",
		StoreId:      "6789",
		QuantitySold: 10,
		SalePrice:    19.99,
		SaleDate:     time.Date(2024, 6, 15, 14, 30, 0, 0, time.UTC),
	}

	jsonData, _ := json.Marshal(sale)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/data", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.AddData(c)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestDataHandler_Calculate(t *testing.T) {
	handler := setupHandler()

	sale1 := &models.Sale{
		ID:           "1",
		ProductId:    "12345",
		StoreId:      "6789",
		QuantitySold: 10,
		SalePrice:    19.99,
		SaleDate:     time.Date(2024, 6, 15, 14, 30, 0, 0, time.UTC),
	}
	sale2 := &models.Sale{
		ID:           "2",
		ProductId:    "54321",
		StoreId:      "9876",
		QuantitySold: 5,
		SalePrice:    9.99,
		SaleDate:     time.Date(2024, 6, 16, 10, 0, 0, 0, time.UTC),
	}

	handler.service.AddSale(sale1)
	handler.service.AddSale(sale2)

	calculateRequest := CalculateRequest{
		Operation: "total_sales",
		StoreId:   "6789",
		StartDate: time.Date(2024, 6, 1, 14, 30, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 6, 16, 14, 30, 0, 0, time.UTC),
	}

	jsonData, _ := json.Marshal(calculateRequest)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/calculate", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Calculate(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var calculateResponse CalculateResponse
	err := json.Unmarshal(w.Body.Bytes(), &calculateResponse)
	assert.NoError(t, err)
	assert.Equal(t, "6789", calculateResponse.StoreId)
	assert.NotNil(t, calculateResponse.TotalSales)
}

func TestDataHandler_AddData_InvalidJSON(t *testing.T) {
	handler := setupHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/data", bytes.NewBufferString("invalid json"))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.AddData(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDataHandler_Calculate_InvalidJSON(t *testing.T) {
	handler := setupHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/calculate", bytes.NewBufferString("invalid json"))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Calculate(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDataHandler_Calculate_InvalidOperation(t *testing.T) {
	handler := setupHandler()

	calculateRequest := CalculateRequest{
		Operation: "invalid_operation",
		StoreId:   "9876",
		StartDate: time.Date(2024, 6, 1, 14, 30, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 6, 16, 14, 30, 0, 0, time.UTC),
	}

	jsonData, _ := json.Marshal(calculateRequest)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/calculate", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Calculate(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "unsupported operation")
}

func TestDataHandler_Calculate_InvalidDate(t *testing.T) {
	handler := setupHandler()

	calculateRequest := CalculateRequest{
		Operation: "total_sales",
		StoreId:   "9876",
		StartDate: time.Date(2024, 6, 20, 14, 30, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 6, 16, 14, 30, 0, 0, time.UTC),
	}

	jsonData, _ := json.Marshal(calculateRequest)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/calculate", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Calculate(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "start date must be before end date")
}
