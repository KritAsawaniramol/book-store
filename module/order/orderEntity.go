package order

import (
	"gorm.io/gorm"
)

type (
	Orders struct {
		gorm.Model
		UserID uint
		Status string
		OrdersBooks  []OrdersBooks
		Note   string
		Total  uint
	}

	OrdersBooks struct {
		gorm.Model
		OrdersID uint
		BookID   uint
		//Price when book was buyed
		Price    uint
	}
)
