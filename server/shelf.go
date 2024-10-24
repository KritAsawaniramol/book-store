package server

import (
	"context"
	"log"

	"github.com/kritAsawaniramol/book-store/module/shelf/shelfHandler"
	"github.com/kritAsawaniramol/book-store/module/shelf/shelfPb"
	"github.com/kritAsawaniramol/book-store/module/shelf/shelfRepository"
	"github.com/kritAsawaniramol/book-store/module/shelf/shelfUsecase"
	"github.com/kritAsawaniramol/book-store/pkg/grpccon"
	"github.com/kritAsawaniramol/book-store/pkg/queue"
)

func (g *ginServer) shelfService() {
	repo := shelfRepository.NewUserRepositoryImpl(g.db)
	usecase := shelfUsecase.NewShelfUsecaseImpl(repo)
	httpHandler := shelfHandler.NewShelfHttpHandlerImpl(g.cfg, usecase)
	queueHandler := shelfHandler.NewShelfQueueHandlerImpl(g.cfg, usecase)
	grpcHandler := shelfHandler.NewShelfGrpcHandlerImpl(g.cfg, usecase)

	queueConn, err := queue.ConnectConsumer([]string{g.cfg.Kafka.Url}, g.cfg.Kafka.ApiKey, g.cfg.Kafka.Secret,g.cfg.Kafka.GroupID)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		grpcServer, listener := grpccon.NewGrpcServer(g.cfg.Grpc.ShelfUrl)
		shelfPb.RegisterShelfGrpcServiceServer(grpcServer, grpcHandler)
		log.Printf("gRPC server listening on %s\n", g.cfg.Grpc.ShelfUrl)
		grpcServer.Serve(listener)
	}()

	consumerHandeler := shelfHandler.NewShelfConsumerHandler(usecase, g.cfg)
	go func() {
		for {
			if err := queueConn.Consume(context.Background(), []string{"shelf"}, consumerHandeler); err != nil {
				log.Fatal(err)
			}
		}
	}()

	shelf := g.app.Group("/shelf_v1")
	
	
	shelf.GET("", g.healthCheck)
	shelf.GET("/shelf", g.middleware.JwtAuthorization(), httpHandler.GetMyShelf)

	_ = queueHandler
}
