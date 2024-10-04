package orderHandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/order"
	"github.com/kritAsawaniramol/book-store/module/order/orderUsecase"
	"github.com/kritAsawaniramol/book-store/pkg/request"
)

type (
	OrderHttpHandler interface {
		BuyBooks(ctx *gin.Context)
	}

	orderHttpHandlerImpl struct {
		orderUsecase orderUsecase.OrderUsecase
		cfg          *config.Config
	}
)

// BuyBooks implements OrderHttpHandler.
func (o *orderHttpHandlerImpl) BuyBooks(ctx *gin.Context) {
	wrapper := request.ContextWrapper(ctx)
	req := &order.BuyBooksReq{}
	if err := wrapper.Bind(req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	req.UserID = ctx.GetUint("userID")
	if _, err := o.orderUsecase.BuyBooks(o.cfg, req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "create order"})
}

func NewOrderHttpHandlerImpl(cfg *config.Config, orderUsecase orderUsecase.OrderUsecase) OrderHttpHandler {
	return &orderHttpHandlerImpl{
		orderUsecase: orderUsecase,
		cfg:          cfg,
	}
}
