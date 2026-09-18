package notificationsvc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/product"
)

func (s *Service) ViewOrder(
	ctx context.Context,
	chatID msginfo.ChatID,
	ord order.Order,
	products map[product.ProductID]product.Product,
	curr currency.Currency,
) error {
	if err := s.sender.SendMessage(
		ctx,
		msginfo.Message{
			ChatID: chatID,
			Type:   msginfo.MessageTypeMarkdown,
			Text: s.formatOrder(
				ord,
				products,
				curr,
				0,
			),
		},
	); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func (s *Service) formatOrder(
	ord order.Order,
	productsInfo map[product.ProductID]product.Product,
	curr currency.Currency,
	queuePosition int,
) string {
	format := []string{
		fmt.Sprintf("order id: *%s*", s.escaper.EscapeMarkdown(ord.ID.String())),
		fmt.Sprintf("status: *%s*", ord.Status.HumanReadable()),
		fmt.Sprintf("verification code: *%s*", s.escaper.EscapeMarkdown(ord.VerificationCode)),
		fmt.Sprintf("daily position: *%d*", ord.DailyPosition),
		fmt.Sprintf("total price: *%s*", curr.FormatPrice(ord.TotalPrice)),
		fmt.Sprintf("created\\_at: *%s*", s.escaper.EscapeMarkdown(ord.CreatedAt.Format(time.RFC3339))),
		fmt.Sprintf("updated\\_at: *%s*", s.escaper.EscapeMarkdown(ord.UpdatedAt.Format(time.RFC3339))),
	}

	for _, t := range ord.Timeline {
		format = append(format, fmt.Sprintf(
			"%s Time: *%s*",
			t.Status.HumanReadable(),
			s.escaper.EscapeMarkdown(t.Time.Format(time.RFC3339))),
		)
	}

	for _, v := range ord.Products {
		format = append(format, fmt.Sprintf("%s x%d %s",
			s.escaper.EscapeMarkdown(productsInfo[v.ProductID].Title), v.Count, curr.FormatPrice(v.Price)))
	}

	if queuePosition > 0 {
		format = append(format, fmt.Sprintf("position in queue: *%d*", queuePosition))
	}

	return strings.Join(format, "\n")
}
