package authHandler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/auth"
	"github.com/kritAsawaniramol/book-store/module/auth/authHandler"
	"github.com/kritAsawaniramol/book-store/module/auth/authUsecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setup() (*gin.Context, *httptest.ResponseRecorder, *gin.Engine, *authUsecase.AuthUsecaseMock) {
	cfg := config.LoadConfig("../../../env/test/.env")
	usecaseMock := authUsecase.NewAuthUsecaseMock()
	httpHandler := authHandler.NewAuthHttpHandlerImpl(cfg, usecaseMock)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	ctx, router := gin.CreateTestContext(w)

	a := router.Group("/auth_v1")
	a.POST("/auth/login", httpHandler.Login)
	a.POST("/auth/logout", httpHandler.Logout)
	a.POST("/auth/refresh-token", httpHandler.RefreshToken)
	return ctx, w, router, usecaseMock
}

func TestLogin(t *testing.T) {
	ctx, w, router, usecaseMock := setup()
	loginReq := auth.LoginReq{
		Username: "user1",
		Password: "password",
	}
	byteLoginReq, _ := json.Marshal(loginReq)

	t.Run("success", func(t *testing.T) {

		loginRes := auth.LoginRes{
			ID:        1,
			Username:  "user1",
			RoleID:    1,
			Coin:      20,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			Credential: &auth.CredentialRes{
				ID:           1,
				UserID:       1,
				AccessToken:  "access_token",
				RefreshToken: "refresh_token",
				CreatedAt:    time.Time{},
				UpdatedAt:    time.Time{},
			},
		}

		usecaseMock.On("Login",
			mock.AnythingOfType("*config.Config"),
			mock.AnythingOfType("*auth.LoginReq")).Return(&loginRes, nil)
		ctx.Request = httptest.NewRequest(http.MethodPost, "/auth_v1/auth/login", bytes.NewBuffer(byteLoginReq))
		ctx.Request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, ctx.Request)
		assert.Equal(t, 200, w.Code)
		loginResJson, _ := json.Marshal(loginRes)
		assert.Equal(t, string(loginResJson), string(w.Body.String()))
	})

	t.Run("bind request fail", func(t *testing.T) {
		usecaseMock.ExpectedCalls = nil
		w := httptest.NewRecorder()
		ctx.Request = httptest.NewRequest(http.MethodPost, "/auth_v1/auth/login", bytes.NewBuffer([]byte{}))
		ctx.Request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, ctx.Request)
		responseData, _ := io.ReadAll(w.Body)
		assert.Equal(t, 400, w.Code)
		assert.Equal(t, `{"message":"errors: bad requset"}`, string(responseData))
	})

	t.Run("authUsecase.Login() fail", func(t *testing.T) {
		usecaseMock.ExpectedCalls = nil
		mockErr := errors.New("error: login failed")
		w := httptest.NewRecorder()
		usecaseMock.On("Login",
			mock.AnythingOfType("*config.Config"),
			mock.AnythingOfType("*auth.LoginReq")).Return(nil, mockErr)
		ctx.Request = httptest.NewRequest(http.MethodPost, "/auth_v1/auth/login", bytes.NewBuffer(byteLoginReq))
		ctx.Request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, ctx.Request)
		responseData, _ := io.ReadAll(w.Body)
		assert.Equal(t, 400, w.Code)
		assert.Equal(t, `{"message":"error: login failed"}`, string(responseData))
	})
}

