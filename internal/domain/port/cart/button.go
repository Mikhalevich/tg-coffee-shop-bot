package cart

import (
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/product"
)

type CartCancelPayload struct {
	CartID ID
}

func CartCancel(caption string, cartID ID) (button.Button, error) {
	//nolint:wrapcheck
	return button.CreateButton(
		caption,
		button.OperationCartCancel,
		button.WithPayload(
			CartCancelPayload{
				CartID: cartID,
			},
		),
	)
}

type CartConfirmPayload struct {
	CartID     ID
	CurrencyID currency.ID
}

func CartConfirm(
	caption string,
	cartID ID,
	currencyID currency.ID,
) (button.Button, error) {
	//nolint:wrapcheck
	return button.CreateButton(
		caption,
		button.OperationCartConfirm,
		button.WithPayload(
			CartConfirmPayload{
				CartID:     cartID,
				CurrencyID: currencyID,
			},
		),
	)
}

type CartViewCategoryProductsPayload struct {
	CartID     ID
	CategoryID product.CategoryID
	CurrencyID currency.ID
}

func CartViewCategoryProducts(
	caption string,
	cartID ID,
	categoryID product.CategoryID,
	currencyID currency.ID,
) (button.Button, error) {
	//nolint:wrapcheck
	return button.CreateButton(
		caption,
		button.OperationCartViewCategoryProducts,
		button.WithPayload(
			CartViewCategoryProductsPayload{
				CartID:     cartID,
				CategoryID: categoryID,
				CurrencyID: currencyID,
			},
		),
	)
}
