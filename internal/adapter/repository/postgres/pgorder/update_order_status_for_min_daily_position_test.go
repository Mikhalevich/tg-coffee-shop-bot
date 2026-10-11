package pgorder_test

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/perror"
)

func (s *PgOrderSuit) TestUpdateOrderStatusForMinDailyPosition() {
	s.Run("picks min daily position instead of min id", func() {
		var (
			ctx        = s.T().Context()
			day        = time.Date(2026, time.October, 10, 9, 0, 0, 0, time.UTC)
			now        = day.Add(time.Hour)
			currencyID = s.insertCurrency(ctx)
		)

		s.insertOrder(ctx, 1, order.StatusConfirmed, currencyID, 3, day)
		expectedID := s.insertOrder(ctx, 2, order.StatusConfirmed, currencyID, 1, day)
		s.insertOrder(ctx, 3, order.StatusConfirmed, currencyID, 2, day)

		actualID, err := s.pgOrder.UpdateOrderStatusForMinDailyPosition(
			ctx, now, order.StatusInProgress, order.StatusConfirmed)

		s.Require().NoError(err)
		s.Require().Equal(expectedID, actualID)
		s.requireOrderStatus(ctx, expectedID, order.StatusInProgress, now)
	})

	s.Run("earlier day goes first", func() {
		var (
			ctx        = s.T().Context()
			prevDay    = time.Date(2026, time.October, 9, 22, 0, 0, 0, time.UTC)
			nextDay    = prevDay.AddDate(0, 0, 1)
			now        = nextDay.Add(time.Hour)
			currencyID = s.insertCurrency(ctx)
		)

		s.insertOrder(ctx, 1, order.StatusConfirmed, currencyID, 1, nextDay)
		expectedID := s.insertOrder(ctx, 2, order.StatusConfirmed, currencyID, 5, prevDay)

		actualID, err := s.pgOrder.UpdateOrderStatusForMinDailyPosition(
			ctx, now, order.StatusInProgress, order.StatusConfirmed)

		s.Require().NoError(err)
		s.Require().Equal(expectedID, actualID)
		s.requireOrderStatus(ctx, expectedID, order.StatusInProgress, now)
	})

	s.Run("ignores orders in other statuses", func() {
		var (
			ctx        = s.T().Context()
			day        = time.Date(2026, time.October, 10, 9, 0, 0, 0, time.UTC)
			now        = day.Add(time.Hour)
			currencyID = s.insertCurrency(ctx)
		)

		inProgressID := s.insertOrder(ctx, 1, order.StatusInProgress, currencyID, 1, day)
		expectedID := s.insertOrder(ctx, 2, order.StatusConfirmed, currencyID, 2, day)

		actualID, err := s.pgOrder.UpdateOrderStatusForMinDailyPosition(
			ctx, now, order.StatusInProgress, order.StatusConfirmed)

		s.Require().NoError(err)
		s.Require().Equal(expectedID, actualID)
		s.requireOrderStatus(ctx, inProgressID, order.StatusInProgress, day)
	})

	s.Run("not found", func() {
		var (
			ctx        = s.T().Context()
			day        = time.Date(2026, time.October, 10, 9, 0, 0, 0, time.UTC)
			currencyID = s.insertCurrency(ctx)
		)

		s.insertOrder(ctx, 1, order.StatusInProgress, currencyID, 1, day)

		actualID, err := s.pgOrder.UpdateOrderStatusForMinDailyPosition(
			ctx, day, order.StatusInProgress, order.StatusConfirmed)

		s.Require().Error(err)
		s.Require().True(perror.IsType(err, perror.TypeNotFound))
		s.Require().Equal(order.ID(0), actualID)
	})
}

func (s *PgOrderSuit) insertCurrency(ctx context.Context) int {
	var currencyID int

	s.Require().NoError(sqlx.GetContext(ctx, s.transactor.ExtContext(ctx), &currencyID, `
		INSERT INTO currency (
			code,
			exp,
			decimal_sep,
			min_amount,
			max_amount,
			is_enabled
		) VALUES (
			'USD', 2, '.', 100, 100000, true
		)
		RETURNING id
	`))

	return currencyID
}

func (s *PgOrderSuit) insertOrder(
	ctx context.Context,
	chatID int64,
	status order.Status,
	currencyID int,
	dailyPosition int,
	updatedAt time.Time,
) order.ID {
	trx := s.transactor.ExtContext(ctx)

	query, args, err := sqlx.Named(`
		INSERT INTO orders (
			chat_id,
			status,
			currency_id,
			daily_position,
			total_price,
			created_at,
			updated_at
		) VALUES (
			:chat_id,
			:status,
			:currency_id,
			:daily_position,
			100,
			:updated_at,
			:updated_at
		)
		RETURNING id
	`, map[string]any{
		"chat_id":        chatID,
		"status":         status,
		"currency_id":    currencyID,
		"daily_position": dailyPosition,
		"updated_at":     updatedAt,
	})
	s.Require().NoError(err)

	var orderID int

	s.Require().NoError(sqlx.GetContext(ctx, trx, &orderID, trx.Rebind(query), args...))

	return order.IDFromInt(orderID)
}

func (s *PgOrderSuit) requireOrderStatus(
	ctx context.Context,
	orderID order.ID,
	expectedStatus order.Status,
	expectedUpdatedAt time.Time,
) {
	var row struct {
		Status    order.Status `db:"status"`
		UpdatedAt time.Time    `db:"updated_at"`
	}

	s.Require().NoError(sqlx.GetContext(ctx, s.transactor.ExtContext(ctx), &row,
		"SELECT status, updated_at FROM orders WHERE id = $1", orderID.Int()))

	s.Require().Equal(expectedStatus, row.Status)
	s.Require().True(expectedUpdatedAt.Equal(row.UpdatedAt))
}
