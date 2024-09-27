package authHandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/auth"
	"github.com/kritAsawaniramol/book-store/module/auth/authUsecase"
	"github.com/kritAsawaniramol/book-store/pkg/request"
)

type (
	AuthHttpHandler interface {
		Login(ctx *gin.Context)
		Logout(ctx *gin.Context)
	}

	authHttpHandlerImpl struct {
		cfg         *config.Config
		authUsecase authUsecase.AuthUsecase
	}
)

// Logout implements AuthHttpHandler.
func (a *authHttpHandlerImpl) Logout(ctx *gin.Context) {
	wrapper := request.ContextWrapper(ctx)
	req := &auth.LogoutReq{}
	err := wrapper.Bind(req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := a.authUsecase.Logout(req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.Status(http.StatusOK)
}

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
