package payment

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/infra/logger"
)

func (p *Payment) InProgress(
	ctx context.Context,
	paymentID string,
	orderID order.ID,
	totalAmount int,
) error {
	storeInfo, err := p.storeService.GetStoreInfo(ctx)
	if err != nil {
		return fmt.Errorf("get store info: %w", err)
	}

	if !storeInfo.IsActive {
		if err := p.notificationService.PaymentStoreClosed(
			ctx,
			paymentID,
			storeInfo.CurrentTime,
			storeInfo.NextWorkingTime,
		); err != nil {
			return fmt.Errorf("payment store closed: %w", err)
		}

		return nil
	}

	if err := p.transactor.Transaction(ctx, func(ctx context.Context) error {
		err := p.orderService.SetOrderPaymentInProgress(
			ctx,
			orderID,
			totalAmount,
			storeInfo.CurrentTime,
		)
		if err != nil {
			logger.FromContext(ctx).
				WithError(err).
				Warn("set order in progress")
		}

		if err := p.notificationService.PaymentAnswer(ctx, paymentID, err == nil); err != nil {
			return fmt.Errorf("payment success msg: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}
