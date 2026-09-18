package cartorder

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/cart"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/product"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/store"
)

type Transactor interface {
	Transaction(ctx context.Context, trxFn func(ctx context.Context) error) error
}

type StoreService interface {
	GetStoreInfo(ctx context.Context) (store.StoreInfo, error)
}

type ProductService interface {
	GetCategories(ctx context.Context) ([]product.Category, error)
	GetProductsByCategoryID(
		ctx context.Context,
		categoryID product.CategoryID,
		currencyID currency.ID,
	) ([]product.Product, error)
	GetProductsByIDs(
		ctx context.Context,
		ids []product.ProductID,
		currencyID currency.ID,
	) (map[product.ProductID]product.Product, error)
}

type CartService interface {
	Create(
		ctx context.Context,
		chatID msginfo.ChatID,
	) (cart.ID, error)
	Clear(ctx context.Context, chatID msginfo.ChatID, cartID cart.ID) error
	AddProduct(
		ctx context.Context,
		id cart.ID,
		p cart.CartProduct,
	) error
	GetProducts(ctx context.Context, id cart.ID) ([]cart.CartProduct, error)
}

type OrderService interface {
	CreateOrder(
		ctx context.Context,
		info order.CreateOrderInfo,
	) (order.Order, error)
}

type CurrencyService interface {
	GetCurrencyByID(
		ctx context.Context,
		id currency.ID,
	) (currency.Currency, error)
}

type NotificationService interface {
	SendStoreClosed(
		ctx context.Context,
		chatID msginfo.ChatID,
		currentTime time.Time,
		nextWorkingTime time.Time,
	) error
	ViewOrderCategories(
		ctx context.Context,
		chatID msginfo.ChatID,
		messageID msginfo.MessageID,
		cartID cart.ID,
		categories []product.Category,
		orderedProducts []order.OrderedProduct,
		curr currency.Currency,
	) error
	ViewCategoryProducts(
		ctx context.Context,
		chatID msginfo.ChatID,
		messageID msginfo.MessageID,
		cartID cart.ID,
		categoryID product.CategoryID,
		categoryProducts []product.Product,
		cartProducts []cart.CartProduct,
		curr currency.Currency,
	) error
	CartOrderUnavailable(
		ctx context.Context,
		chatID msginfo.ChatID,
	) error
	NoProductsForOrder(ctx context.Context, chatID msginfo.ChatID) error
	OrderAlreadyExists(ctx context.Context, chatID msginfo.ChatID) error
	SendInvoice(
		ctx context.Context,
		chatID msginfo.ChatID,
		ord order.Order,
		productsInfo map[product.ProductID]product.Product,
		curr currency.Currency,
	) error
}

type CartOrder struct {
	transactor          Transactor
	storeService        StoreService
	productService      ProductService
	cartService         CartService
	orderService        OrderService
	currencyService     CurrencyService
	notificationService NotificationService
}

func New(
	transator Transactor,
	storeService StoreService,
	productService ProductService,
	cartService CartService,
	orderService OrderService,
	currencyService CurrencyService,
	notificationService NotificationService,
) *CartOrder {
	return &CartOrder{
		transactor:          transator,
		storeService:        storeService,
		productService:      productService,
		cartService:         cartService,
		orderService:        orderService,
		currencyService:     currencyService,
		notificationService: notificationService,
	}
}
