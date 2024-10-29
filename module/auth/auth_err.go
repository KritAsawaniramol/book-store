package auth

import "errors"

var (
	ErrCredentialNotFound       = errors.New("error: credential not found")
	ErrDeleteUserCredentialFail = errors.New("error: delete user credential failed")
	ErrUserNotFound             = errors.New("error: user not found")
	ErrUpdateCredential         = errors.New("error: update credential failed")
	ErrCreateUserCredential     = errors.New("error: create user credential failed")
	ErrEmailOrPasswordIncorrect = errors.New("error: email or password are incorrect")
)
