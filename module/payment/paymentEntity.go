package payment

import "gorm.io/gorm"

type (
	PaymentQueue struct {
		gorm.Model
		Offset int64
	}
)
