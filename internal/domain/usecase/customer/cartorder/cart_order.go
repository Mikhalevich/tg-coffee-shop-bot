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

type CurrencyService interface {
	GetCurrencyByID(
		ctx context.Context,
		id currency.ID,
	) (*currency.Currency, error)
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
		curr *currency.Currency,
	) error
	ViewCategoryProducts(
		ctx context.Context,
		chatID msginfo.ChatID,
		messageID msginfo.MessageID,
		cartID cart.ID,
		categoryID product.CategoryID,
		categoryProducts []product.Product,
		cartProducts []cart.CartProduct,
		curr *currency.Currency,
	) error
	CartOrderUnavailable(
		ctx context.Context,
		chatID msginfo.ChatID,
	) error
}

type CartOrder struct {
	storeService        StoreService
	productService      ProductService
	cartService         CartService
	currencyService     CurrencyService
	notificationService NotificationService
}

func New(
	storeService StoreService,
	productService ProductService,
	cartService CartService,
	currencyService CurrencyService,
	notificationService NotificationService,
) *CartOrder {
	return &CartOrder{
		storeService:        storeService,
		productService:      productService,
		cartService:         cartService,
		currencyService:     currencyService,
		notificationService: notificationService,
	}
}
