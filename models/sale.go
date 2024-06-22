package models

import "time"

type Sale struct {
	ID           string    `json:"id"`
	ProductId    string    `json:"product_id"`
	StoreId      string    `json:"store_id"`
	QuantitySold int       `json:"quantity_sold"`
	SalePrice    float64   `json:"sale_price"`
	SaleDate     time.Time `json:"sale_date"`
}
