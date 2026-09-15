package notificationsvc

import (
	"context"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/product"
)

func (s *Service) SendInvoice(
	ctx context.Context,
	chatID msginfo.ChatID,
	ord order.Order,
	productsInfo map[product.ProductID]product.Product,
	curr currency.Currency,
) error {
	return nil
}
