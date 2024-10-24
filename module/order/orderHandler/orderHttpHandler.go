package orderHandler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/order"
	"github.com/kritAsawaniramol/book-store/module/order/orderUsecase"
	"github.com/kritAsawaniramol/book-store/pkg/request"
)

type (
	OrderHttpHandler interface {
		BuyBooks(ctx *gin.Context)
		SearchOneMyOrder(ctx *gin.Context)
		GetMyOrders(ctx *gin.Context)
	}

	orderHttpHandlerImpl struct {
		orderUsecase orderUsecase.OrderUsecase
		cfg          *config.Config
	}
)

// GetMyOrders implements OrderHttpHandler.
func (o *orderHttpHandlerImpl) GetMyOrders(ctx *gin.Context) {
	res, err := o.orderUsecase.GetMyOrders(o.cfg, ctx.GetUint("userID"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "query \"book_id=\" is required"})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

// SearchMyOrder implements OrderHttpHandler.
func (o *orderHttpHandlerImpl) SearchOneMyOrder(ctx *gin.Context) {
	bookIDStr := ctx.Query("book_id")
	if bookIDStr == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "query \"book_id=\" is required"})
		return
	}
	bookIDUint64, err := strconv.ParseUint(bookIDStr, 10, 64)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	req := &order.SearchOneMyOrderReq{
		UserID: ctx.GetUint("userID"),
		BookID: uint(bookIDUint64),
	}
	orders, err := o.orderUsecase.SearchOneUserOrderByBookID(req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, orders)
}

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
	ctx.JSON(http.StatusOK, gin.H{"message": "create order success"})
}

func NewOrderHttpHandlerImpl(cfg *config.Config, orderUsecase orderUsecase.OrderUsecase) OrderHttpHandler {
	return &orderHttpHandlerImpl{
		orderUsecase: orderUsecase,
		cfg:          cfg,
	}
}
