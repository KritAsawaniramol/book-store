package paymentUsecase

import "github.com/kritAsawaniramol/book-store/module/payment/paymentRepository"

type paymentUsecaseImpl struct {
	paymentRepository paymentRepository.PaymentRepository
}

func NewPaymentUsecaseImpl(paymentRepository paymentRepository.PaymentRepository) PaymentUsecase {
	return &paymentUsecaseImpl{paymentRepository: paymentRepository}
}