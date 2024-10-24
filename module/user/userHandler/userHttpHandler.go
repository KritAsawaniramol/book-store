package userHandler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/user"
	"github.com/kritAsawaniramol/book-store/module/user/userUsecase"
	"github.com/kritAsawaniramol/book-store/pkg/request"
	"github.com/stripe/stripe-go/v80"
	"github.com/stripe/stripe-go/v80/webhook"
)

type UserHttpHandler interface {
	Register(ctx *gin.Context)
	AddUserTransaction(ctx *gin.Context)
	GetUserBalance(ctx *gin.Context)
	GetUserProfile(ctx *gin.Context)
	SearchUserTransaction(ctx *gin.Context)
	TopUp(ctx *gin.Context)
	GetOneTopUpOrder(ctx *gin.Context)
	StripeWebhook(ctx *gin.Context)
}

type userHttpHandlerImpl struct {
	userUsecase userUsecase.UserUsecase
	cfg         *config.Config
}

// StripeWebhook implements UserHttpHandler.
func (u *userHttpHandlerImpl) StripeWebhook(ctx *gin.Context) {
	stripe.Key = u.cfg.Stripe.SecretKey

	payload, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Replace this endpoint secret with your endpoint's unique secret
	// If you are testing with the CLI, find the secret by running 'stripe listen'
	// If you are using an endpoint defined with the API or dashboard, look in your webhook settings
	// at https://dashboard.stripe.com/webhooks
	endpointSecret := u.cfg.Stripe.EndPointSecret
	signatureHeader := ctx.Request.Header.Get("Stripe-Signature")
	event, err := webhook.ConstructEvent(payload, signatureHeader, endpointSecret)
	if err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Webhook signature verification failed. %v\n", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Webhook signature verification failed"})
		return
	}
	switch event.Type {
	case "checkout.session.completed":
		paymentData := event.Data.Object
		sessionID := paymentData["id"].(string)
		sessionStatus := paymentData["status"].(string)
		if err := u.userUsecase.HandleStripeWebhook(sessionID, sessionStatus); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		}

	default:
		fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
	}

	// Unmarshal the event data into an appropriate struct depending on its Type

	ctx.Status(http.StatusOK)

}

// GetOneTopUpOrder implements UserHttpHandler.
func (u *userHttpHandlerImpl) GetOneTopUpOrder(ctx *gin.Context) {
	topUpIDStr := ctx.Param("id")
	topUpID, err := strconv.ParseUint(topUpIDStr, 10, 64)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	res, err := u.userUsecase.GetOneTopUpOrderByID(uint(topUpID))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

// TopUp implements UserHttpHandler.
func (u *userHttpHandlerImpl) TopUp(ctx *gin.Context) {
	wrapper := request.ContextWrapper(ctx)
	req := &user.TopUpReq{}
	if err := wrapper.Bind(req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	req.UserID = ctx.GetUint("userID")
	sessionID, err := u.userUsecase.TopUp(req, u.cfg)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"sessionID": sessionID})
}

// SearchUserTransaction implements UserHttpHandler.
func (u *userHttpHandlerImpl) SearchUserTransaction(ctx *gin.Context) {
	userIDStr := ctx.Query("user_id")
	var userID uint = 0
	if userIDStr != "" {
		userIDUint64, err := strconv.ParseUint(userIDStr, 10, 64)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		userID = uint(userIDUint64)
	}
	req := &user.SearchUserTransactionReq{UsersID: userID}
	res, err := u.userUsecase.SearchUserTransaction(req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

// GetUserProfile implements UserHttpHandler.
func (u *userHttpHandlerImpl) GetUserProfile(ctx *gin.Context) {
	userProfile, err := u.userUsecase.GetUserProfile(ctx.GetUint("userID"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, userProfile)
}

// GetUserBalance implements UserHttpHandler.
func (u *userHttpHandlerImpl) GetUserBalance(ctx *gin.Context) {
	balance, err := u.userUsecase.GetUserBalance(ctx.GetUint("userID"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, balance)
}

// AddUserTransaction implements UserHttpHandler.
func (u *userHttpHandlerImpl) AddUserTransaction(ctx *gin.Context) {
	wrapper := request.ContextWrapper(ctx)
	req := &user.CreateUserTransactionReq{}
	if err := wrapper.Bind(req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	balance, err := u.userUsecase.CreateUserTransaction(req, "admin add user money.")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, balance)
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

func NewUserHttpHandler(userUsecase userUsecase.UserUsecase, cfg *config.Config) UserHttpHandler {
	return &userHttpHandlerImpl{userUsecase: userUsecase, cfg: cfg}
}
