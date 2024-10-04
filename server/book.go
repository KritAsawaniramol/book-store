package server

import (
	"log"

	"github.com/kritAsawaniramol/book-store/module/book/bookHandler"
	"github.com/kritAsawaniramol/book-store/module/book/bookPb"
	"github.com/kritAsawaniramol/book-store/module/book/bookRepository"
	"github.com/kritAsawaniramol/book-store/module/book/bookUsecase"
	"github.com/kritAsawaniramol/book-store/pkg/grpccon"
)

func (g *ginServer) bookService() {
	repo := bookRepository.NewBookRepositoryImpl(g.db)
	usecase := bookUsecase.NewBookUsecaseImpl(g.cfg, repo)
	httpHandler := bookHandler.NewBookHttpHandlerImpl(g.cfg, usecase)
	grpcHandler := bookHandler.NewBookGrpcHandlerImpl(g.cfg, usecase)

	go func() {
		grpcServer, listener := grpccon.NewGrpcServer(g.cfg.Grpc.BookUrl)
		bookPb.RegisterBookGrpcServiceServer(grpcServer, grpcHandler)
		log.Printf("gRPC server listening on %s\n", g.cfg.Grpc.BookUrl)
		grpcServer.Serve(listener)
	}()

	// g.app.Static("/book/cover", "./asset/image/bookCover")
	g.app.StaticFile("/book/file", "./asset/book")
	g.app.GET("/book/cover/:fileName", httpHandler.GetBookCoverImage)
	g.app.GET("/book", httpHandler.SearchBooks)
	g.app.GET("/book/:id", httpHandler.GetOneBook)
	g.app.GET("/book/tags", httpHandler.GetTags)

	g.app.POST("/book",
		g.middleware.JwtAuthorization(),
		g.middleware.RbacAuthorization(
			map[uint]bool{
				1: true,
			},
		), httpHandler.CreateOneBook)
}
