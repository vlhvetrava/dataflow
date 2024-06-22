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
	service *services.DataService
}

func NewDataHandler(service *services.DataService) *DataHandler {
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
	Operation string    `json:"operation"`
	StoreId   string    `json:"store_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

type CalculateResponse struct {
	StoreId    string     `json:"store_id"`
	StartDate  time.Time  `json:"start_date"`
	EndDate    time.Time  `json:"end_date"`
	TotalSales *big.Float `json:"total_sales"`
}

func (h *DataHandler) Calculate(c *gin.Context) {
	var calculateRequest CalculateRequest
	err := c.ShouldBindJSON(&calculateRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	totalSales, err := h.service.CalculateSales(calculateRequest.StartDate, calculateRequest.EndDate, calculateRequest.StoreId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
