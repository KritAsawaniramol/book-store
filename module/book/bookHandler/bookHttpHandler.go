package bookHandler

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/book"
	"github.com/kritAsawaniramol/book-store/module/book/bookUsecase"
	"github.com/kritAsawaniramol/book-store/pkg/request"
)

type (
	BookHttpHandler interface {
		CreateOneBook(ctx *gin.Context)
		SearchBooks(ctx *gin.Context)
		GetOneBook(ctx *gin.Context)
		GetTags(ctx *gin.Context)
		GetBookCoverImage(ctx *gin.Context)
		ReadBook(ctx *gin.Context)
		UpdateOneBook(ctx *gin.Context)
		UpdateOneBookCover(ctx *gin.Context)
		UpdateOneBookFile(ctx *gin.Context)
	}

	bookHttpHandlerImpl struct {
		cfg         *config.Config
		bookUsecase bookUsecase.BookUsecase
	}
)

// UpdateOneBookFile implements BookHttpHandler.
func (b *bookHttpHandlerImpl) UpdateOneBookFile(ctx *gin.Context) {
	wrapper := request.ContextWrapper(ctx)
	bookIDStr := ctx.Param("id")
	bookID, err := strconv.ParseUint(bookIDStr, 10, 64)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	bookFilePath, err := wrapper.SavePdfFormFile("book_file", "asset/book")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := b.bookUsecase.UpdateOneBookFile(uint(bookID), bookFilePath); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

// UpdateOneBookCover implements BookHttpHandler.
func (b *bookHttpHandlerImpl) UpdateOneBookCover(ctx *gin.Context) {
	wrapper := request.ContextWrapper(ctx)
	bookIDStr := ctx.Param("id")
	bookID, err := strconv.ParseUint(bookIDStr, 10, 64)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	imagePath, err := wrapper.SaveImageFormFile("book_cover", "asset/image/bookCover")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := b.bookUsecase.UpdateOneBookCover(uint(bookID), imagePath); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

// ReadBook implements BookHttpHandler.
func (b *bookHttpHandlerImpl) ReadBook(ctx *gin.Context) {
	bookID := ctx.GetUint("bookID")

	filePath, err := b.bookUsecase.GetOneBookFilePath(bookID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("error: ReadBook: %s\n", err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	defer file.Close()
	ctx.Writer.Header().Set("Content-type", "application/pdf")
	if _, err := io.Copy(ctx.Writer, file); err != nil {
		log.Printf("error: ReadBook: %s\n", err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.Status(http.StatusOK)
}

// GetBookCoverImage implements BookHttpHandler.
func (b *bookHttpHandlerImpl) GetBookCoverImage(ctx *gin.Context) {
	fileName := ctx.Param("fileName")
	if fileName == "book-store_default_bookCover.png" {
		fileName = "default/" + fileName
	}
	file, err := os.Open("./asset/image/bookCover/" + fileName)
	log.Println(fileName)
	if err != nil {
		log.Printf("error: %s\n", err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		log.Printf("error: %s\n", err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	// i, err := jpeg.Decode(file)
	// if err != nil {
	// 	log.Printf("error: GetBookCoverImage: %s\n", err.Error())
	// 	ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	// 	return
	// }

	// jpeg.Encode(ctx.Writer, img, nil)
	if err := serveImage(img, format, ctx); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
}

func serveImage(img image.Image, format string, ctx *gin.Context) error {
	format = strings.ToLower(format)
	ctx.Header("Content-type", fmt.Sprintf("image/%s", format))
	var err error = nil
	switch format {
	case "jpeg", "jpg":
		err = jpeg.Encode(ctx.Writer, img, nil)
	case "png":
		err = png.Encode(ctx.Writer, img)
	default:
		err = errors.New("error: upsupported image format")
	}
	if err != nil {
		log.Printf("error: serveImage: %s\n", err.Error())
		return err
	}
	return nil
}

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
	res, err := b.bookUsecase.GetOneBook(uint(bookID), ctx.GetUint("roleID"))
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
	res, err := b.bookUsecase.SearchBooks(req, ctx.GetUint("roleID"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
}


// UpdateOneBook implements BookHttpHandler.
func (b *bookHttpHandlerImpl) UpdateOneBook(ctx *gin.Context) {
	wrapper := request.ContextWrapper(ctx)
	req := &book.UpdateBookDetailReq{}

	bookIDStr := ctx.Param("id")
	bookID, err := strconv.ParseUint(bookIDStr, 10, 64)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	req.BookID = uint(bookID)

	if err := wrapper.Bind(req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := b.bookUsecase.UpdateOneBookDetail(req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}

	ctx.Status(http.StatusOK)
}

// CreateOneBook implements BookHttpHandler.
func (b *bookHttpHandlerImpl) CreateOneBook(ctx *gin.Context) {
	wrapper := request.ContextWrapper(ctx)
	req := &book.CreateBookReq{}

	if err := wrapper.BindPostForm(req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	imageNotFound := false
	imagePath, err := wrapper.SaveImageFormFile("book_cover", "asset/image/bookCover")
	if err != nil {
		if err.Error() == "error: image not found" {
			imageNotFound = true
			imagePath = "asset/image/bookCover/default/book-store_default_bookCover.png"
		} else {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
	}

	bookFilePath, err := wrapper.SavePdfFormFile("book_file", "asset/book")
	if err != nil {
		if !imageNotFound {
			rollBackSaveFile(imagePath)
		}
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
