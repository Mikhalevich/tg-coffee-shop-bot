package notificationsvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/internal/message"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/cart"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/product"
)

func (s *Service) ViewOrderCategories(
	ctx context.Context,
	chatID msginfo.ChatID,
	messageID msginfo.MessageID,
	cartID cart.ID,
	categories []product.Category,
	orderedProducts []order.OrderedProduct,
	curr currency.Currency,
) error {
	btns, err := s.makeCartCategoriesButtons(
		cartID,
		categories,
		orderedProducts,
		curr,
	)

	if err != nil {
		return fmt.Errorf("make buttons: %w", err)
	}

	if err := s.sender.SendMessage(
		ctx,
		msginfo.Message{
			ChatID:     chatID,
			ReplyMsgID: messageID,
			Type:       msginfo.MessageTypePlain,
			Text:       "Select category to view products to order",
			Buttons:    btns,
		},
	); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func (s *Service) makeCartCategoriesButtons(
	cartID cart.ID,
	categories []product.Category,
	orderedProducts []order.OrderedProduct,
	curr currency.Currency,
) ([]button.ButtonRow, error) {
	buttons := make([]button.ButtonRow, 0, len(categories)+1)

	for _, v := range categories {
		title := makeViewCategoryButtonTitle(v, orderedProducts, curr)

		b, err := cart.CartViewCategoryProducts(title, cartID, v.ID, curr.ID)
		if err != nil {
			return nil, fmt.Errorf("create cart view button: %w", err)
		}

		buttons = append(buttons, button.Row(b))
	}

	cancelCartBtn, err := cart.CartCancel(message.Cancel(), cartID)
	if err != nil {
		return nil, fmt.Errorf("cancel cart button: %w", err)
	}

	confirmCartBtn, err := cart.CartConfirm(
		makePriceButtonTitle(orderedProducts, curr),
		cartID,
		curr.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("confirm cart button: %w", err)
	}

	buttons = append(buttons, []button.Button{
		cancelCartBtn,
		confirmCartBtn,
	})

	return buttons, nil
}

func makeViewCategoryButtonTitle(
	category product.Category,
	orderedProducts []order.OrderedProduct,
	curr currency.Currency,
) string {
	var (
		count int
		price int
	)

	for _, v := range orderedProducts {
		if category.ID == v.CategoryID {
			price += v.Price * v.Count
			count += v.Count
		}
	}

	if count > 0 {
		return fmt.Sprintf("%s [x%d %s]", category.Title, count, curr.FormatPrice(price))
	}

	return category.Title
}

func makePriceButtonTitle(
	orderedProducts []order.OrderedProduct,
	curr currency.Currency,
) string {
	price := 0
	for _, v := range orderedProducts {
		price += v.Price * v.Count
	}

	if price > 0 {
		return fmt.Sprintf("%s [%s]", message.Confirm(), curr.FormatPrice(price))
	}

	return message.Confirm()
}
