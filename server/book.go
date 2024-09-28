package server

import (
	"image/jpeg"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kritAsawaniramol/book-store/module/book/bookHandler"
	"github.com/kritAsawaniramol/book-store/module/book/bookRepository"
	"github.com/kritAsawaniramol/book-store/module/book/bookUsecase"
)

func (g *ginServer) bookService() {
	repo := bookRepository.NewBookRepositoryImpl(g.db)
	usecase := bookUsecase.NewBookUsecaseImpl(repo)
	httpHandler := bookHandler.NewBookHttpHandlerImpl(g.cfg, usecase)
	grpcHandler := bookHandler.NewBookGrpcHandlerImpl(g.cfg, usecase)
	_ = grpcHandler

	// g.app.Static("/book/cover", "./asset/image/bookCover")
	g.app.StaticFile("/book/file", "./asset/book")
	g.app.GET("/book/cover/:fileName", func(ctx *gin.Context) {
		fileName := ctx.Param("fileName")
		file, err := os.Open("./asset/image/bookCover/" + fileName)
		if err != nil {
			log.Printf("error: %s\n", err.Error())
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		defer file.Close()
		ctx.Header("Content-type", "image")

		image, err := jpeg.Decode(file)
		if err != nil {
			log.Printf("error: %s\n", err.Error())
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		jpeg.Encode(ctx.Writer, image, nil)
	})
	g.app.POST("/book", g.middleware.JwtAuthorization(), g.middleware.RbacAuthorization(
		map[uint]bool{
			1: true,
		},
	), httpHandler.CreateOneBook)
}
