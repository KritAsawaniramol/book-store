package request

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type (
	contextWrapperService interface {
		Bind(data any) error
	}

	contextWrapper struct {
		Context   *gin.Context
		validator *validator.Validate
	}
)

// Bind implements contextWrapperService.
func (c *contextWrapper) Bind(data any) error {
	if err := c.Context.Bind(data); err != nil {
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
