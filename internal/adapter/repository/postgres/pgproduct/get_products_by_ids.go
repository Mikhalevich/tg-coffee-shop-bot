package pgproduct

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/internal/model"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/product"
)

func (p *PgProduct) GetProductsByIDs(
	ctx context.Context,
	ids []product.ProductID,
	currencyID currency.ID,
) (map[product.ProductID]product.Product, error) {
	var (
		query = `
			SELECT
				p.id,
				p.title,
				pp.currency_id,
				pp.price,
				p.is_enabled,
				p.created_at,
				p.updated_at
			FROM
				product p INNER JOIN product_price pp ON p.id = pp.product_id
			WHERE
				p.id IN(?) AND
				pp.currency_id = ?
		`
	)

	query, args, err := sqlx.In(query, ids, currencyID.Int())

	if err != nil {
		return nil, fmt.Errorf("sqlx in: %w", err)
	}

	var (
		trx        = p.transactor.ExtContext(ctx)
		dbProducts []model.Product
	)

	if err := sqlx.SelectContext(
		ctx,
		trx,
		&dbProducts,
		trx.Rebind(query),
		args...,
	); err != nil {
		return nil, fmt.Errorf("select context: %w", err)
	}

	if len(dbProducts) == 0 {
		//nolint:nilnil
		return nil, nil
	}

	output := make(map[product.ProductID]product.Product, len(dbProducts))

	for _, v := range dbProducts {
		output[product.ProductIDFromInt(v.ID)] = v.ToPortProduct()
	}

	return output, nil
}
