package order

import "github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/button"

type CancelOrderPayload struct {
	OrderID       ID
	IsTextMessage bool
}

func CancelOrder(caption string, orderID ID, isTextMsg bool) (button.Button, error) {
	//nolint:wrapcheck
	return button.CreateButton(
		caption,
		button.OperationOrderCancel,
		button.WithPayload(
			CancelOrderPayload{
				OrderID:       orderID,
				IsTextMessage: isTextMsg,
			},
		),
	)
}
