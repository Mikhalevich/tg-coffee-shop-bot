package messageprocessor

import (
	"context"
	"fmt"
	"strings"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/internal/message"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor/button"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/product"
)

type LabeledPrice struct {
	Label  string
	Amount int
}

func (m *MessageProcessor) SendInvoice(
	ctx context.Context,
	chatID msginfo.ChatID,
	title string,
	ord *order.Order,
	productsInfo map[product.ProductID]product.Product,
	curr *currency.Currency,
) error {
	rows, err := makeInvoiceButtons(chatID, ord, curr)
	if err != nil {
		return fmt.Errorf("make invoice buttons: %w", err)
	}

	inlineButtons, err := m.SetButtonRows(ctx, rows...)
	if err != nil {
		return fmt.Errorf("set button rows: %w", err)
	}

	if err := m.sender.SendOrderInvoice(
		ctx,
		chatID,
		title,
		makeOrderDescription(ord.Products, productsInfo),
		curr.Code,
		ord.ID.String(),
		makeLabeledPrices(ord.Products, productsInfo),
		inlineButtons...,
	); err != nil {
		return fmt.Errorf("send order invoice: %w", err)
	}

	return nil
}

func makeInvoiceButtons(
	chatID msginfo.ChatID,
	ord *order.Order,
	curr *currency.Currency,
) ([]button.ButtonRow, error) {
	payBtn := button.Pay(fmt.Sprintf("%s, %s", message.Pay(), curr.FormatPrice(ord.TotalPrice)))

	cancelBtn, err := button.CancelOrder(chatID, message.Cancel(), ord.ID, false)
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
		positions = append(positions, fmt.Sprintf("%s x%d", productsInfo[v.ProductID].Title, v.Count))
	}

	return strings.Join(positions, ", ")
}

func makeLabeledPrices(
	orderedProducts []order.OrderedProduct,
	productsInfo map[product.ProductID]product.Product,
) []LabeledPrice {
	prices := make([]LabeledPrice, 0, len(orderedProducts))

	for _, v := range orderedProducts {
		prices = append(prices, LabeledPrice{
			Label:  fmt.Sprintf("%s x%d", productsInfo[v.ProductID].Title, v.Count),
			Amount: v.Count * v.Price,
		})
	}

	return prices
}
