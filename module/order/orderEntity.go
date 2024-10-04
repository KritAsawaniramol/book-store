package order

import (
	"gorm.io/gorm"
)

type (
	Orders struct {
		gorm.Model
		UserID uint
		Status string
		Books  []OrderBook
		Note   string
		Total  uint
	}

	OrderBook struct {
		gorm.Model
		OrdersID uint
		BookID   uint
		Price    uint
	}
)
