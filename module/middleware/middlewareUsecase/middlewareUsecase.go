package middlewareUsecase

import "github.com/kritAsawaniramol/book-store/config"

type MiddlewareUsecase interface {
	JwtAuthorization(cfg *config.Config, accessToken string) (uint, uint, error)
	RbacAuthorization(roleID uint, expectedRoleID map[uint]bool) error
}
