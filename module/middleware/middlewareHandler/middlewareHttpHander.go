package middlewareHandler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/middleware/middlewareUsecase"
)

type (
	MiddlewareHttpHandler interface {
		JwtAuthorization() gin.HandlerFunc
	}

	middlewareHttpHandlerImpl struct {
		cfg               *config.Config
		middlewareUsecase middlewareUsecase.MiddlewareUsecase
	}
)

// JwtAuthorization implements MiddlewareHttpHandler.
func (m *middlewareHttpHandlerImpl) JwtAuthorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken := strings.TrimPrefix(ctx.GetHeader("Authorization"), "Bearer ")
		if accessToken == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized access",
			})
			return
		}

		userID, roleID, err := m.middlewareUsecase.JwtAuthorization(m.cfg, accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return
		}

		ctx.Set("userID", userID)
		ctx.Set("roleID", roleID)
		ctx.Next()
	}
}

func NewMiddlewareHttpHandlerImpl(
	cfg *config.Config,
	middlewareUsecase middlewareUsecase.MiddlewareUsecase,
) MiddlewareHttpHandler {
	return &middlewareHttpHandlerImpl{
		cfg:               cfg,
		middlewareUsecase: middlewareUsecase,
	}
}