package orderpayment

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/product"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/infra/logger"
)

func (p *OrderPayment) PaymentConfirmed(
	ctx context.Context,
	chatID msginfo.ChatID,
	orderID order.ID,
	currency string,
	totalAmount int,
) error {
	now := p.timeProvider.Now()

	position, err := p.dailyPosition.Position(ctx, now)
	if err != nil {
		return fmt.Errorf("daily position: %w", err)
	}

	if err := p.transactor.Transaction(ctx, func(ctx context.Context) error {
		ord, err := p.repository.UpdateOrderByChatAndID(
			ctx,
			orderID,
			chatID,
			port.UpdateOrderData{
				Status:              order.StatusConfirmed,
				StatusOperationTime: now,
				VerificationCode:    p.codeGenerator.Generate(),
				DailyPosition:       position,
			},
			order.StatusPaymentInProgress,
		)
		if err != nil {
			return fmt.Errorf("update order status: %w", err)
		}

		queuePosition := p.orderQueuePosition(ctx, ord)

		if err := p.sendOrderQRImage(ctx, chatID, ord, queuePosition); err != nil {
			return fmt.Errorf("send order qr: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}

func (p *OrderPayment) sendOrderQRImage(
	ctx context.Context,
	chatID msginfo.ChatID,
	ord *order.Order,
	queuePosition int,
) error {
	png, err := p.qrCode.GeneratePNG(ord.ID.String())
	if err != nil {
		return fmt.Errorf("qrcode generate png: %w", err)
	}

	productsInfo, err := p.repository.GetProductsByIDs(ctx, ord.ProductIDs(), ord.CurrencyID)
	if err != nil {
		return fmt.Errorf("get products by ids: %w", err)
	}

	curr, err := p.repository.GetCurrencyByID(ctx, ord.CurrencyID)
	if err != nil {
		return fmt.Errorf("get currency by id: %w", err)
	}

	if err := p.sender.SendMessage(
		ctx,
		messageprocessor.Message{
			ChatID:  chatID,
			Text:    formatOrder(ord, curr, productsInfo, queuePosition, p.escaper.EscapeMarkdown),
			Type:    messageprocessor.MessageTypePNG,
			Payload: png,
		},
	); err != nil {
		return fmt.Errorf("send png: %w", err)
	}

	return nil
}

func (p *OrderPayment) orderQueuePosition(ctx context.Context, activeOrder *order.Order) int {
	if !activeOrder.InQueue() {
		return 0
	}

	pos, err := p.repository.GetOrderPositionByStatus(
		ctx,
		activeOrder.ID,
		order.StatusConfirmed,
		order.StatusInProgress,
	)

	if err != nil {
		if p.repository.IsNotFoundError(err) {
			return 0
		}

		logger.FromContext(ctx).WithError(err).Error("failed to get order position")

		return 0
	}

	return pos
}

func formatOrder(
	ord *order.Order,
	curr *currency.Currency,
	productsInfo map[product.ProductID]product.Product,
	queuePosition int,
	escaper func(string) string,
) string {
	format := []string{
		fmt.Sprintf("order id: *%s*", escaper(ord.ID.String())),
		fmt.Sprintf("status: *%s*", ord.Status.HumanReadable()),
		fmt.Sprintf("verification code: *%s*", escaper(ord.VerificationCode)),
		fmt.Sprintf("daily position: *%d*", ord.DailyPosition),
		fmt.Sprintf("total price: *%s*", curr.FormatPrice(ord.TotalPrice)),
		fmt.Sprintf("created\\_at: *%s*", escaper(ord.CreatedAt.Format(time.RFC3339))),
		fmt.Sprintf("updated\\_at: *%s*", escaper(ord.UpdatedAt.Format(time.RFC3339))),
	}

	for _, t := range ord.Timeline {
		format = append(format, fmt.Sprintf(
			"%s Time: *%s*",
			t.Status.HumanReadable(),
			escaper(t.Time.Format(time.RFC3339))),
		)
	}

	for _, v := range ord.Products {
		format = append(format, fmt.Sprintf("%s x%d %s",
			escaper(productsInfo[v.ProductID].Title), v.Count, curr.FormatPrice(v.Price)))
	}

	if queuePosition > 0 {
		format = append(format, fmt.Sprintf("position in queue: *%d*", queuePosition))
	}

	return strings.Join(format, "\n")
}
