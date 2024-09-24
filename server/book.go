package server

import (
	"github.com/kritAsawaniramol/book-store/module/book/bookHandler"
	"github.com/kritAsawaniramol/book-store/module/book/bookRepository"
	"github.com/kritAsawaniramol/book-store/module/book/bookUsecase"
)

func (g *ginServer) bookService() {
	repo := bookRepository.NewBookRepositoryImpl(g.db)
	usecase := bookUsecase.NewBookUsecaseImpl(repo)
	httpHandler := bookHandler.NewBookHttpHandlerImpl(g.cfg, usecase)
	grpcHandler := bookHandler.NewBookGrpcHandlerImpl(g.cfg, usecase)

	_ = httpHandler
	_ = grpcHandler
}
