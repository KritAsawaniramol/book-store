package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"mime/multipart"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/nfnt/resize"
)

type (
	contextWrapperService interface {
		Bind(data any) error
		BindPostForm(data any) error
		SaveImageFormFile(name string, dst string) (string, error)
		SavePdfFormFile(name string, dst string) (string, error)
	}

	contextWrapper struct {
		Context   *gin.Context
		validator *validator.Validate
	}
)

func resizeImage(file *multipart.FileHeader) (image.Image, string, error) {
	f, err := file.Open()
	if err != nil {
		log.Printf("error: IscontentTypeImage: %s\n", err.Error())
		return nil, "", errors.New("error: cound not open file")
	}
	defer f.Close()
	image, format, err := image.Decode(f)
	if err != nil {
		log.Printf("error: IscontentTypeImage: %s\n", err.Error())
		return nil, "", errors.New("error: invalid image format")
	}
	resizedImage := resize.Resize(1000, 0, image, resize.Lanczos3)
	return resizedImage, format, nil
}

func (c *contextWrapper) SaveImageFormFile(name string, dst string) (string, error) {

	file, err := c.Context.FormFile(name)
	if err != nil {
		log.Printf("error: SaveImageFormFile: %s\n", err.Error())
		return "", errors.New("error: image not found")
	}
	if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
		return "", errors.New("error: invalid file type. only images are allowed")
	}

	resizedImage, format, err := resizeImage(file)
	if err != nil {
		return "", err
	}

	imageName := uuid.New().String()
	imagePath := fmt.Sprintf("%s/%s.%s", dst, imageName, format)

	if err := os.MkdirAll(dst, os.ModePerm); err != nil {
		log.Printf("error: SaveImageFormFile: %s\n", err.Error())
		return "", errors.New("error: save image failed")
	}
	out, err := os.Create(imagePath)
	if err != nil {
		log.Printf("error: SaveImageFormFile: %s\n", err.Error())
		return "", errors.New("error: save image failed")
	}
	defer out.Close()

	switch strings.ToLower(format) {
	case "jpeg", "jpg":
		err = jpeg.Encode(out, resizedImage, nil)
	case "png":
		err = png.Encode(out, resizedImage)
	default:
		err = errors.New("error: upsupported image format")
	}
	if err != nil {
		log.Printf("error: SaveImageFormFile: %s\n", err.Error())
		return "", err
	}

	return imagePath, nil
}

func (c *contextWrapper) SavePdfFormFile(name string, dst string) (string, error) {
	file, err := c.Context.FormFile(name)
	if err != nil {
		log.Printf("error: SavePdfFormFile: %s\n", err.Error())
		return "", errors.New("error: .pdf file not found")
	}

	if err := isPdf(file); err != nil {
		return "", err
	}

	bookFilePath := fmt.Sprintf("%s/%s.pdf", dst, uuid.New().String())

	if err := c.Context.SaveUploadedFile(file, bookFilePath); err != nil {
		log.Printf("error: SavePdfFormFile: %s\n", err.Error())
		return "", errors.New("error: save .pdf file failed")
	}

	return bookFilePath, nil
}

// Bind implements contextWrapperService.
func (c *contextWrapper) Bind(data any) error {
	defer c.Context.Request.Body.Close()
	if err := c.Context.Bind(data); err != nil {
		log.Printf("Error: Bind data failed: %s", err.Error())
		return ErrBadReq
	}

	if err := c.validator.Struct(data); err != nil {
		log.Printf("Error: Validate data failed: %s", err.Error())
		return ErrValidateDataFail
	}
	return nil
}

func (c *contextWrapper) BindPostForm(data any) error {
	jsonString := c.Context.PostForm("json")
	fmt.Printf("jsonString: %v\n", jsonString)

	err := json.Unmarshal([]byte(jsonString), data)
	if err != nil {
		log.Printf("Error: Bind data failed: %s", err.Error())
		return errors.New("errors: bad requset")
	}

	if err := c.validator.Struct(data); err != nil {
		log.Printf("Error: Validate data failed: %s", err.Error())
		return errors.New("errors: validate data failed")
	}
	return nil
}

func ContextWrapper(ctx *gin.Context) contextWrapperService {
	v := validator.New()
	return &contextWrapper{
		Context:   ctx,
		validator: v,
	}
}

func isPdf(file *multipart.FileHeader) error {
	if file.Header.Get("Content-Type") != "application/pdf" {
		return errors.New("error: file is not pdf")
	}
	f, err := file.Open()
	if err != nil {
		log.Printf("error: IsPdf: %s\n", err.Error())
		return errors.New("error: cound not open file")
	}
	defer f.Close()

	buf := make([]byte, 5)
	_, err = f.Read(buf)
	if err != nil {
		log.Printf("error: IsPdf: %s\n", err.Error())
		return errors.New("error: cound not read file")
	}

	if string(buf) != "%PDF-" {
		return errors.New("error: invalid pdf file")
	}

	return nil
}
