package payment

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
)

func (p *Payment) Confirmed(
	ctx context.Context,
	chatID msginfo.ChatID,
	orderID order.ID,
	currency string,
	totalAmount int,
) error {
	var (
		now              = p.timeProvider.Now()
		verificationCode = p.codeGeneratorService.Generate()
	)

	if err := p.transactor.Transaction(ctx, func(ctx context.Context) error {
		position, err := p.positionService.Position(ctx, now)
		if err != nil {
			return fmt.Errorf("daily position: %w", err)
		}

		if err := p.confimOrder(
			ctx,
			chatID,
			orderID,
			verificationCode,
			position,
			now,
		); err != nil {
			return fmt.Errorf("confirm order: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}

func (p *Payment) confimOrder(
	ctx context.Context,
	chatID msginfo.ChatID,
	orderID order.ID,
	verificationCode string,
	dailyPosition int,
	operationTime time.Time,
) error {
	queuePosition, err := p.orderService.SetOrderPaymentConfirmed(
		ctx,
		chatID,
		orderID,
		verificationCode,
		dailyPosition,
		operationTime,
	)
	if err != nil {
		return fmt.Errorf("set payment confirmed: %w", err)
	}

	ord, err := p.orderService.GetOrderByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("get order by id: %w", err)
	}

	curr, err := p.currencyService.GetCurrencyByID(ctx, ord.CurrencyID)
	if err != nil {
		return fmt.Errorf("get currency by id: %w", err)
	}

	png, err := p.qrCodeService.GeneratePNG(orderID.String())
	if err != nil {
		return fmt.Errorf("generate qr code: %w", err)
	}

	productsInfo, err := p.productsService.GetProductsByIDs(
		ctx,
		ord.ProductIDs(),
		ord.CurrencyID,
	)

	if err != nil {
		return fmt.Errorf("get products by ids: %w", err)
	}

	if err := p.notificationService.PaymentConfirmed(
		ctx,
		chatID,
		ord,
		curr,
		productsInfo,
		queuePosition,
		png,
	); err != nil {
		return fmt.Errorf("payments confirmed msg: %w", err)
	}

	return nil
}
