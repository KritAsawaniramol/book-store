package authUsecase_test

import (
	"errors"
	"testing"
	"time"

	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/auth"
	"github.com/kritAsawaniramol/book-store/module/auth/authPb"
	"github.com/kritAsawaniramol/book-store/module/auth/authRepository"
	"github.com/kritAsawaniramol/book-store/module/auth/authUsecase"
	"github.com/kritAsawaniramol/book-store/module/user/userPb"
	"github.com/kritAsawaniramol/book-store/pkg/jwtAuth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type (
	testLogin struct {
		name     string
		cfg      *config.Config
		req      *auth.LoginReq
		expected *auth.LoginRes
		isErr    bool
	}
)

func TestLogout(t *testing.T) {
	repoMock := authRepository.NewAuthRepositoryMock()
	usecase := authUsecase.NewAuthUsecaseImpl(repoMock)
	t.Run("success", func(t *testing.T) {
		repoMock.On("DeleteOneUserCredentialByID", uint(1)).Return(nil)
		err := usecase.Logout(&auth.LogoutReq{CredentialID: 1})
		assert.NoError(t, err)
	})

	t.Run("DeleteOneUserCredentialByID fail", func(t *testing.T) {
		repoMock.ExpectedCalls = nil
		repoMock.On("DeleteOneUserCredentialByID", uint(1)).Return(auth.ErrDeleteUserCredentialFail)
		err := usecase.Logout(&auth.LogoutReq{CredentialID: 1})
		assert.ErrorIs(t, auth.ErrDeleteUserCredentialFail, err)
	})
}

