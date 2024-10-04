package server

import (
	"context"
	"log"

	"github.com/kritAsawaniramol/book-store/module/shelf/shelfHandler"
	"github.com/kritAsawaniramol/book-store/module/shelf/shelfRepository"
	"github.com/kritAsawaniramol/book-store/module/shelf/shelfUsecase"
	"github.com/kritAsawaniramol/book-store/pkg/queue"
)

func (g *ginServer) shelfService() {
	repo := shelfRepository.NewUserRepositoryImpl(g.db)
	usecase := shelfUsecase.NewShelfUsecaseImpl(repo)
	httpHandler := shelfHandler.NewShelfHttpHandlerImpl(g.cfg, usecase)
	queueHandler := shelfHandler.NewShelfQueueHandlerImpl(g.cfg, usecase)
	queueConn, err := queue.ConnectConsumer([]string{g.cfg.Kafka.Url}, g.cfg.Kafka.GroupID)
	if err != nil {
		log.Fatal(err)
	}

	consumerHandeler := shelfHandler.NewShelfConsumerHandler(usecase, g.cfg)
	go func() {
		for {
			if err := queueConn.Consume(context.Background(), []string{"shelf"}, consumerHandeler); err != nil {
				log.Fatal(err)
			}
		}
	}()

	// go func() {
	// 	err := <-queueConn.Errors()
	// 	log.Printf("error: queueConn: %s\n", err.Error())
	// }()

	_ = httpHandler
	_ = queueHandler
}
