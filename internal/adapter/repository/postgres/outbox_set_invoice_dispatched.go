package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

func (p *Postgres) OutboxSetInvoiceDispatched(
	ctx context.Context,
	ids []int,
	dispatchedAt time.Time,
) error {
	if len(ids) == 0 {
		return nil
	}

	var (
		query = `
			UPDATE outbox_order_invoice
			SET
				is_dispatched = TRUE,
				dispatched_at = ?
			WHERE
				id IN(?)
		`
	)

	query, args, err := sqlx.In(query, dispatchedAt, ids)
	if err != nil {
		return fmt.Errorf("sqlx in: %w", err)
	}

	ext := p.transactor.ExtContext(ctx)

	res, err := ext.ExecContext(ctx, ext.Rebind(query), args...)
	if err != nil {
		return fmt.Errorf("exec context: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if affected == 0 {
		return errors.New("no rows affected")
	}

	return nil
}
