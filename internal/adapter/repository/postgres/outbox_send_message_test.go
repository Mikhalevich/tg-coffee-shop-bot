package postgres_test

import (
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/internal/jsonb"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/internal/model"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
)

func (s *PostgresSuit) TestSendMessageWitoutPayloadAndButtons() {
	s.Run("store plain text", func() {
		inputMsg := messageprocessor.Message{
			ChatID:     msginfo.ChatIDFromInt(1),
			ReplyMsgID: msginfo.MessageIDFromInt(1),
			Text:       "test text",
			Type:       messageprocessor.MessageTypePlain,
		}

		err := s.pgDB.SendMessage(s.T().Context(), inputMsg)

		s.Require().NoError(err)

		actualMsg, err := s.pgDB.TestGetOutboxMessageByChatID(s.T().Context(), int(inputMsg.ChatID.Int64()))

		s.Require().NoError(err)

		s.Require().Positive(actualMsg.ID)
		s.Require().Equal(inputMsg.ChatID.Int64(), actualMsg.ChatID)
		s.Require().Equal(int64(inputMsg.ReplyMsgID.Int()), actualMsg.ReplyMessageID.Int64)
		s.Require().Equal(inputMsg.Text, actualMsg.Text)
		s.Require().Equal(model.MessageTypePlain, actualMsg.Type)
		s.Require().Nil(actualMsg.Payload)
		s.Require().Equal(jsonb.NewString("[]"), actualMsg.Button)
		s.Require().False(actualMsg.IsDispatched)
	})

	s.Run("store markdown text", func() {
		inputMsg := messageprocessor.Message{
			ChatID:     msginfo.ChatIDFromInt(2),
			ReplyMsgID: msginfo.MessageIDFromInt(2),
			Text:       "*test text*",
			Type:       messageprocessor.MessageTypeMarkdown,
		}

		err := s.pgDB.SendMessage(s.T().Context(), inputMsg)

		s.Require().NoError(err)

		actualMsg, err := s.pgDB.TestGetOutboxMessageByChatID(s.T().Context(), int(inputMsg.ChatID.Int64()))

		s.Require().NoError(err)

		s.Require().Positive(actualMsg.ID, 0)
		s.Require().Equal(inputMsg.ChatID.Int64(), actualMsg.ChatID)
		s.Require().Equal(int64(inputMsg.ReplyMsgID.Int()), actualMsg.ReplyMessageID.Int64)
		s.Require().Equal(inputMsg.Text, actualMsg.Text)
		s.Require().Equal(model.MessageTypeMarkdown, actualMsg.Type)
		s.Require().Nil(actualMsg.Payload)
		s.Require().Equal(jsonb.NewString("[]"), actualMsg.Button)
		s.Require().False(actualMsg.IsDispatched)
	})
}
