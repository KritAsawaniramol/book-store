package middlewareUsecase

import (
	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/middleware/middlewareRepository"
	"github.com/kritAsawaniramol/book-store/pkg/jwtAuth"
)

type middlewareUsecaseImpl struct {
	middlewareRepository middlewareRepository.MiddlewareRepository
}

// JwtAuthorization implements MiddlewareUsecase.
func (m *middlewareUsecaseImpl) JwtAuthorization(cfg *config.Config, accessToken string) (uint, uint, error) {
	claims, err := jwtAuth.ParseToken(cfg.Jwt.AccessSecretKey, accessToken)
	if err != nil {
		return 0, 0, err
	}

	if  err := m.middlewareRepository.AccessTokenSearch(cfg.Grpc.AuthUrl,accessToken); err != nil {
		return 0, 0, err
	}

	return claims.UserID, claims.RoleID, err
}

func NewMiddlewareUsecaseImpl(middlewareRepository middlewareRepository.MiddlewareRepository) MiddlewareUsecase {
	return &middlewareUsecaseImpl{
		middlewareRepository: middlewareRepository,
	}
}
