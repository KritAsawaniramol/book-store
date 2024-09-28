package middlewareUsecase

import (
	"errors"

	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/middleware/middlewareRepository"
	"github.com/kritAsawaniramol/book-store/pkg/jwtAuth"
)

type middlewareUsecaseImpl struct {
	middlewareRepository middlewareRepository.MiddlewareRepository
}

// RbacAuthorization implements MiddlewareUsecase.
func (m *middlewareUsecaseImpl) RbacAuthorization(roleID uint, expectedRoleID map[uint]bool)  error {
	v, ok := expectedRoleID[roleID]
	if !ok || v == false {
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
