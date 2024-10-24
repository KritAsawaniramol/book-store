package middlewareUsecase

import (
	"errors"
	"log"

	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/middleware/middlewareRepository"
	"github.com/kritAsawaniramol/book-store/pkg/jwtAuth"
)

type middlewareUsecaseImpl struct {
	middlewareRepository middlewareRepository.MiddlewareRepository
}

// BookOwnershipAuthorization implements MiddlewareUsecase.
func (m *middlewareUsecaseImpl) BookOwnershipAuthorization(cfg *config.Config, roleID uint, userID uint, bookID uint) error {
	// admin bypass
	if roleID == 1 {
		log.Println("admin read book")
		return nil
	}

	return m.middlewareRepository.BookShelfSearch(cfg.Grpc.ShelfUrl, userID, bookID)
}

// RbacAuthorization implements MiddlewareUsecase.
func (m *middlewareUsecaseImpl) RbacAuthorization(roleID uint, expectedRoleID map[uint]bool) error {
	v, ok := expectedRoleID[roleID]
	if !ok || !v {
		return errors.New("error: permission denied")
	}
	return nil
}

// JwtAuthorization implements MiddlewareUsecase.
func (m *middlewareUsecaseImpl) JwtAuthorization(cfg *config.Config, accessToken string) (uint, uint, error) {
	claims, err := jwtAuth.ParseToken(cfg.Jwt.AccessSecretKey, accessToken)
	if err != nil {
		return 0, 0, err
	}

	if err := m.middlewareRepository.AccessTokenSearch(cfg.Grpc.AuthUrl, accessToken); err != nil {
		return 0, 0, err
	}

	return claims.UserID, claims.RoleID, err
}

func NewMiddlewareUsecaseImpl(middlewareRepository middlewareRepository.MiddlewareRepository) MiddlewareUsecase {
	return &middlewareUsecaseImpl{
		middlewareRepository: middlewareRepository,
	}
}
