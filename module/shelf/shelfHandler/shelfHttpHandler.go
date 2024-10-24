package shelfHandler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/shelf/shelfUsecase"
)

type (
	ShelfHttpHandler interface {
		GetMyShelf(ctx *gin.Context)
	}

	shelfHttpHandlerImpl struct {
		cfg          *config.Config
		shelfUsecase shelfUsecase.ShelfUsecase
	}
)

// GetMyShelf implements ShelfHttpHandler.
func (s *shelfHttpHandlerImpl) GetMyShelf(ctx *gin.Context) {
	bookIDStr := ctx.Query("book_id")
	var bookID uint = 0
	if bookIDStr != "" {
		bookIDUint64, err := strconv.ParseUint(bookIDStr, 10, 64)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		bookID = uint(bookIDUint64)
	}
	res, err := s.shelfUsecase.GetMyShelves(s.cfg, ctx.GetUint("userID"), bookID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func NewShelfHttpHandlerImpl(cfg *config.Config, shelfUescase shelfUsecase.ShelfUsecase) ShelfHttpHandler {
	return &shelfHttpHandlerImpl{
		cfg:          cfg,
		shelfUsecase: shelfUescase,
	}
}
