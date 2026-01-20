package cartprocessing

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/internal/message"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/cart"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/perror"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/product"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/infra/logger"
)

func (c *CartProcessing) Confirm(
	ctx context.Context,
	info msginfo.Info,
	cartID cart.ID,
	currencyID currency.ID,
) error {
	storeInfo, err := c.storeInfoByID(ctx, c.storeID)
	if err != nil {
		return fmt.Errorf("check for active: %w", err)
	}

	if !storeInfo.IsActive {
		c.sendPlainText(ctx, info.ChatID, storeInfo.ClosedStoreMessage)

		return nil
	}

	orderedProducts, productsInfo, err := c.orderedProductsFromCart(ctx, cartID, currencyID)
	if err != nil {
		if perror.IsType(err, perror.TypeNotFound) {
			c.editPlainText(ctx, info.ChatID, info.MessageID, message.CartOrderUnavailable())

			return nil
		}

		return fmt.Errorf("ordered products from cart: %w", err)
	}

	if len(orderedProducts) == 0 {
		c.sendPlainText(ctx, info.ChatID, message.NoProductsForOrder())

		return nil
	}

	createdOrder, err := c.repository.CreateOrder(ctx, c.makeCreateOrderInput(info.ChatID, orderedProducts, currencyID))
	if err != nil {
		if c.repository.IsAlreadyExistsError(err) {
			c.sendPlainText(ctx, info.ChatID, message.AlreadyHasActiveOrder())

			return nil
		}

		return fmt.Errorf("repository create order: %w", err)
	}

	if err := c.cart.Clear(ctx, info.ChatID, cartID); err != nil {
		return fmt.Errorf("clear cart: %w", err)
	}

	if err := c.sendOrderInvoice(ctx, info.ChatID, currencyID, createdOrder, productsInfo); err != nil {
		return fmt.Errorf("send order invoice: %w", err)
	}

	c.deleteMessage(ctx, info.ChatID, info.MessageID)

	return nil
}

func (c *CartProcessing) deleteMessage(
	ctx context.Context,
	chatID msginfo.ChatID,
	messageID msginfo.MessageID,
) {
	if err := c.sender.DeleteMessage(ctx, chatID, messageID); err != nil {
		logger.FromContext(ctx).WithError(err).Error("delete message")
	}
}

func (c *CartProcessing) sendOrderInvoice(
	ctx context.Context,
	chatID msginfo.ChatID,
	currencyID currency.ID,
	createdOrder *order.Order,
	productsInfo map[product.ProductID]product.Product,
) error {
	curr, err := c.repository.GetCurrencyByID(ctx, currencyID)
	if err != nil {
		return fmt.Errorf("get currency by id: %w", err)
	}

	if err := c.sender.SendInvoice(
		ctx,
		chatID,
		message.OrderInvoice(),
		createdOrder,
		productsInfo,
		curr,
	); err != nil {
		return fmt.Errorf("send order invoice: %w", err)
	}

	return nil
}

func (c *CartProcessing) makeCreateOrderInput(
	chatID msginfo.ChatID,
	orderedProducts []order.OrderedProduct,
	currencyID currency.ID,
) CreateOrderInput {
	totalPrice := 0
	for _, v := range orderedProducts {
		totalPrice += v.Count * v.Price
	}

	return CreateOrderInput{
		ChatID:              chatID,
		Status:              order.StatusWaitingPayment,
		StatusOperationTime: c.timeProvider.Now(),
		VerificationCode:    "",
		TotalPrice:          totalPrice,
		Products:            orderedProducts,
		CurrencyID:          currencyID,
	}
}

func (c *CartProcessing) orderedProductsFromCart(
	ctx context.Context,
	cartID cart.ID,
	currencyID currency.ID,
) ([]order.OrderedProduct, map[product.ProductID]product.Product, error) {
	cartProducts, err := c.cart.GetProducts(ctx, cartID)
	if err != nil {
		if c.cart.IsNotFoundError(err) {
			return nil, nil, perror.NotFound("cart not found")
		}

		return nil, nil, fmt.Errorf("get cart products: %w", err)
	}

	if len(cartProducts) == 0 {
		return nil, nil, nil
	}

	productIDs := make([]product.ProductID, 0, len(cartProducts))
	for _, v := range cartProducts {
		productIDs = append(productIDs, v.ProductID)
	}

	productsInfo, err := c.repository.GetProductsByIDs(ctx, productIDs, currencyID)
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
