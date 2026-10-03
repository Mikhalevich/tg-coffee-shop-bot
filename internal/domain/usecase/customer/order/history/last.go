package history

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
)

func (h *History) Last(
	ctx context.Context,
	info msginfo.Info,
) error {
	if err := h.loadPage(ctx, info, -1); err != nil {
		return fmt.Errorf("load last page: %w", err)
	}

	return nil
}
