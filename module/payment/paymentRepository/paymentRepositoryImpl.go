package paymentRepository

import "gorm.io/gorm"

type paymentRepositoryImpl struct {
	db *gorm.DB
}

func NewPaymentRepositoryImpl(db *gorm.DB) PaymentRepository {
	return &paymentRepositoryImpl{db: db}
}