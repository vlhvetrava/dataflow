package handlers

import (
	"dataflow/models"
	"dataflow/repo"
	"dataflow/services"
	"errors"
	"github.com/gin-gonic/gin"
	"math/big"
	"net/http"
	"time"
)

type DataHandler struct {
	service services.DataService
}

func NewDataHandler(service services.DataService) *DataHandler {
	return &DataHandler{service: service}
}

func (h *DataHandler) GetData(c *gin.Context) {
	sales, err := h.service.GetAllSales()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
	}
	c.JSON(http.StatusOK, sales)
}

func (h *DataHandler) AddData(c *gin.Context) {
	var sale *models.Sale
	err := c.ShouldBindJSON(&sale)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.service.AddSale(sale)
	if err != nil {
		if errors.Is(err, repo.ErrSaleAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error(), "status": http.StatusConflict})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "status": http.StatusInternalServerError})
		}
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "success"})
}

type CalculateRequest struct {
	Operation string `json:"operation"`
	StoreId   string `json:"store_id"`
	StartDate string `json:"start_date,omitempty"`
	EndDate   string `json:"end_date,omitempty"`
}

type CalculateResponse struct {
	StoreId    string     `json:"store_id"`
	StartDate  string     `json:"start_date"`
	EndDate    string     `json:"end_date"`
	TotalSales *big.Float `json:"total_sales"`
}

func (h *DataHandler) Calculate(c *gin.Context) {
	var calculateRequest CalculateRequest
	err := c.ShouldBindJSON(&calculateRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if calculateRequest.Operation != "total_sales" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported operation"})
		return
	}

	var startDate time.Time
	var endDate time.Time

	if calculateRequest.StartDate != "" {
		t, err := time.Parse(time.RFC3339, calculateRequest.StartDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		startDate = t
	}
	if calculateRequest.EndDate != "" {
		t, err := time.Parse(time.RFC3339, calculateRequest.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		endDate = t
	}

	totalSales, err := h.service.CalculateSales(startDate, endDate, calculateRequest.StoreId)
	if err != nil {
		if errors.Is(err, services.ErrWrongDate) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "status": http.StatusInternalServerError})
		}
		return
	}
	calculateResponse := CalculateResponse{
		StoreId:    calculateRequest.StoreId,
		TotalSales: totalSales,
		StartDate:  calculateRequest.StartDate,
		EndDate:    calculateRequest.EndDate,
	}
	c.JSON(http.StatusOK, calculateResponse)
}
