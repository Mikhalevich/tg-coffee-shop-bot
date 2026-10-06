package tghandler

import (
	"context"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/cart"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/product"
)

type CartUsecase interface {
	Create(ctx context.Context, info msginfo.Info) error
	ViewCategoryProducts(ctx context.Context, info msginfo.Info, cartID cart.ID, categoryID product.CategoryID,
		currencyID currency.ID) error
	ViewCategories(ctx context.Context, info msginfo.Info, cartID cart.ID, currencyID currency.ID) error
	Add(
		ctx context.Context,
		info msginfo.Info,
		cartID cart.ID,
		categoryID product.CategoryID,
		productID product.ProductID,
		currencyID currency.ID,
	) error
	Cancel(
		ctx context.Context,
		chatID msginfo.ChatID,
		cartID cart.ID,
	) error
	Confirm(ctx context.Context, info msginfo.Info, cartID cart.ID, currencyID currency.ID) error
}

type ViewActiveOrderUsecase interface {
	ViewActiveOrder(ctx context.Context, chatID msginfo.ChatID) error
}

type CancelOrderUsecase interface {
	Cancel(
		ctx context.Context,
		info msginfo.Info,
		orderID order.ID,
	) error
}

type QueueSizeUsecase interface {
	Size(ctx context.Context, chatID msginfo.ChatID) error
}

type HistoryOrderUsecase interface {
	First(ctx context.Context, info msginfo.Info) error
	Last(ctx context.Context, info msginfo.Info) error
	Page(ctx context.Context, info msginfo.Info, pageNumber int) error
}

type OrderPaymentUsecase interface {
	InProgress(
		ctx context.Context,
		paymentID string,
		orderID order.ID,
		totalAmount int,
	) error
	Confirmed(
		ctx context.Context,
		chatID msginfo.ChatID,
		orderID order.ID,
		currency string,
		totalAmount int,
	) error
}

type ButtonProvider interface {
	GetButton(ctx context.Context, id button.ID) (*button.Button, error)
}

type cbHandler func(ctx context.Context, info msginfo.Info, btn button.Button) error

type TGHandler struct {
	cartUsecase            CartUsecase
	viewActiveOrderUsecase ViewActiveOrderUsecase
	cancelOrderUsecase     CancelOrderUsecase
	queueSizeUsecase       QueueSizeUsecase
	historyOrderUsecase    HistoryOrderUsecase
	orderPaymentUsecase    OrderPaymentUsecase
	buttonProvider         ButtonProvider

	cbHandlers map[button.Operation]cbHandler
}

func New(
	cartUsecase CartUsecase,
	viewActiveOrderUsecase ViewActiveOrderUsecase,
	cancelOrderUsecase CancelOrderUsecase,
	queueSizeUsecase QueueSizeUsecase,
	historyOrderUsecase HistoryOrderUsecase,
	orderPaymentUsecase OrderPaymentUsecase,
	buttonProvider ButtonProvider,
) *TGHandler {
	handler := &TGHandler{
		cartUsecase:            cartUsecase,
		viewActiveOrderUsecase: viewActiveOrderUsecase,
		cancelOrderUsecase:     cancelOrderUsecase,
		queueSizeUsecase:       queueSizeUsecase,
		historyOrderUsecase:    historyOrderUsecase,
		orderPaymentUsecase:    orderPaymentUsecase,
		buttonProvider:         buttonProvider,
	}

	handler.initCBHandlers()

	return handler
}

func (t *TGHandler) initCBHandlers() {
	t.cbHandlers = map[button.Operation]cbHandler{
		button.OperationOrderCancel:              t.cancelOrder,
		button.OperationCartCancel:               t.cancelCart,
		button.OperationCartConfirm:              t.confirmCart,
		button.OperationCartViewCategories:       t.viewCategories,
		button.OperationCartViewCategoryProducts: t.viewCategoryProducts,
		button.OperationCartAddProduct:           t.addProduct,

		button.OperationOrderHistoryByPageFirst: t.historyFirstV2,
		button.OperationOrderHistoryByPageLast:  t.historyLastV2,
		button.OperationOrderHistoryByPage:      t.historyPageV2,
	}
}
