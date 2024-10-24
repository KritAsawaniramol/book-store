package orderUsecase

import (
	"fmt"
	"sort"

	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/book/bookPb"
	"github.com/kritAsawaniramol/book-store/module/order"
	"github.com/kritAsawaniramol/book-store/module/order/orderRepository"
	"github.com/kritAsawaniramol/book-store/module/shelf"
	"github.com/kritAsawaniramol/book-store/module/user"
	"github.com/kritAsawaniramol/book-store/util"
)

type orderUsecaseImpl struct {
	orderRepository orderRepository.OrderRepository
}

// GetMyOrders implements OrderUsecase.
func (o *orderUsecaseImpl) GetMyOrders(cfg *config.Config, userID uint) (*order.GetMyOrdersRes, error) {
	result, err := o.orderRepository.GetOrders([]uint{}, []uint{userID}, []string{}, true)
	if err != nil {
		return nil, err
	}

	util.PrintObjInJson(result)

	bookIDs := []uint64{}
	uniqueBooks := map[uint]*bookPb.Book{}

	for _, r := range result {
		for _, ob := range r.OrdersBooks {
			if _, ok := uniqueBooks[ob.BookID]; ok {
				continue
			} else {
				uniqueBooks[ob.BookID] = &bookPb.Book{}
				bookIDs = append(bookIDs, uint64(ob.BookID))
			}
		}
	}

	findBooksInIdsRes, err := o.orderRepository.FindBookInIds(cfg.Grpc.BookUrl, &bookPb.FindBooksInIdsReq{
		Ids: bookIDs,
	})
	if err != nil {
		return nil, err
	}
	for _, b := range findBooksInIdsRes.Book {
		uniqueBooks[uint(b.Id)] = b
	}

	res := []order.GetMyOrdersResOrder{}

	for _, r := range result {
		orderRes := order.GetMyOrdersResOrder{
			OrderID:    r.ID,
			UserID:     r.UserID,
			TotalPrice: r.Total,
			Status:     r.Status,
			CreatedAt:  r.CreatedAt,
			UpdatedAt:  r.UpdatedAt,
		}

		bookDatum := []*order.GetMyOrdersResBookDatum{}
		for _, b := range r.OrdersBooks {
			if v, ok := uniqueBooks[b.BookID]; ok && v != nil {
				bookDatum = append(bookDatum, &order.GetMyOrdersResBookDatum{
					BookID: b.BookID,
					Title:  v.Title,
					Price:  b.Price,
				})
			}
		}

		orderRes.Books = bookDatum

		res = append(res, orderRes)
	}

	return &order.GetMyOrdersRes{Orders: res}, nil
}

func (o *orderUsecaseImpl) CanAddBooksToCart(req *order.SearchOneMyOrderReq) bool {
	return false
}

// SearchOrder implements OrderUsecase.
func (o *orderUsecaseImpl) SearchOneUserOrderByBookID(req *order.SearchOneMyOrderReq) (*order.SearchOneMyOrderRes, error) {
	orders, err := o.orderRepository.FindOrdersWithBookDetail(
		&order.Orders{UserID: req.UserID},
		[]string{order.Completed, order.Pendding})
	if err != nil {
		return nil, err
	}
	for _, v := range orders {
		sort.Slice(v.OrdersBooks, func(i, j int) bool {
			return v.OrdersBooks[i].BookID < v.OrdersBooks[j].BookID
		})
		//binary search
		i := sort.Search(len(v.OrdersBooks), func(i int) bool { return v.OrdersBooks[i].BookID >= req.BookID })
		if i < len(v.OrdersBooks) && v.OrdersBooks[i].BookID == req.BookID {
			bookResDatum := []*order.BookResDatum{}
			for _, b := range v.OrdersBooks {
				bookResDatum = append(bookResDatum, &order.BookResDatum{
					BookID: b.BookID,
					Price:  b.Price,
				})
			}
			return &order.SearchOneMyOrderRes{
				OrderID:    v.ID,
				UserID:     v.UserID,
				Books:      bookResDatum,
				TotalPrice: v.Total,
				Status:     v.Status,
			}, nil
		}
	}
	return nil, nil
}

