package authUsecase

import (
	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/auth"
	"github.com/kritAsawaniramol/book-store/module/auth/authRepository"
	"github.com/kritAsawaniramol/book-store/module/user/userPb"
	"github.com/kritAsawaniramol/book-store/pkg/jwtAuth"
)

type authUsecaseImpl struct {
	authRepository authRepository.AuthRepository
}

// RefreshToken implements AuthUsecase.
func (a *authUsecaseImpl) RefreshToken(cfg *config.Config, req *auth.RefreshTokenReq) (*auth.CredentialRes, error) {
	claims, err := jwtAuth.ParseToken(cfg.Jwt.RefreshSecretKey, req.RefreshToken)
	if err != nil {
		return nil, err
	}

	userProfile, err := a.authRepository.FindOneUserProfileToRefresh(cfg.Grpc.UserUrl, &userPb.FindOneUserProfileToRefreshReq{
		UserId: uint64(claims.UserID),
	})
	if err != nil {
		return nil, err
	}

	accessToken := jwtAuth.NewAccessToken(cfg.Jwt.AccessSecretKey, cfg.Jwt.AccessDuration, &jwtAuth.Claims{
		UserID: uint(userProfile.Id),
		RoleID: uint(userProfile.RoleId),
	}).SignToken()

	refreshToken := jwtAuth.ReloadToken(cfg.Jwt.RefreshSecretKey, claims.ExpiresAt.Unix(), &jwtAuth.Claims{
		UserID: uint(userProfile.Id),
		RoleID: uint(userProfile.RoleId),
	})

	if err := a.authRepository.UpdateOneCredentialByID(req.CredentialID, &auth.Credential{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}); err != nil {
		return nil, err
	}

	condition := &auth.Credential{}
	condition.ID = req.CredentialID
	credential, err := a.authRepository.GetOneUserCredential(condition)
	if err != nil {
		return nil, err
	}

	return &auth.CredentialRes{
		ID:           credential.ID,
		UserID:       credential.UserID,
		AccessToken:  credential.AccessToken,
		RefreshToken: credential.RefreshToken,
		CreatedAt:    credential.CreatedAt,
		UpdatedAt:    credential.UpdatedAt,
	}, nil
}

// Logout implements AuthUsecase.
func (a *authUsecaseImpl) Logout(req *auth.LogoutReq) error {
	return a.authRepository.DeleteOneUserCredentialByID(req.CredentialID)
}

// Login implements AuthUsecase.
func (a *authUsecaseImpl) Login(cfg *config.Config, req *auth.LoginReq) (*auth.LoginRes, error) {
	userProfile, err := a.authRepository.FindOneUserProfileToLogin(
		cfg.Grpc.UserUrl,
		&userPb.FindUserProfileToLoginReq{
			Username: req.Username,
			Password: req.Password,
		})
	if err != nil {
		return nil, err
	}

	credentialID, err := a.authRepository.CreateOneUserCredential(&auth.Credential{
		UserID: uint(userProfile.Id),
		RoleID: uint(userProfile.RoleId),
		AccessToken: jwtAuth.NewAccessToken(cfg.Jwt.AccessSecretKey, cfg.Jwt.AccessDuration, &jwtAuth.Claims{
			UserID: uint(userProfile.Id),
			RoleID: uint(userProfile.RoleId),
		}).SignToken(),
		RefreshToken: jwtAuth.NewRefreshToken(cfg.Jwt.RefreshSecretKey, cfg.Jwt.RefreshDuration, &jwtAuth.Claims{
			UserID: uint(userProfile.Id),
			RoleID: uint(userProfile.RoleId),
		}).SignToken(),
	})
	if err != nil {
		return nil, err
	}

	condition := &auth.Credential{}
	condition.ID = credentialID
	credential, err := a.authRepository.GetOneUserCredential(condition)
	if err != nil {
		return nil, err
	}

	return &auth.LoginRes{
		ID:        uint(userProfile.Id),
		Username:  userProfile.Username,
		RoleID:    uint(userProfile.RoleId),
		Coin:      userProfile.Coin,
		CreatedAt: userProfile.CreatedAt.AsTime().Local(),
		UpdatedAt: userProfile.UpdatedAt.AsTime().Local(),
		Credential: &auth.CredentialRes{
			ID:           credential.ID,
			UserID:       credential.UserID,
			AccessToken:  credential.AccessToken,
			RefreshToken: credential.RefreshToken,
			CreatedAt:    credential.CreatedAt,
			UpdatedAt:    credential.UpdatedAt,
		},
	}, nil
}

func NewAuthUsecaseImpl(authRepository authRepository.AuthRepository) AuthUsecase {
	return &authUsecaseImpl{authRepository: authRepository}
}
