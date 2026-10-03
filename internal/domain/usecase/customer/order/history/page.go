package history

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
)

func (h *History) Page(
	ctx context.Context,
	info msginfo.Info,
	page int,
) error {
	if err := h.loadPage(ctx, info, page); err != nil {
		return fmt.Errorf("load page %d: %w", page, err)
	}

	return nil
}
