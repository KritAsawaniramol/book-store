package bookHandler

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/book"
	"github.com/kritAsawaniramol/book-store/module/book/bookUsecase"
	"github.com/kritAsawaniramol/book-store/pkg/request"
)

type (
	BookHttpHandler interface {
		CreateOneBook(ctx *gin.Context)
		GetBookCover(ctx *gin.Context)
		SearchBooks(ctx *gin.Context)
		GetOneBook(ctx *gin.Context)
		GetTags(ctx *gin.Context)
	}

	bookHttpHandlerImpl struct {
		cfg         *config.Config
		bookUsecase bookUsecase.BookUsecase
	}
)

// GetTags implements BookHttpHandler.
func (b *bookHttpHandlerImpl) GetTags(ctx *gin.Context) {
	tags, err := b.bookUsecase.GetTags()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, tags)
}

// GetOneBook implements BookHttpHandler.
func (b *bookHttpHandlerImpl) GetOneBook(ctx *gin.Context) {
	bookIDStr := ctx.Param("id")
	bookID, err := strconv.ParseUint(bookIDStr, 10, 64)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	res, err := b.bookUsecase.GetOneBook(uint(bookID))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

// SearchBooks implements BookHttpHandler.
func (b *bookHttpHandlerImpl) SearchBooks(ctx *gin.Context) {
	wrapper := request.ContextWrapper(ctx)
	req := &book.SearchBooksReq{}
	if err := wrapper.Bind(req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	res, err := b.bookUsecase.SearchBooks(req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

// GetBookCover implements BookHttpHandler.
func (b *bookHttpHandlerImpl) GetBookCover(ctx *gin.Context) {
	bookIDStr := ctx.Param("id")
	bookID, err := strconv.ParseUint(bookIDStr, 10, 64)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	res, err := b.bookUsecase.GetOneBook(uint(bookID))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

// CreateOneBook implements BookHttpHandler.
func (b *bookHttpHandlerImpl) CreateOneBook(ctx *gin.Context) {
	wrapper := request.ContextWrapper(ctx)
	req := &book.CreateBookReq{}
	if err := wrapper.BindPostForm(req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	imagePath, err := wrapper.SaveImageFormFile("book_cover", "asset/image/bookCover")
	if err != nil {
		if err.Error() == "error: image not found" {
			imagePath = "asset/image/bookCover/default/book-store_default_bookCover.png"
		} else {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
	}

	bookFilePath, err := wrapper.SavePdfFormFile("book_file", "asset/book")
	if err != nil {
		rollBackSaveFile(imagePath)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	req.CoverImagePath = imagePath
	req.FilePath = bookFilePath

	bookID, err := b.bookUsecase.CreateOneBook(req)
	if err != nil {
		rollBackSaveFile(imagePath, bookFilePath)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"book_id": bookID})
}

func NewBookHttpHandlerImpl(cfg *config.Config, bookUsecase bookUsecase.BookUsecase) BookHttpHandler {
	return &bookHttpHandlerImpl{cfg: cfg, bookUsecase: bookUsecase}
}

func rollBackSaveFile(filePaths ...string) error {
	for _, v := range filePaths {
		err := os.Remove(v)
		if err != nil {
			log.Printf("error: rollBackRemoveFile: %s\n", err.Error())
			return err
		}
	}
	return nil
}
