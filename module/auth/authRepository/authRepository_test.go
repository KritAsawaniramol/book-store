package authRepository

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kritAsawaniramol/book-store/module/auth"
	"github.com/kritAsawaniramol/book-store/module/user/userPb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGetOneUserCredential(t *testing.T) {
	// Mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))
	if err != nil {
		t.Fatalf("Failed to open gorm database: %v", err)
	}

	userGrpcClientMock := userPb.NewUserGrpcServiceClientMock()
	repo := NewAuthRepositoryImpl(gormDB, userGrpcClientMock)
	testCredential := &auth.Credential{UserID: 1}
	t.Run("success", func(t *testing.T) {
		expectedResult := auth.Credential{
			UserID:       testCredential.ID,
			RoleID:       1,
			AccessToken:  "access_token",
			RefreshToken: "refresh_token",
		}

		rows := sqlmock.NewRows([]string{"user_id", "role_id", "access_token", "refresh_token"}).
			AddRow(expectedResult.UserID, expectedResult.RoleID, expectedResult.AccessToken, expectedResult.RefreshToken)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "credentials" WHERE "credentials"."user_id" = $1 AND "credentials"."deleted_at" IS NULL ORDER BY "credentials"."id" LIMIT $2`)).
			WithArgs(1, 1).
			WillReturnRows(rows)

		credential, err := repo.GetOneUserCredential(testCredential)
		assert.NoError(t, err)
		assert.Equal(t, expectedResult, *credential)
		assert.Equal(t, expectedResult.UserID, credential.UserID)
		assert.Equal(t, expectedResult.RoleID, credential.RoleID)
		assert.Equal(t, expectedResult.AccessToken, credential.AccessToken)
		assert.Equal(t, expectedResult.RefreshToken, credential.RefreshToken)
	})

	t.Run(".First fail", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "credentials" WHERE "credentials"."user_id" = $1 AND "credentials"."deleted_at" IS NULL ORDER BY "credentials"."id" LIMIT $2`)).
			WithArgs(1, 1).
			WillReturnError(auth.ErrCredentialNotFound)

		credential, err := repo.GetOneUserCredential(testCredential)
		assert.Equal(t, auth.ErrCredentialNotFound, err)
		assert.Nil(t, credential)

		// Ensure all expectations were met
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCreateOneUserCredential(t *testing.T) {
	// Mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))
	if err != nil {
		t.Fatalf("Failed to open gorm database: %v", err)
	}
	userGrpcClientMock := userPb.NewUserGrpcServiceClientMock()
	repo := NewAuthRepositoryImpl(gormDB, userGrpcClientMock)

	testCredential := &auth.Credential{
		UserID:       1,
		RoleID:       1,
		AccessToken:  "access_token",
		RefreshToken: "refresh_token",
	}

	t.Run("success", func(t *testing.T) {
		mockCredentialID := uint(1)
		mock.ExpectBegin()
		mock.ExpectQuery("^INSERT INTO \"credentials\" (.+)$").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(mockCredentialID))
		mock.ExpectCommit()

		credentialID, err := repo.CreateOneUserCredential(testCredential)
		assert.NoError(t, err)
		assert.Equal(t, mockCredentialID, credentialID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run(".Create() fail", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery("^INSERT INTO \"credentials\" (.+)$").
			WillReturnError(auth.ErrCreateUserCredential)
		mock.ExpectRollback()

		credentialID, err := repo.CreateOneUserCredential(testCredential)
		assert.Error(t, err)
		assert.Equal(t, uint(0), credentialID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUpdateOneCredentialByID(t *testing.T) {
	// Mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))
	if err != nil {
		t.Fatalf("Failed to open gorm database: %v", err)
	}

	userGrpcClientMock := userPb.NewUserGrpcServiceClientMock()
	repo := NewAuthRepositoryImpl(gormDB, userGrpcClientMock)
	mockCredentialID := uint(1)
	testCredential := &auth.Credential{
		UserID:       1,
		RoleID:       1,
		AccessToken:  "access_token",
		RefreshToken: "refresh_token",
	}

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("^UPDATE \"credentials\" SET (.+)$").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		err := repo.UpdateOneCredentialByID(mockCredentialID, testCredential)
		assert.NoError(t, err)
	})
	t.Run(".Updates() fail", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("^UPDATE \"credentials\" SET (.+)$").WillReturnError(auth.ErrUpdateCredential)
		mock.ExpectRollback()
		err := repo.UpdateOneCredentialByID(mockCredentialID, testCredential)
		assert.Error(t, err)
	})
}

func TestDeleteOneUserCredentialByID(t *testing.T) {
	// Mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))
	if err != nil {
		t.Fatalf("Failed to open gorm database: %v", err)
	}

	userGrpcClientMock := userPb.NewUserGrpcServiceClientMock()
	repo := NewAuthRepositoryImpl(gormDB, userGrpcClientMock)
	mockCredentialID := uint(1)

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("^UPDATE \"credentials\" SET (.+)$").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		err := repo.DeleteOneUserCredentialByID(mockCredentialID)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run(".Delete() fail", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("^UPDATE \"credentials\" SET (.+)$").WillReturnError(auth.ErrUpdateCredential)
		mock.ExpectRollback()
		err := repo.DeleteOneUserCredentialByID(mockCredentialID)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestFindOneUserProfileToRefresh(t *testing.T) {
	userGrpcClientMock := userPb.NewUserGrpcServiceClientMock()
	mockReq := &userPb.FindOneUserProfileToRefreshReq{UserId: 1}
	createAt := timestamppb.New(time.Time{})
	updatedAt := timestamppb.New(time.Time{})
	mockRes := &userPb.UserProfile{
		Id:        1,
		Username:  "user1",
		RoleId:    1,
		Coin:      20,
		CreatedAt: createAt,
		UpdatedAt: updatedAt,
	}

	repo := NewAuthRepositoryImpl(nil, userGrpcClientMock)
	t.Run("success", func(t *testing.T) {
		userGrpcClientMock.On("FindOneUserProfileToRefresh", mock.Anything, mock.AnythingOfType("*userPb.FindOneUserProfileToRefreshReq")).
			Return(mockRes, nil)
		userProfile, err := repo.FindOneUserProfileToRefresh(mockReq)
		assert.NoError(t, err)
		assert.Equal(t, mockRes, userProfile)
	})
	t.Run("FindOneUserProfileToRefresh fail", func(t *testing.T) {
		userGrpcClientMock.ExpectedCalls = nil
		userGrpcClientMock.On("FindOneUserProfileToRefresh", mock.Anything, mock.AnythingOfType("*userPb.FindOneUserProfileToRefreshReq")).
			Return(nil, auth.ErrUserNotFound)
		userProfile, err := repo.FindOneUserProfileToRefresh(mockReq)
		assert.Error(t, err)
		assert.Nil(t, userProfile)
	})
}
func TestFindOneUserProfileToLogin(t *testing.T) {
	userGrpcClientMock := userPb.NewUserGrpcServiceClientMock()
	mockReq := &userPb.FindUserProfileToLoginReq{Username: "user1", Password: "password"}
	createAt := timestamppb.New(time.Time{})
	updatedAt := timestamppb.New(time.Time{})
	mockRes := &userPb.UserProfile{
		Id:        1,
		Username:  "user1",
		RoleId:    1,
		Coin:      20,
		CreatedAt: createAt,
		UpdatedAt: updatedAt,
	}

	repo := NewAuthRepositoryImpl(nil, userGrpcClientMock)
	t.Run("success", func(t *testing.T) {
		userGrpcClientMock.On("FindUserProfileToLogin", mock.Anything, mock.AnythingOfType("*userPb.FindUserProfileToLoginReq")).
			Return(mockRes, nil)
		userProfile, err := repo.FindOneUserProfileToLogin(mockReq)
		assert.NoError(t, err)
		assert.Equal(t, mockRes, userProfile)
	})
	t.Run("FindUserProfileToLogin fail", func(t *testing.T) {
		userGrpcClientMock.ExpectedCalls = nil
		userGrpcClientMock.On("FindUserProfileToLogin", mock.Anything, mock.AnythingOfType("*userPb.FindUserProfileToLoginReq")).
			Return(nil, auth.ErrUserNotFound)
		userProfile, err := repo.FindOneUserProfileToLogin(mockReq)
		assert.Error(t, err)
		assert.Nil(t, userProfile)
	})
}