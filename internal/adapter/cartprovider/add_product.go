package cartprovider

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"strings"

	"github.com/redis/go-redis/v9"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/cart"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/perror"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/product"
)

func (c *CartProvider) AddProduct(ctx context.Context, cartID cart.ID, p cart.CartProduct) error {
	encodedProduct, err := encodeCartProduct(cartProduct{
		ProductID:  p.ProductID,
		CategoryID: p.CategoryID,
	})
	if err != nil {
		return fmt.Errorf("encode cart product: %w", err)
	}

	if err := c.addProductToExistingList(
		ctx,
		makeCartProductsKey(cartID.String()),
		encodedProduct,
	); err != nil {
		return fmt.Errorf("add product to existing list: %w", err)
	}

	return nil
}

type cartProduct struct {
	ProductID  product.ProductID
	CategoryID product.CategoryID
}

func encodeCartProduct(p cartProduct) (string, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(p); err != nil {
		return "", fmt.Errorf("gob encode: %w", err)
	}

	return buf.String(), nil
}

func decodeCartProduct(s string) (cartProduct, error) {
	var p cartProduct
	if err := gob.NewDecoder(strings.NewReader(s)).Decode(&p); err != nil {
		return cartProduct{}, fmt.Errorf("gob decode: %w", err)
	}

	return p, nil
}

// addProductToExistingList appends the product to an existing cart list
// and returns a not-found error if the cart has expired.
func (c *CartProvider) addProductToExistingList(ctx context.Context, key string, data string) error {
	newLen, err := c.client.RPushX(ctx, key, data).Result()
	if err != nil {
		return fmt.Errorf("rpushx: %w", err)
	}

	if newLen == 0 {
		return perror.NotFound("cart not found")
	}

	return nil
}

//nolint:unused
func (c *CartProvider) addProductToNotExistingList(ctx context.Context, key string, id product.ProductID) error {
	if _, err := c.client.Pipelined(ctx, func(pipline redis.Pipeliner) error {
		if err := pipline.RPush(ctx, key, id.String()).Err(); err != nil {
			return fmt.Errorf("rpush: %w", err)
		}

		if err := pipline.Expire(ctx, key, c.ttl).Err(); err != nil {
			return fmt.Errorf("expire: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("piplined: %w", err)
	}

	return nil
}
