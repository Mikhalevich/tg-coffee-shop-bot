package pgcurrency

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/pgcurrency/internal/model"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/perror"
)

func (p *PgCurrency) GetCurrencyByID(
	ctx context.Context,
	currencyID currency.ID,
) (*currency.Currency, error) {
	var (
		query = `
			SELECT
				id,
				code,
				exp,
				decimal_sep,
				min_amount,
				max_amount,
				is_enabled
			FROM
				currency
			WHERE
				id = :id
		`

		trx = p.transactor.ExtContext(ctx)
	)

	query, args, err := sqlx.Named(query, map[string]any{
		"id": currencyID.Int(),
	})

	if err != nil {
		return nil, fmt.Errorf("sqlx named: %w", err)
	}

	var curr model.Currency
	if err := sqlx.GetContext(
		ctx,
		trx,
		&curr,
		trx.Rebind(query),
		args...,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, perror.NotFound("currency not found")
		}

		return nil, fmt.Errorf("get context: %w", err)
	}

	return curr.ToDom(), nil
}
