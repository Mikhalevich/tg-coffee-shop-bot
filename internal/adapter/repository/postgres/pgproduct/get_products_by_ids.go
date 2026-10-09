package pgproduct

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/pgproduct/internal/model"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/perror"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/product"
)

// GetProductsByIDs returns the products with the given ids keyed by product ID, each with its price
// in currencyID. Products are returned regardless of their enabled state; ids that do not exist or
// have no price in currencyID are omitted. It returns a nil map and nil error if nothing matches,
// and an error if ids is empty.
func (p *PgProduct) GetProductsByIDs(
	ctx context.Context,
	ids []product.ProductID,
	currencyID currency.ID,
) (map[product.ProductID]product.Product, error) {
	if len(ids) == 0 {
		return nil, perror.InvalidParam("ids is empty")
	}

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
		output[product.ProductIDFromInt(v.ID)] = v.ToDom()
	}

	return output, nil
}
