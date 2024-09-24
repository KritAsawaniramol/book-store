package server

import (
	"github.com/kritAsawaniramol/book-store/module/shelf/shelfHandler"
	"github.com/kritAsawaniramol/book-store/module/shelf/shelfRepository"
	"github.com/kritAsawaniramol/book-store/module/shelf/shelfUsecase"
)

func (g *ginServer) shelfService() {
	repo := shelfRepository.NewUserRepositoryImpl(g.db)
	usecase := shelfUsecase.NewShelfUsecaseImpl(repo)
	httpHandler := shelfHandler.NewShelfHttpHandlerImpl(g.cfg, usecase)
	queueHandler := shelfHandler.NewShelfQueueHandlerImpl(g.cfg, usecase)

	_ = httpHandler
	_ = queueHandler
}
