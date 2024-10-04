package server

import (
	"context"
	"log"

	"github.com/kritAsawaniramol/book-store/module/order/orderHandler"
	"github.com/kritAsawaniramol/book-store/module/order/orderRepository"
	"github.com/kritAsawaniramol/book-store/module/order/orderUsecase"
	"github.com/kritAsawaniramol/book-store/pkg/queue"
)

func (g *ginServer) orderServer() {
	repo := orderRepository.NewOrderRepositoryImpl(g.db)
	usecase := orderUsecase.NewOrderUsecaseImpl(repo)
	httpHandler := orderHandler.NewOrderHttpHandlerImpl(g.cfg, usecase)
	queueConn, err := queue.ConnectConsumer([]string{g.cfg.Kafka.Url}, g.cfg.Kafka.GroupID)
	if err != nil {
		log.Fatal(err)
	}

	consumerHandler := orderHandler.NewOrderConsumerHandler(usecase, g.cfg)
	go func() {
		for {
			if err := queueConn.Consume(context.Background(), []string{"order"}, consumerHandler); err != nil {
				log.Fatal(err)
			}
		}
	}()

	// go func() {
	// 	err := <-queueConn.Errors()
	// 	log.Printf("error: queueConn: %s\n", err.Error())
	// }()

	// consumerHandler := orderHandler.NewOrderHttpHandlerImpl()

	g.app.POST("/order/buy", g.middleware.JwtAuthorization(), httpHandler.BuyBooks)
}
