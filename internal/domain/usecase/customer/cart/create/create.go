package create

import (
	"context"
	"fmt"
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
}

type CartService interface {
	Create(
		ctx context.Context,
		chatID msginfo.ChatID,
	) (cart.ID, error)
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
}

type Create struct {
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
) *Create {
	return &Create{
		storeService:        storeService,
		productService:      productService,
		cartService:         cartService,
		currencyService:     currencyService,
		notificationService: notificationService,
	}
}

func (c *Create) Create(
	ctx context.Context,
	info msginfo.Info,
) error {
	storeInfo, err := c.storeService.GetStoreInfo(ctx)
	if err != nil {
		return fmt.Errorf("get store info: %w", err)
	}

	if !storeInfo.IsActive {
		if err := c.notificationService.SendStoreClosed(
			ctx,
			info.ChatID,
			storeInfo.CurrentTime,
			storeInfo.NextWorkingTime,
		); err != nil {
			return fmt.Errorf("send store closed msg: %w", err)
		}
	}

	categories, err := c.productService.GetCategories(ctx)
	if err != nil {
		return fmt.Errorf("get product categories: %w", err)
	}

	cartID, err := c.cartService.Create(ctx, info.ChatID)
	if err != nil {
		return fmt.Errorf("create cart: %w", err)
	}

	curr, err := c.currencyService.GetCurrencyByID(ctx, storeInfo.DefaultCurrencyID)
	if err != nil {
		return fmt.Errorf("get store currency: %w", err)
	}

	if err := c.notificationService.ViewOrderCategories(
		ctx,
		info.ChatID,
		info.MessageID,
		cartID,
		categories,
		nil,
		curr,
	); err != nil {
		return fmt.Errorf("view order categories msg: %w", err)
	}

	return nil
}
