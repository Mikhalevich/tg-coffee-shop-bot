package cartprovider

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/cart"
)

func (c *CartProvider) GetProducts(ctx context.Context, id cart.ID) ([]cart.CartProduct, error) {
	items, err := c.client.LRange(ctx, makeCartProductsKey(id.String()), 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("lrange: %w", err)
	}

	if len(items) == 0 {
		return nil, redis.Nil
	}

	cartItems, err := convertToCartItems(combineDuplicateItems(items))
	if err != nil {
		return nil, fmt.Errorf("convert to cart items: %w", err)
	}

	return cartItems, nil
}

func combineDuplicateItems(items []string) map[string]int {
	itemsMap := make(map[string]int, len(items))

	for _, item := range items {
		if _, ok := itemsMap[item]; !ok {
			itemsMap[item] = 1

			continue
		}

		itemsMap[item]++
	}

	return itemsMap
}

func convertToCartItems(itemsMap map[string]int) ([]cart.CartProduct, error) {
	cartItems := make([]cart.CartProduct, 0, len(itemsMap))

	for encodedProduct, count := range itemsMap {
		if isEmptyListPlaceholder(encodedProduct) {
			continue
		}

		decodedProduct, err := decodeCartProduct(encodedProduct)
		if err != nil {
			return nil, fmt.Errorf("decode cart product: %w", err)
		}

		cartItems = append(cartItems, cart.CartProduct{
			ProductID:  decodedProduct.ProductID,
			CategoryID: decodedProduct.CategoryID,
			Count:      count,
		})
	}

	return cartItems, nil
}
