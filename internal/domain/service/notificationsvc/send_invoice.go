package notificationsvc

import (
	"context"
	"fmt"
	"strings"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/internal/message"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/button"
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
	btns, err := makeInvoiceButtons(ord, curr)
	if err != nil {
		return fmt.Errorf("make buttons: %w", err)
	}

	if err := s.sender.SendInvoice(
		ctx,
		msginfo.Invoice{
			ChatID:      chatID,
			Title:       message.OrderInvoice(),
			Description: makeOrderDescription(ord.Products, productsInfo),
			Currency:    curr.Code,
			Payload:     ord.ID.String(),
			Labels:      makeLabeledPrices(ord.Products, productsInfo),
			Buttons:     btns,
		},
	); err != nil {
		return fmt.Errorf("send invoice: %w", err)
	}

	return nil
}

func makeInvoiceButtons(
	ord order.Order,
	curr currency.Currency,
) ([]button.ButtonRow, error) {
	payBtn := button.Pay(
		fmt.Sprintf(
			"%s, %s",
			message.Pay(),
			curr.FormatPrice(ord.TotalPrice),
		),
	)

	cancelBtn, err := order.CancelOrder(message.Cancel(), ord.ID, false)
	if err != nil {
		return nil, fmt.Errorf("cancel order button: %w", err)
	}

	return []button.ButtonRow{
		button.Row(payBtn),
		button.Row(cancelBtn),
	}, nil
}

func makeOrderDescription(
	orderedProducts []order.OrderedProduct,
	productsInfo map[product.ProductID]product.Product,
) string {
	positions := make([]string, 0, len(orderedProducts))

	for _, v := range orderedProducts {
		positions = append(positions,
			fmt.Sprintf("%s x%d", productsInfo[v.ProductID].Title, v.Count))
	}

	return strings.Join(positions, ", ")
}

func makeLabeledPrices(
	orderedProducts []order.OrderedProduct,
	productsInfo map[product.ProductID]product.Product,
) []msginfo.LabeledPrice {
	prices := make([]msginfo.LabeledPrice, 0, len(orderedProducts))

	for _, v := range orderedProducts {
		prices = append(prices, msginfo.LabeledPrice{
			Label:  fmt.Sprintf("%s x%d", productsInfo[v.ProductID].Title, v.Count),
			Amount: v.Count * v.Price,
		})
	}

	return prices
}
