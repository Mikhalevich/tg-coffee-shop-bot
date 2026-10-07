package notificationsvc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/internal/message"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
)

func (s *Service) ShowHistoryPage(
	ctx context.Context,
	chatID msginfo.ChatID,
	messageID msginfo.MessageID,
	orders []order.HistoryOrder,
	pageInfo order.Page,
	curr currency.Currency,
) error {
	buttons, err := makeHistoryButtons(pageInfo)
	if err != nil {
		return fmt.Errorf("make history buttons: %w", err)
	}

	if err := s.sender.SendMessage(
		ctx,
		msginfo.Message{
			ChatID:     chatID,
			ReplyMsgID: messageID,
			Type:       msginfo.MessageTypePlain,
			Text:       formatHistoryOrders(orders, curr, pageInfo),
			Buttons:    buttons,
		},
	); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func formatHistoryOrders(
	orders []order.HistoryOrder,
	curr currency.Currency,
	currentPage order.Page,
) string {
	formattedOrders := make([]string, 0, len(orders)+1)
	formattedOrders = append(formattedOrders, fmt.Sprintf("Page %d/%d", currentPage.Number, currentPage.Total))

	for _, ord := range orders {
		formattedOrders = append(formattedOrders,
			fmt.Sprintf("%d) %s, %s, %s",
				ord.SerialNumber,
				ord.CreatedAt.Format(time.RFC3339),
				ord.Status.HumanReadable(),
				curr.FormatPrice(ord.TotalPrice),
			),
		)
	}

	return strings.Join(formattedOrders, "\n")
}

func makeHistoryButtons(
	currentPage order.Page,
) ([]button.ButtonRow, error) {
	var buttons button.ButtonRow

	if currentPage.HasPrevious() {
		nextBtn, err := order.OrderHistoryByPage(message.HistoryNext(), currentPage.Previous())
		if err != nil {
			return nil, fmt.Errorf("next history button: %w", err)
		}

		firstBtn := order.OrderHistoryByPageFirst(message.HistoryFirst())

		buttons = append(buttons, firstBtn, nextBtn)
	}

	if currentPage.HasNext() {
		previousBtn, err := order.OrderHistoryByPage(message.HistoryPrevious(), currentPage.Next())
		if err != nil {
			return nil, fmt.Errorf("previous history button: %w", err)
		}

		lastBtn := order.OrderHistoryByPageLast(message.HistoryLast())

		buttons = append(buttons, previousBtn, lastBtn)
	}

	return []button.ButtonRow{buttons}, nil
}
