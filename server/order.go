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
	queueConn, err := queue.ConnectConsumer([]string{g.cfg.Kafka.Url}, g.cfg.Kafka.ApiKey, g.cfg.Kafka.Secret, g.cfg.Kafka.GroupID)
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

	order := g.app.Group("/order_v1")
	
	order.GET("", g.healthCheck)
	order.POST("/order/buy",
		g.middleware.JwtAuthorization(),
		g.middleware.RbacAuthorization(
			map[uint]bool{
				2: true,
			},
		),
		httpHandler.BuyBooks,
	)
	order.GET("/order/myorder",
		g.middleware.JwtAuthorization(),
		g.middleware.RbacAuthorization(
			map[uint]bool{
				2: true,
			},
		),
		httpHandler.SearchOneMyOrder,
	)
	order.GET("/order",
		g.middleware.JwtAuthorization(),
		g.middleware.RbacAuthorization(
			map[uint]bool{
				2: true,
			},
		),
		httpHandler.GetMyOrders,
	)
}
