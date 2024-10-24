package middlewareHandler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/middleware/middlewareUsecase"
)

type (
	MiddlewareHttpHandler interface {
		JwtAuthorization() gin.HandlerFunc
		RbacAuthorization(expectedRoleID map[uint]bool) gin.HandlerFunc
		BookOwnershipAuthorization() gin.HandlerFunc
	}

	middlewareHttpHandlerImpl struct {
		cfg               *config.Config
		middlewareUsecase middlewareUsecase.MiddlewareUsecase
	}
)

// OwnershipAuthorization implements MiddlewareHttpHandler.
func (m *middlewareHttpHandlerImpl) BookOwnershipAuthorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.GetUint("userID")
		bookIDStr := ctx.Param("bookID")
		roleID := ctx.GetUint("roleID")

		if bookIDStr == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "bookID is required"})
			return
		}
		bookIDUint64, err := strconv.ParseUint(bookIDStr, 10, 64)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		bookID := uint(bookIDUint64)
		if err := m.middlewareUsecase.BookOwnershipAuthorization(m.cfg, roleID, userID, bookID); err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}
		ctx.Set("bookID", bookID)
		ctx.Next()
	}
}

// RbacAuthorization implements MiddlewareHttpHandler.
func (m *middlewareHttpHandlerImpl) RbacAuthorization(expectedRoleID map[uint]bool) gin.HandlerFunc {
	fmt.Println("Call RbacAuthorization")
	return func(ctx *gin.Context) {
		roleID := ctx.GetUint("roleID")
		if err := m.middlewareUsecase.RbacAuthorization(roleID, expectedRoleID); err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return
		}
		ctx.Next()
	}
}

// JwtAuthorization implements MiddlewareHttpHandler.
func (m *middlewareHttpHandlerImpl) JwtAuthorization() gin.HandlerFunc {
	fmt.Println("Call JwtAuthorization")

	return func(ctx *gin.Context) {
		accessToken := strings.TrimPrefix(ctx.GetHeader("Authorization"), "Bearer ")
		if accessToken == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized access",
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
