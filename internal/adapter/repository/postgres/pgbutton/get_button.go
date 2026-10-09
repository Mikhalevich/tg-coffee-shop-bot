package pgbutton

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/pgbutton/internal/model"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/perror"
)

func (p *PgButton) GetButton(ctx context.Context, buttonID button.ID) (*button.Button, error) {
	var (
		query = `
			SELECT
				id,
				caption,
				operation,
				is_delete_after_process,
				style,
				url,
				payload,
				pay
			FROM
				buttons
			WHERE
				id = :id
		`

		trx = p.transactor.ExtContext(ctx)
	)

	query, args, err := sqlx.Named(query, map[string]any{
		"id": buttonID.String(),
	})

	if err != nil {
		return nil, fmt.Errorf("sqlx named: %w", err)
	}

	var dbButton model.Button
	if err := sqlx.GetContext(
		ctx,
		trx,
		&dbButton,
		trx.Rebind(query),
		args...,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, perror.NotFound("button not found")
		}

		return nil, fmt.Errorf("get context: %w", err)
	}

	btn := dbButton.ToDom()

	return &btn, nil
}
