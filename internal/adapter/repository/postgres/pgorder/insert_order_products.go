package pgorder

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/pgorder/internal/model"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/perror"
)

func (p *PgOrder) InsertProductsOrder(
	ctx context.Context,
	orderID order.ID,
	products []order.OrderedProduct,
) error {
	var (
		query = `
			INSERT INTO order_products(
				order_id,
				product_id,
				count,
				price
			) VALUES (
				:order_id,
				:product_id,
				:count,
				:price
			)
		`
	)

	res, err := sqlx.NamedExecContext(
		ctx,
		p.transactor.ExtContext(ctx),
		query,
		model.ToDBOrderProducts(orderID, products),
	)

	if err != nil {
		return fmt.Errorf("insert order products: %w", err)
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