// HandleAddBookRes implements OrderUsecase.
func (o *orderUsecaseImpl) HandleAddBookRes(cfg *config.Config, res *shelf.AddBooksRes) {
	if res.Error != "" {
		o.orderRepository.UpdateOrderByID(
			res.OrderID, &order.Orders{Status: order.Failed, Note: res.Error},
		)
		o.orderRepository.RollbackUserTransaction(cfg, &user.RollbackUserTransactionReq{
			TransactionID: res.TransactionID,
		})
		return
	}

	if err := o.orderRepository.UpdateOrderByID(
		res.OrderID, &order.Orders{Status: order.Completed},
	); err != nil {
		o.orderRepository.RollbackUserTransaction(cfg, &user.RollbackUserTransactionReq{
			TransactionID: res.TransactionID,
		})
		o.orderRepository.RollbackAddBooks(cfg, &shelf.RollbackAddBooks{
			ShelfIDs: res.ShelfIDs,
		})
	}
}

// HandleBuyBooksRes implements OrderUsecase.
func (o *orderUsecaseImpl) HandleBuyBooksRes(cfg *config.Config, res *user.BuyBookRes) {
	if res.Error != "" {
		o.orderRepository.UpdateOrderByID(
			res.OrderID, &order.Orders{Status: order.Failed, Note: res.Error},
		)
		return
	}

	// add book to user shelf
	if err := o.orderRepository.AddBookToShelf(cfg, &shelf.AddBooksReq{
		OrderID:       res.OrderID,
		UserID:        res.UserID,
		TransactionID: res.TransactionID,
		BookIDs:       res.BookIDs,
	}); err != nil {
		//rollback transaction
		o.orderRepository.RollbackUserTransaction(cfg, &user.RollbackUserTransactionReq{
			TransactionID: res.TransactionID,
		})
	}
}

// BuyBooks implements OrderUsecase.
func (o *orderUsecaseImpl) BuyBooks(cfg *config.Config, req *order.BuyBooksReq) (*order.BuyBooksRes, error) {
	bookIDsUint64 := []uint64{}
	for _, b := range req.BookIDs {
		bookIDsUint64 = append(bookIDsUint64, uint64(b))
	}

	// check book isn't in completed or pendding order of this user
	orders, err := o.orderRepository.FindOrdersWithBookDetail(
		&order.Orders{UserID: req.UserID},
		[]string{order.Completed, order.Pendding})
	if err != nil {
		return nil, err
	}

	mapBookIDsInOrder := map[uint]struct{}{}
	for _, v := range req.BookIDs {
		mapBookIDsInOrder[v] = struct{}{}
	}

	for _, v := range orders {
		for _, orderBook := range v.OrdersBooks {
			if _, ok := mapBookIDsInOrder[orderBook.BookID]; ok {
				return nil, fmt.Errorf("error: book id %v already in another order (pendding, completed)", orderBook.BookID)
			}
		}
	}

	//Get book info
	fmt.Printf("bookIDsUint64: %v\n", bookIDsUint64)
	booksInfo, err := o.orderRepository.FindBookInIds(cfg.Grpc.BookUrl, &bookPb.FindBooksInIdsReq{Ids: bookIDsUint64})
	if err != nil {
		return nil, err
	}

	fmt.Printf("booksInfo: %v\n", booksInfo)

	ordersBooks := []order.OrdersBooks{}
	var totalPrice uint = 0
	for _, b := range booksInfo.Book {
		totalPrice += uint(b.Price)
		ordersBooks = append(ordersBooks, order.OrdersBooks{
			BookID: uint(b.Id),
			Price:  uint(b.Price),
		})
	}

	//create createOrderReq in database, status is "pendding"
	createOrderReq := &order.Orders{
		UserID:      req.UserID,
		Status:      order.Pendding,
		OrdersBooks: ordersBooks,
		Total:       totalPrice,
	}
	if err := o.orderRepository.CreateOrder(createOrderReq); err != nil {
		return nil, err
	}

	//decress user money
	decressUserMoneyReq := &user.BuyBookReq{
		OrderID: createOrderReq.ID,
		UserID:  req.UserID,
		BookIDs: req.BookIDs,
		Total:   totalPrice,
	}
	if err := o.orderRepository.DecreaseUserMoney(cfg, decressUserMoneyReq); err != nil {
		o.orderRepository.UpdateOrderByID(
			createOrderReq.ID,
			&order.Orders{Status: order.Failed, Note: err.Error()},
		)
		return nil, err
	}
	return nil, nil
}

func NewOrderUsecaseImpl(orderRepository orderRepository.OrderRepository) OrderUsecase {
	return &orderUsecaseImpl{orderRepository: orderRepository}
}
