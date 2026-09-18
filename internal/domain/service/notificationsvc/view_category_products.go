package notificationsvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/internal/message"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/cart"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/product"
)

func (s *Service) ViewCategoryProducts(
	ctx context.Context,
	chatID msginfo.ChatID,
	messageID msginfo.MessageID,
	cartID cart.ID,
	categoryID product.CategoryID,
	categoryProducts []product.Product,
	cartProducts []cart.CartProduct,
	curr currency.Currency,
) error {
	btns, err := s.makeCartProductsButtons(
		cartID,
		categoryID,
		categoryProducts,
		cartProducts,
		curr,
	)

	if err != nil {
		return fmt.Errorf("view category products: %w", err)
	}

	if err := s.sender.SendMessage(
		ctx,
		msginfo.Message{
			ChatID:     chatID,
			ReplyMsgID: messageID,
			Type:       msginfo.MessageTypePlain,
			Text:       message.OrderProductPage(),
			Buttons:    btns,
		},
	); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func (s *Service) makeCartProductsButtons(
	cartID cart.ID,
	categoryID product.CategoryID,
	categoryProducts []product.Product,
	cartProducts []cart.CartProduct,
	curr currency.Currency,
) ([]button.ButtonRow, error) {
	buttons := make([]button.ButtonRow, 0, len(categoryProducts)+1)

	for _, v := range categoryProducts {
		title := makeProductButtonTitle(v, cartProducts, curr)
		btn, err := cart.CartAddProduct(title, cartID, v.ID, categoryID, curr.ID)

		if err != nil {
			return nil, fmt.Errorf("add product button: %w", err)
		}

		buttons = append(buttons, button.Row(btn))
	}

	viewCategoriesBtn, err := cart.CartViewCategories(message.Done(), cartID, curr.ID)
	if err != nil {
		return nil, fmt.Errorf("cart view categories button: %w", err)
	}

	buttons = append(buttons, button.Row(viewCategoriesBtn))

	return buttons, nil
}

func makeProductButtonTitle(
	prod product.Product,
	cartProducts []cart.CartProduct,
	curr currency.Currency,
) string {
	for _, cartProduct := range cartProducts {
		if cartProduct.ProductID == prod.ID {
			return fmt.Sprintf("%s %s [x%d %s]",
				prod.Title,
				curr.FormatPrice(prod.Price),
				cartProduct.Count,
				curr.FormatPrice(prod.Price*cartProduct.Count),
			)
		}
	}

	return fmt.Sprintf("%s %s",
		prod.Title,
		curr.FormatPrice(prod.Price),
	)
}
