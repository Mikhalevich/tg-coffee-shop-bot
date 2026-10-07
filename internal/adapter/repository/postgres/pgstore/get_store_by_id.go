package pgstore

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/pgstore/internal/model"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/perror"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/store"
)

func (p *PgStore) GetStoreByID(
	ctx context.Context,
	storeID store.ID,
) (*store.Store, error) {
	trx := p.transactor.ExtContext(ctx)
	modelStore, err := selectStoreByID(ctx, trx, storeID)
	if err != nil {
		return nil, fmt.Errorf("select store by id: %w", err)
	}

	modelSchedule, err := selectStoreSchedule(ctx, trx, storeID)
	if err != nil {
		return nil, fmt.Errorf("select store schedule: %w", err)
	}

	portStore, err := modelStore.ToDom(modelSchedule)
	if err != nil {
		return nil, fmt.Errorf("convert to port store: %w", err)
	}

	return portStore, nil
}

func selectStoreByID(ctx context.Context, trx sqlx.ExtContext, storeID store.ID) (*model.Store, error) {
	query, args, err := sqlx.Named(`
		SELECT
			id,
			description,
			default_currency_id
		FROM
			store
		WHERE
			id = :id
	`, map[string]any{
		"id": storeID.Int(),
	})

	if err != nil {
		return nil, fmt.Errorf("named: %w", err)
	}

	var store model.Store
	if err := sqlx.GetContext(ctx, trx, &store, trx.Rebind(query), args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, perror.NotFound("store not found")
		}

		return nil, fmt.Errorf("get context: %w", err)
	}

	return &store, nil
}

func selectStoreSchedule(ctx context.Context, trx sqlx.ExtContext, storeID store.ID) ([]model.StoreSchedule, error) {
	query, args, err := sqlx.Named(`
		SELECT
			store_id,
			day_of_week,
			start_time,
			end_time
		FROM
			store_schedule
		WHERE
			store_id = :store_id
	`, map[string]any{
		"store_id": storeID.Int(),
	})

	if err != nil {
		return nil, fmt.Errorf("named: %w", err)
	}

	var schedule []model.StoreSchedule
	if err := sqlx.SelectContext(ctx, trx, &schedule, trx.Rebind(query), args...); err != nil {
		return nil, fmt.Errorf("select context: %w", err)
	}

	return schedule, nil
}
