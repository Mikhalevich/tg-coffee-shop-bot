package cartorder

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/cart"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/perror"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/product"
)

func (c *CartOrder) ViewCategories(
	ctx context.Context,
	info msginfo.Info,
	cartID cart.ID,
	currencyID currency.ID,
) error {
	orderedProducts, _, err := c.orderedProductsFromCart(ctx, cartID, currencyID)
	if err != nil {
		if !perror.IsType(err, perror.TypeNotFound) {
			return fmt.Errorf("get ordered products: %w", err)
		}

		if err := c.notificationService.CartOrderUnavailable(ctx, info.ChatID); err != nil {
			return fmt.Errorf("order unavailable msg: %w", err)
		}

		return nil
	}

	categories, err := c.productService.GetCategories(ctx)
	if err != nil {
		return fmt.Errorf("get products: %w", err)
	}

	curr, err := c.currencyService.GetCurrencyByID(ctx, currencyID)
	if err != nil {
		return fmt.Errorf("get currency by id: %w", err)
	}

	if err := c.notificationService.ViewOrderCategories(
		ctx,
		info.ChatID,
		info.MessageID,
		cartID,
		categories,
		orderedProducts,
		curr,
	); err != nil {
		return fmt.Errorf("view order categories msg: %w", err)
	}

	return nil
}

func (c *CartOrder) orderedProductsFromCart(
	ctx context.Context,
	cartID cart.ID,
	currencyID currency.ID,
) ([]order.OrderedProduct, map[product.ProductID]product.Product, error) {
	cartProducts, err := c.cartService.GetProducts(ctx, cartID)
	if err != nil {
		return nil, nil, fmt.Errorf("get cart products: %w", err)
	}

	if len(cartProducts) == 0 {
		return nil, nil, nil
	}

	productIDs := make([]product.ProductID, 0, len(cartProducts))
	for _, v := range cartProducts {
		productIDs = append(productIDs, v.ProductID)
	}

	productsInfo, err := c.productService.GetProductsByIDs(ctx, productIDs, currencyID)
	if err != nil {
		return nil, nil, fmt.Errorf("get products by ids: %w", err)
	}

	output := make([]order.OrderedProduct, 0, len(cartProducts))

	for _, prod := range cartProducts {
		productInfo, ok := productsInfo[prod.ProductID]
		if !ok {
			return nil, nil, fmt.Errorf("missing product id: %d", prod.ProductID.Int())
		}

		output = append(output, order.OrderedProduct{
			ProductID:  prod.ProductID,
			CategoryID: prod.CategoryID,
			Count:      prod.Count,
			Price:      productInfo.Price,
		})
	}

	return output, productsInfo, nil
}
