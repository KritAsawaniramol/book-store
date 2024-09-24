package userRepository

import (
	"github.com/kritAsawaniramol/book-store/module/user"
)

type UserRepository interface {
	CreateOneUser(in *user.User) (uint, error)
}
