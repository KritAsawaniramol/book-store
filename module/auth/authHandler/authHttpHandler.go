package authHandler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/auth"
	"github.com/kritAsawaniramol/book-store/module/auth/authUsecase"
	"github.com/kritAsawaniramol/book-store/pkg/request"
)

type (
	AuthHttpHandler interface {
		Login(*gin.Context)
	}

	authHttpHandlerImpl struct {
		cfg         *config.Config
		authUsecase authUsecase.AuthUsecase
	}
)

// Login implements AuthHttpHandler.
func (a *authHttpHandlerImpl) Login(ctx *gin.Context) {
	wrapper := request.ContextWrapper(ctx)
	req := &auth.LoginReq{}
	err := wrapper.Bind(req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}


	res, err := a.authUsecase.Login(a.cfg, req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func NewAuthHttpHandlerImpl(cfg *config.Config, authUsecase authUsecase.AuthUsecase) AuthHttpHandler {
	return &authHttpHandlerImpl{cfg: cfg, authUsecase: authUsecase}
}
