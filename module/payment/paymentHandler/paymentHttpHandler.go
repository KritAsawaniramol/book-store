package paymentHandler

import "github.com/kritAsawaniramol/book-store/module/payment/paymentUsecase"

type (
	PaymentHttpHandler interface {
	}

	paymentHttpHandlerImpl struct {
		paymentUsecase paymentUsecase.PaymentUsecase
	}
)

func NewPaymentHttpHandlerImpl(paymentUsecase paymentUsecase.PaymentUsecase) PaymentHttpHandler {
	return &paymentHttpHandlerImpl{
		paymentUsecase: paymentUsecase,
	}
}
