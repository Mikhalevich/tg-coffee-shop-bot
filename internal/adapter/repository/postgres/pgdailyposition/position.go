package pgdailyposition

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

func (p *PgDailyPosition) Position(ctx context.Context, operationTime time.Time) (int, error) {
	var (
		query = `
			INSERT INTO daily_positions(
				day,
				position
			) VALUES (
				:day,
				1
			)
			ON CONFLICT (day) DO UPDATE SET
				position = daily_positions.position + 1
			RETURNING
				position
		`

		trx = p.transactor.ExtContext(ctx)
	)

	query, args, err := sqlx.Named(query, map[string]any{
		"day": makeDay(operationTime),
	})

	if err != nil {
		return 0, fmt.Errorf("sqlx named: %w", err)
	}

	var position int
	if err := sqlx.GetContext(
		ctx,
		trx,
		&position,
		trx.Rebind(query),
		args...,
	); err != nil {
		return 0, fmt.Errorf("get context: %w", err)
	}

	return position, nil
}

// makeDay returns the date of t in its own location, so the day boundary doesn't depend on the db timezone.
func makeDay(t time.Time) string {
	return t.Format(time.DateOnly)
}
