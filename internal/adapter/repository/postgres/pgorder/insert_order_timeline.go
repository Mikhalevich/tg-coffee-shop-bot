package pgorder

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/pgorder/internal/model"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/perror"
)

func (p *PgOrder) InsertOrderTimeline(
	ctx context.Context,
	orderID order.ID,
	status order.Status,
	createdAt time.Time,
) error {
	var (
		query = `
			INSERT INTO order_status_timeline(
				order_id,
				status,
				updated_at
			) VALUES(
				:order_id,
				:status,
				:updated_at
			)
		`
	)
	res, err := sqlx.NamedExecContext(
		ctx,
		p.transactor.ExtContext(ctx),
		query,
		model.OrderTimeline{
			ID:        orderID.Int(),
			Status:    status.String(),
			UpdatedAt: createdAt,
		},
	)

	if err != nil {
		return fmt.Errorf("named exec context: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return perror.NoRowsUpdated()
	}

	return nil
}