func TestRefreshToken(t *testing.T) {
	cfg := config.LoadConfig("../../../env/test/.env")
	mockUserProfile := &userPb.UserProfile{
		Id:        1,
		Username:  "user1",
		RoleId:    1,
		Coin:      100,
		CreatedAt: timestamppb.New(time.Time{}),
		UpdatedAt: timestamppb.New(time.Time{}),
	}
	mockCredential := &auth.Credential{
		UserID:       1,
		AccessToken:  "access_token",
		RefreshToken: "refresh_token",
	}
	mockCredential.ID = 1
	mockCredential.CreatedAt = time.Time{}
	mockCredential.UpdatedAt = time.Time{}

	mockSuccessCredentialRes := &auth.CredentialRes{
		ID:           mockCredential.ID,
		UserID:       mockCredential.UserID,
		AccessToken:  mockCredential.AccessToken,
		RefreshToken: mockCredential.RefreshToken,
		CreatedAt:    mockCredential.CreatedAt,
		UpdatedAt:    mockCredential.UpdatedAt,
	}

	repoMock := authRepository.NewAuthRepositoryMock()
	mockRefreshToken := jwtAuth.NewRefreshToken(cfg.Jwt.RefreshSecretKey, cfg.Jwt.RefreshDuration, &jwtAuth.Claims{}).SignToken()

	usecase := authUsecase.NewAuthUsecaseImpl(repoMock)
	t.Run("success", func(t *testing.T) {
		repoMock.On("FindOneUserProfileToRefresh", cfg, mock.AnythingOfType("*userPb.FindOneUserProfileToRefreshReq")).
			Return(mockUserProfile, nil)

		repoMock.On("UpdateOneCredentialByID", uint(1), mock.AnythingOfType("*auth.Credential")).
			Return(nil)

		repoMock.On("GetOneUserCredential", mock.AnythingOfType("*auth.Credential")).
			Return(mockCredential, nil)

		res, err := usecase.RefreshToken(cfg, &auth.RefreshTokenReq{
			CredentialID: 1,
			RefreshToken: mockRefreshToken,
		})

		assert.NoError(t, err)
		assert.Equal(t, mockSuccessCredentialRes, res)
	})

	t.Run("ParseToken fail", func(t *testing.T) {
		repoMock.ExpectedCalls = nil

		res, err := usecase.RefreshToken(cfg, &auth.RefreshTokenReq{
			CredentialID: 1,
			RefreshToken: "",
		})

		repoMock.AssertNotCalled(t, "FindOneUserProfileToRefresh")
		repoMock.AssertNotCalled(t, "UpdateOneCredentialByID")
		repoMock.AssertNotCalled(t, "GetOneUserCredential")

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("FindOneUserProfileToRefresh fail", func(t *testing.T) {
		repoMock.ExpectedCalls = nil

		repoMock.On("FindOneUserProfileToRefresh", cfg, mock.AnythingOfType("*userPb.FindOneUserProfileToRefreshReq")).
			Return(nil, auth.ErrUserNotFound)

		repoMock.AssertNotCalled(t, "UpdateOneCredentialByID")
		repoMock.AssertNotCalled(t, "GetOneUserCredential")

		res, err := usecase.RefreshToken(cfg, &auth.RefreshTokenReq{
			CredentialID: 1,
			RefreshToken: mockRefreshToken,
		})

		assert.ErrorIs(t, auth.ErrUserNotFound, err)
		assert.Nil(t, res)
	})

	t.Run("UpdateOneCredentialByID fail", func(t *testing.T) {
		repoMock.ExpectedCalls = nil

		repoMock.On("FindOneUserProfileToRefresh", cfg, mock.AnythingOfType("*userPb.FindOneUserProfileToRefreshReq")).
			Return(mockUserProfile, nil)

		repoMock.On("UpdateOneCredentialByID", uint(1), mock.AnythingOfType("*auth.Credential")).
			Return(auth.ErrUpdateCredential)

		repoMock.AssertNotCalled(t, "GetOneUserCredential")

		res, err := usecase.RefreshToken(cfg, &auth.RefreshTokenReq{
			CredentialID: 1,
			RefreshToken: mockRefreshToken,
		})

		assert.ErrorIs(t, auth.ErrUpdateCredential, err)
		assert.Nil(t, res)
	})

	t.Run("GetOneUserCredential fail", func(t *testing.T) {
		repoMock.ExpectedCalls = nil

		repoMock.On("FindOneUserProfileToRefresh", cfg, mock.AnythingOfType("*userPb.FindOneUserProfileToRefreshReq")).
			Return(mockUserProfile, nil)

		repoMock.On("UpdateOneCredentialByID", uint(1), mock.AnythingOfType("*auth.Credential")).
			Return(nil)

		repoMock.On("GetOneUserCredential", mock.AnythingOfType("*auth.Credential")).
			Return(nil, auth.ErrCredentialNotFound)

		res, err := usecase.RefreshToken(cfg, &auth.RefreshTokenReq{
			CredentialID: 1,
			RefreshToken: mockRefreshToken,
		})

		assert.ErrorIs(t, auth.ErrCredentialNotFound, err)
		assert.Nil(t, res)
	})

}

func TestAccessTokenSearch(t *testing.T) {
	mockCredential := &auth.Credential{
		UserID:       1,
		AccessToken:  "access_token",
		RefreshToken: "refresh_token",
	}
	mockCredential.ID = 1
	mockCredential.CreatedAt = time.Time{}
	mockCredential.UpdatedAt = time.Time{}
	repoMock := authRepository.NewAuthRepositoryMock()
	usecase := authUsecase.NewAuthUsecaseImpl(repoMock)

	t.Run("success", func(t *testing.T) {
		repoMock.On("GetOneUserCredential", mock.AnythingOfType("*auth.Credential")).Return(mockCredential, nil)
		res, err := usecase.AccessTokenSearch("access_token")
		assert.NoError(t, err)
		assert.Equal(t, &authPb.AccessTokenSearchRes{IsValid: true}, res)
	})

	t.Run("credential is nil", func(t *testing.T) {
		repoMock.ExpectedCalls = nil
		repoMock.On("GetOneUserCredential", mock.AnythingOfType("*auth.Credential")).Return(nil, nil)
		res, err := usecase.AccessTokenSearch("access_token")
		assert.ErrorIs(t, auth.ErrCredentialNotFound, err)
		assert.Equal(t, &authPb.AccessTokenSearchRes{IsValid: false}, res)
	})
	t.Run("GetOneUserCredential error", func(t *testing.T) {
		repoMock.ExpectedCalls = nil
		repoMock.On("GetOneUserCredential", mock.AnythingOfType("*auth.Credential")).Return(nil, auth.ErrCredentialNotFound)
		_, err := usecase.AccessTokenSearch("access_token")
		assert.ErrorIs(t, auth.ErrCredentialNotFound, err)
	})

}

func TestLogin(t *testing.T) {
	repoMock := authRepository.NewAuthRepositoryMock()
	usecase := authUsecase.NewAuthUsecaseImpl(repoMock)
	cfg := config.LoadConfig("../../../env/test/.env")

	testCases := []testLogin{
		{
			name: "success",
			cfg:  cfg,
			req: &auth.LoginReq{
				Username: "user1",
				Password: "password1",
			},
			expected: &auth.LoginRes{
				ID:        1,
				Username:  "user1",
				RoleID:    1,
				Coin:      100,
				CreatedAt: timestamppb.New(time.Time{}).AsTime().Local(),
				UpdatedAt: timestamppb.New(time.Time{}).AsTime().Local(),
				Credential: &auth.CredentialRes{
					ID:           1,
					UserID:       1,
					AccessToken:  "access_token",
					RefreshToken: "refresh_token",
					CreatedAt:    time.Time{},
					UpdatedAt:    time.Time{},
				},
			},
			isErr: false,
		},
	}

	mockUserProfile := &userPb.UserProfile{
		Id:        1,
		Username:  "user1",
		RoleId:    1,
		Coin:      100,
		CreatedAt: timestamppb.New(time.Time{}),
		UpdatedAt: timestamppb.New(time.Time{}),
	}

	mockCredential := &auth.Credential{
		UserID:       1,
		AccessToken:  "access_token",
		RefreshToken: "refresh_token",
	}
	mockCredential.ID = 1
	mockCredential.CreatedAt = time.Time{}
	mockCredential.UpdatedAt = time.Time{}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			repoMock.On(
				"FindOneUserProfileToLogin",
				cfg,
				mock.AnythingOfType("*userPb.FindUserProfileToLoginReq"),
			).Return(mockUserProfile, nil)

			repoMock.On(
				"CreateOneUserCredential",
				mock.AnythingOfType("*auth.Credential"),
			).Return(uint(1), nil)

			repoMock.On(
				"GetOneUserCredential",
				mock.AnythingOfType("*auth.Credential"),
			).Return(
				mockCredential,
				nil,
			)

			result, err := usecase.Login(cfg, test.req)

			if test.isErr {
				assert.NotEmpty(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, result)
			}
		})
	}

	t.Run("FindOneUserProfileToLogin fail", func(t *testing.T) {
		// remove cache in mock
		repoMock.ExpectedCalls = nil
		repoMock.On(
			"FindOneUserProfileToLogin",
			cfg,
			mock.AnythingOfType("*userPb.FindUserProfileToLoginReq"),
		).Return(nil, errors.New("error: email or password are incorrect"))

		_, err := usecase.Login(cfg, &auth.LoginReq{
			Username: "user1",
			Password: "password",
		})

		repoMock.AssertNotCalled(t, "CreateOneUserCredential")
		repoMock.AssertNotCalled(t, "GetOneUserCredential")
		assert.NotEmpty(t, err)
	})
	t.Run("CreateOneUserCredential fail", func(t *testing.T) {
		// remove cache in mock
		repoMock.ExpectedCalls = nil

		repoMock.On(
			"FindOneUserProfileToLogin",
			cfg,
			mock.AnythingOfType("*userPb.FindUserProfileToLoginReq"),
		).Return(mockUserProfile, nil)

		repoMock.On(
			"CreateOneUserCredential",
			mock.AnythingOfType("*auth.Credential"),
		).Return(uint(0), errors.New("error: create user credential failed"))

		_, err := usecase.Login(cfg, &auth.LoginReq{
			Username: "user1",
			Password: "password",
		})

		repoMock.AssertNotCalled(t, "GetOneUserCredential")
		assert.NotEmpty(t, err)
	})

	t.Run("GetOneUserCredential fail", func(t *testing.T) {
		// remove cache in mock
		repoMock.ExpectedCalls = nil

		repoMock.On(
			"FindOneUserProfileToLogin",
			cfg,
			mock.AnythingOfType("*userPb.FindUserProfileToLoginReq"),
		).Return(mockUserProfile, nil)

		repoMock.On(
			"CreateOneUserCredential",
			mock.AnythingOfType("*auth.Credential"),
		).Return(uint(1), nil)

		repoMock.On(
			"GetOneUserCredential",
			mock.AnythingOfType("*auth.Credential"),
		).Return(nil, errors.New("error: credential not found"))
		_, err := usecase.Login(cfg, &auth.LoginReq{
			Username: "user1",
			Password: "password",
		})
		assert.NotEmpty(t, err)
	})
}
