package order

import (
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/button"
)

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

type OrderHistoryByPagePayload struct {
	Page int
}

func OrderHistoryByPage(caption string, page int) (button.Button, error) {
	//nolint:wrapcheck
	return button.CreateButton(
		caption,
		button.OperationOrderHistoryByPage,
		button.WithPayload(
			OrderHistoryByPagePayload{
				Page: page,
			},
		),
	)
}

func OrderHistoryByPageFirst(caption string) button.Button {
	return button.MustCreateButton(
		caption,
		button.OperationOrderHistoryByPageFirst,
	)
}

func OrderHistoryByPageLast(caption string) button.Button {
	return button.MustCreateButton(
		caption,
		button.OperationOrderHistoryByPageLast,
	)
}
