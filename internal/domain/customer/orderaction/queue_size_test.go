package orderaction_test

import (
	"context"
	"errors"

	"go.uber.org/mock/gomock"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
)

func (s *OrderActionSuite) TestQueueSizeGetOrdersCountByStatusError() {
	var (
		ctx  = context.Background()
		info = msginfo.Info{
			ChatID:    msginfo.ChatIDFromInt(207),
			MessageID: msginfo.MessageIDFromInt(100),
		}
	)

	s.repository.EXPECT().
		GetOrdersCountByStatus(ctx, order.StatusConfirmed, order.StatusInProgress).
		Return(0, errors.New("some error"))

	err := s.orderAction.QueueSize(ctx, info)

	s.Require().EqualError(err, "get orders count by status: some error")
}

func (s *OrderActionSuite) TestQueueSizeSuccess() {
	var (
		ctx  = context.Background()
		info = msginfo.Info{
			ChatID:    msginfo.ChatIDFromInt(207),
			MessageID: msginfo.MessageIDFromInt(100),
		}
	)

	gomock.InOrder(
		s.repository.EXPECT().
			GetOrdersCountByStatus(ctx, order.StatusConfirmed, order.StatusInProgress).
			Return(1, nil),

		s.sender.EXPECT().ReplyTextMarkdown(ctx, info.ChatID, info.MessageID, "*1*"),
	)

	err := s.orderAction.QueueSize(ctx, info)

	s.Require().NoError(err)
}