func TestLogout(t *testing.T) {
	ctx, w, router, usecaseMock := setup()

	logoutReq := &auth.LogoutReq{CredentialID: 1}

	byteLogoutReq, err := json.Marshal(logoutReq)
	if err != nil {
		panic(err)
	}
	t.Run("success", func(t *testing.T) {
		usecaseMock.On("Logout", mock.AnythingOfType("*auth.LogoutReq")).Return(nil)
		ctx.Request = httptest.NewRequest(http.MethodPost, "/auth_v1/auth/logout", bytes.NewBuffer(byteLogoutReq))
		ctx.Request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, ctx.Request)
		assert.Equal(t, 200, w.Code)
	})

	t.Run("bind request fail", func(t *testing.T) {
		usecaseMock.ExpectedCalls = nil
		w := httptest.NewRecorder()
		ctx.Request = httptest.NewRequest(http.MethodPost, "/auth_v1/auth/logout", bytes.NewBuffer([]byte{}))
		ctx.Request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, ctx.Request)
		responseData, _ := io.ReadAll(w.Body)
		assert.Equal(t, 400, w.Code)
		assert.Equal(t, `{"message":"errors: bad requset"}`, string(responseData))
	})

	t.Run("authUsecase.Logout() fail", func(t *testing.T) {
		usecaseMock.ExpectedCalls = nil
		w := httptest.NewRecorder()
		usecaseMock.On("Logout", mock.AnythingOfType("*auth.LogoutReq")).Return(errors.New("error: logout fail"))
		ctx.Request = httptest.NewRequest(http.MethodPost, "/auth_v1/auth/logout", bytes.NewBuffer(byteLogoutReq))
		ctx.Request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, ctx.Request)
		responseData, _ := io.ReadAll(w.Body)
		assert.Equal(t, 400, w.Code)
		assert.Equal(t, `{"message":"error: logout fail"}`, string(responseData))
	})
}

func TestRefreshToken(t *testing.T) {
	ctx, w, router, usecaseMock := setup()
	refreshToken := "refresh_token"
	refreshTokenReq := auth.RefreshTokenReq{
		CredentialID: 1,
		RefreshToken: refreshToken,
	}

	credentialRes := &auth.CredentialRes{
		ID:           1,
		UserID:       1,
		AccessToken:  "access_token",
		RefreshToken: refreshToken,
		CreatedAt:    time.Time{},
		UpdatedAt:    time.Time{},
	}

	byteRefreshTokenReq, _ := json.Marshal(refreshTokenReq)
	t.Run("success", func(t *testing.T) {
		usecaseMock.On("RefreshToken", mock.AnythingOfType("*config.Config"), mock.AnythingOfType("*auth.RefreshTokenReq")).Return(credentialRes, nil)
		ctx.Request = httptest.NewRequest(http.MethodPost, "/auth_v1/auth/refresh-token", bytes.NewBuffer(byteRefreshTokenReq))
		ctx.Request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, ctx.Request)
		credentialRes, _ := json.Marshal(credentialRes)
		assert.Equal(t, string(credentialRes), string(w.Body.String()))
		assert.Equal(t, 200, w.Code)
	})

	t.Run("bind request fail", func(t *testing.T) {
		usecaseMock.ExpectedCalls = nil
		w := httptest.NewRecorder()
		ctx.Request = httptest.NewRequest(http.MethodPost, "/auth_v1/auth/refresh-token", bytes.NewBuffer([]byte{}))
		ctx.Request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, ctx.Request)
		responseData, _ := io.ReadAll(w.Body)
		assert.Equal(t, 400, w.Code)
		assert.Equal(t, `{"message":"errors: bad requset"}`, string(responseData))
	})

	t.Run("authUsecase.RefreshToken() fail", func(t *testing.T) {
		usecaseMock.ExpectedCalls = nil
		mockErr := errors.New("error: refresh token failed")
		w := httptest.NewRecorder()
		usecaseMock.On(
			"RefreshToken",
			mock.AnythingOfType("*config.Config"),
			mock.AnythingOfType("*auth.RefreshTokenReq"),
		).Return(nil, mockErr)
		ctx.Request = httptest.NewRequest(
			http.MethodPost,
			"/auth_v1/auth/refresh-token",
			bytes.NewBuffer(byteRefreshTokenReq),
		)
		ctx.Request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, ctx.Request)
		responseData, _ := io.ReadAll(w.Body)
		assert.Equal(t, 400, w.Code)
		assert.Equal(t, `{"message":"error: refresh token failed"}`, string(responseData))
	})
}
