package server

import (
	"github.com/kritAsawaniramol/book-store/module/payment/paymentHandler"
	"github.com/kritAsawaniramol/book-store/module/payment/paymentRepository"
	"github.com/kritAsawaniramol/book-store/module/payment/paymentUsecase"
)

func (g *ginServer) paymentServer()  {
	repo := paymentRepository.NewPaymentRepositoryImpl(g.db)
	usecase := paymentUsecase.NewPaymentUsecaseImpl(repo)
	httpHandler := paymentHandler.NewPaymentHttpHandlerImpl(usecase)

	_ = httpHandler
}