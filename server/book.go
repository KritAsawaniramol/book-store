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

	book := g.app.Group("/book_v1")

	book.GET("", g.healthCheck)
	book.StaticFile("/book/file", "./asset/book")
	book.GET("/book/cover/:fileName", httpHandler.GetBookCoverImage)
	book.GET("/book", httpHandler.SearchBooks)
	book.GET("/admin/book",
		g.middleware.JwtAuthorization(),
		g.middleware.RbacAuthorization(
			map[uint]bool{
				1: true,
			},
		), httpHandler.SearchBooks)

	book.GET("/book/:id", httpHandler.GetOneBook)
	book.GET("/admin/book/:id",
		g.middleware.JwtAuthorization(),
		g.middleware.RbacAuthorization(
		map[uint]bool{
				1: true,
			},
		), httpHandler.GetOneBook)

	book.GET("/book/tags", httpHandler.GetTags)
	book.GET("/book/read/:bookID",
		g.middleware.JwtAuthorization(),
		g.middleware.BookOwnershipAuthorization(),
		httpHandler.ReadBook)

	book.POST("/book",
		g.middleware.JwtAuthorization(),
		g.middleware.RbacAuthorization(
			map[uint]bool{
				1: true,
			},
		), httpHandler.CreateOneBook)

	book.PATCH("/book/:id",
		g.middleware.JwtAuthorization(),
		g.middleware.RbacAuthorization(
			map[uint]bool{
				1: true,
			},
		), httpHandler.UpdateOneBook)

	book.PATCH("/book/cover/:id",
		g.middleware.JwtAuthorization(),
		g.middleware.RbacAuthorization(
			map[uint]bool{
				1: true,
			},
		), httpHandler.UpdateOneBookCover)

	book.PATCH("/book/file/:id",
		g.middleware.JwtAuthorization(),
		g.middleware.RbacAuthorization(
			map[uint]bool{
				1: true,
			},
		), httpHandler.UpdateOneBookFile)

}
