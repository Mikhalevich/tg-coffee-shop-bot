package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/internal/model"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/customer/cartprocessing"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
)

func (p *Postgres) CreateOrder(ctx context.Context, coi cartprocessing.CreateOrderInput) (*order.Order, error) {
	var orderResult order.Order

	if err := p.transactor.Transaction(ctx, func(ctx context.Context) error {
		trx := p.transactor.ExtContext(ctx)
		orderID, err := p.insertOrder(ctx, trx, model.Order{
			ChatID:           coi.ChatID.Int64(),
			Status:           coi.Status.String(),
			VerificationCode: model.NullString(coi.VerificationCode),
			CurrencyID:       coi.CurrencyID.Int(),
			TotalPrice:       coi.TotalPrice,
			CreatedAt:        coi.StatusOperationTime,
			UpdatedAt:        coi.StatusOperationTime,
		})

		if err != nil {
			return fmt.Errorf("insert order: %w", err)
		}

		if err := insertProductsToOrder(ctx, trx, model.PortToOrderProducts(orderID, coi.Products)); err != nil {
			return fmt.Errorf("insert order products: %w", err)
		}

		if err := insertOrderTimeline(ctx, trx, model.OrderTimeline{
			ID:        orderID.Int(),
			Status:    coi.Status.String(),
			UpdatedAt: coi.StatusOperationTime,
		}); err != nil {
			return fmt.Errorf("insert order timeline: %w", err)
		}

		orderResult = convertToOrder(orderID, coi)

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %w", err)
	}

	return &orderResult, nil
}

func convertToOrder(id order.ID, input cartprocessing.CreateOrderInput) order.Order {
	return order.Order{
		ID:               id,
		ChatID:           input.ChatID,
		Status:           input.Status,
		VerificationCode: input.VerificationCode,
		CurrencyID:       input.CurrencyID,
		TotalPrice:       input.TotalPrice,
		Products:         input.Products,
	}
}

func (p *Postgres) insertOrder(ctx context.Context, ext sqlx.ExtContext, dbOrder model.Order) (order.ID, error) {
	query, args, err := sqlx.Named(`
		INSERT INTO orders(
			chat_id,
			status,
			verification_code,
			currency_id,
			total_price,
			created_at,
			updated_at
		) VALUES (
			:chat_id,
			:status,
			:verification_code,
			:currency_id,
			:total_price,
			:created_at,
			:updated_at
		)
		RETURNING id
	`, dbOrder)

	if err != nil {
		return 0, fmt.Errorf("prepare named: %w", err)
	}

	var orderID int
	if err := sqlx.GetContext(ctx, ext, &orderID, ext.Rebind(query), args...); err != nil {
		if p.driver.IsConstraintError(err, "orders_only_one_active_order_unique_idx") {
			return 0, errAlreadyExists
		}

		return 0, fmt.Errorf("insert order: %w", err)
	}

	return order.IDFromInt(orderID), nil
}

func insertProductsToOrder(ctx context.Context, ext sqlx.ExtContext, products []model.OrderProduct) error {
	res, err := sqlx.NamedExecContext(ctx, ext, `
		INSERT INTO order_products(
			order_id,
			product_id,
			count,
			price
		) VALUES (
			:order_id,
			:product_id,
			:count,
			:price
		)`,
		products,
	)

	if err != nil {
		return fmt.Errorf("insert order products: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return errNotUpdated
	}

	return nil
}
