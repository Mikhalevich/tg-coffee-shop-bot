package pgbutton

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/pgbutton/internal/model"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/button"
)

func (p *PgButton) SetButtonRows(
	ctx context.Context,
	rows ...button.ButtonRow,
) error {
	dbButtons := model.ToDBButtons(rows...)
	if len(dbButtons) == 0 {
		return nil
	}

	var (
		query = `
			INSERT INTO buttons(
				id,
				caption,
				operation,
				is_delete_after_process,
				style,
				url,
				payload,
				pay
			) VALUES (
				:id,
				:caption,
				:operation,
				:is_delete_after_process,
				:style,
				:url,
				:payload,
				:pay
			)
			ON CONFLICT (id) DO UPDATE SET
				caption = EXCLUDED.caption,
				operation = EXCLUDED.operation,
				is_delete_after_process = EXCLUDED.is_delete_after_process,
				style = EXCLUDED.style,
				url = EXCLUDED.url,
				payload = EXCLUDED.payload,
				pay = EXCLUDED.pay
		`
	)

	if _, err := sqlx.NamedExecContext(
		ctx,
		p.transactor.ExtContext(ctx),
		query,
		dbButtons,
	); err != nil {
		return fmt.Errorf("named exec: %w", err)
	}

	return nil
}
