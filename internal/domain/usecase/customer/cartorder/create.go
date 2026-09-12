package cartorder

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
)

func (c *CartOrder) Create(
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
