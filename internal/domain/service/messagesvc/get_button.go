package messagesvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/button"
)

func (s *Service) GetButton(ctx context.Context, id button.ID) (*button.Button, error) {
	btn, err := s.buttonRepository.GetButton(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get button from repository: %w", err)
	}

	return btn, nil
}
