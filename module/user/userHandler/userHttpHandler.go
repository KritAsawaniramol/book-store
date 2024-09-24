package userHandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kritAsawaniramol/book-store/module/user"
	"github.com/kritAsawaniramol/book-store/module/user/userUsecase"
	"github.com/kritAsawaniramol/book-store/pkg/request"
)

type UserHttpHandler interface {
	Register(ctx *gin.Context)
}

type userHttpHandlerImpl struct {
	userUsecase userUsecase.UserUsecase
}

// register implements UserHttpHandler.
func (u *userHttpHandlerImpl) Register(ctx *gin.Context) {
	wrapper := request.ContextWrapper(ctx)
	req := &user.UserRegisterReq{}
	if err := wrapper.Bind(req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	userID, err := u.userUsecase.Register(req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"user_id": userID})
}

func NewUserHttpHandler(userUsecase userUsecase.UserUsecase) UserHttpHandler {
	return &userHttpHandlerImpl{userUsecase: userUsecase}
}
