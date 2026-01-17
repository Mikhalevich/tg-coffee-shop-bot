package postgres_test

import (
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/internal/jsonb"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/internal/model"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor/button"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/cart"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
)

func (s *PostgresSuit) TestSendMessageWitoutPayloadAndButtons() {
	s.Run("store empty text", func() {
		inputMsg := messageprocessor.Message{
			ChatID: msginfo.ChatIDFromInt(1),
			Type:   messageprocessor.MessageTypePlain,
		}

		err := s.pgDB.SendMessage(s.T().Context(), inputMsg)
		s.Require().NoError(err)

		actualMsg, err := s.pgDB.TestGetOutboxMessageByChatID(s.T().Context(), int(inputMsg.ChatID.Int64()))

		s.Require().NoError(err)

		s.Require().Positive(actualMsg.ID)
		s.Require().Equal(inputMsg.ChatID.Int64(), actualMsg.ChatID)
		s.Require().Zero(actualMsg.ReplyMessageID.Int64)
		s.Require().False(actualMsg.ReplyMessageID.Valid)
		s.Require().Empty(inputMsg.Text)
		s.Require().Equal(model.MessageTypePlain, actualMsg.Type)
		s.Require().Nil(actualMsg.Payload)
		s.Require().Equal(jsonb.NewString("[]"), actualMsg.Button)
		s.Require().False(actualMsg.IsDispatched)
		s.Require().NotZero(actualMsg.CreatedAt)
		s.Require().Zero(actualMsg.DispatchedAt.Time)
		s.Require().False(actualMsg.DispatchedAt.Valid)
	})

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
		s.Require().True(actualMsg.ReplyMessageID.Valid)
		s.Require().Equal(inputMsg.Text, actualMsg.Text)
		s.Require().Equal(model.MessageTypePlain, actualMsg.Type)
		s.Require().Nil(actualMsg.Payload)
		s.Require().Equal(jsonb.NewString("[]"), actualMsg.Button)
		s.Require().False(actualMsg.IsDispatched)
		s.Require().NotZero(actualMsg.CreatedAt)
		s.Require().Zero(actualMsg.DispatchedAt.Time)
		s.Require().False(actualMsg.DispatchedAt.Valid)
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
		s.Require().True(actualMsg.ReplyMessageID.Valid)
		s.Require().Equal(inputMsg.Text, actualMsg.Text)
		s.Require().Equal(model.MessageTypeMarkdown, actualMsg.Type)
		s.Require().Nil(actualMsg.Payload)
		s.Require().Equal(jsonb.NewString("[]"), actualMsg.Button)
		s.Require().False(actualMsg.IsDispatched)
		s.Require().NotZero(actualMsg.CreatedAt)
		s.Require().Zero(actualMsg.DispatchedAt.Time)
		s.Require().False(actualMsg.DispatchedAt.Valid)
	})
}

func (s *PostgresSuit) TestSendMessageWithButtons() {
	s.Run("store plain text with empty buttons", func() {
		inputMsg := messageprocessor.Message{
			ChatID:     msginfo.ChatIDFromInt(7),
			ReplyMsgID: msginfo.MessageIDFromInt(7),
			Text:       "test text",
			Type:       messageprocessor.MessageTypePlain,
			Buttons:    []button.ButtonRow{},
		}

		err := s.pgDB.SendMessage(s.T().Context(), inputMsg)

		s.Require().NoError(err)

		actualMsg, err := s.pgDB.TestGetOutboxMessageByChatID(s.T().Context(), int(inputMsg.ChatID.Int64()))

		s.Require().NoError(err)

		s.Require().Positive(actualMsg.ID)
		s.Require().Equal(inputMsg.ChatID.Int64(), actualMsg.ChatID)
		s.Require().Equal(int64(inputMsg.ReplyMsgID.Int()), actualMsg.ReplyMessageID.Int64)
		s.Require().True(actualMsg.ReplyMessageID.Valid)
		s.Require().Equal(inputMsg.Text, actualMsg.Text)
		s.Require().Equal(model.MessageTypePlain, actualMsg.Type)
		s.Require().Nil(actualMsg.Payload)
		s.Require().Equal(jsonb.NewString("[]"), actualMsg.Button)
		s.Require().False(actualMsg.IsDispatched)
		s.Require().NotZero(actualMsg.CreatedAt)
		s.Require().Zero(actualMsg.DispatchedAt.Time)
		s.Require().False(actualMsg.DispatchedAt.Valid)
	})

	s.Run("store plain text with one button without payload", func() {
		inputMsg := messageprocessor.Message{
			ChatID:     msginfo.ChatIDFromInt(8),
			ReplyMsgID: msginfo.MessageIDFromInt(8),
			Text:       "test text",
			Type:       messageprocessor.MessageTypePlain,
			Buttons: []button.ButtonRow{
				{
					{
						ID:        button.IDFromString("button id"),
						ChatID:    msginfo.ChatIDFromInt(8),
						Caption:   "button caption",
						Operation: button.OperationCartCancel,
						Pay:       true,
					},
				},
			},
		}

		err := s.pgDB.SendMessage(s.T().Context(), inputMsg)

		s.Require().NoError(err)

		actualMsg, err := s.pgDB.TestGetOutboxMessageByChatID(s.T().Context(), int(inputMsg.ChatID.Int64()))

		s.Require().NoError(err)

		s.Require().Positive(actualMsg.ID)
		s.Require().Equal(inputMsg.ChatID.Int64(), actualMsg.ChatID)
		s.Require().Equal(int64(inputMsg.ReplyMsgID.Int()), actualMsg.ReplyMessageID.Int64)
		s.Require().True(actualMsg.ReplyMessageID.Valid)
		s.Require().Equal(inputMsg.Text, actualMsg.Text)
		s.Require().Equal(model.MessageTypePlain, actualMsg.Type)
		s.Require().Nil(actualMsg.Payload)
		//nolint:lll
		s.Require().Equal(jsonb.NewString(`[[{"ID": "button id", "Pay": true, "ChatID": 8, "Caption": "button caption", "Payload": null, "Operation": "CartCancel"}]]`), actualMsg.Button)
		s.Require().False(actualMsg.IsDispatched)
		s.Require().NotZero(actualMsg.CreatedAt)
		s.Require().Zero(actualMsg.DispatchedAt.Time)
		s.Require().False(actualMsg.DispatchedAt.Valid)
	})

	s.Run("store plain text with one button with payload", func() {
		cartCancelBtn, err := button.CartCancel(
			msginfo.ChatIDFromInt(9),
			"cart button caption",
			cart.IDFromString("cart id"),
		)
		s.Require().NoError(err)

		inputMsg := messageprocessor.Message{
			ChatID:     msginfo.ChatIDFromInt(9),
			ReplyMsgID: msginfo.MessageIDFromInt(9),
			Text:       "test text",
			Type:       messageprocessor.MessageTypePlain,
			Buttons: []button.ButtonRow{
				{
					cartCancelBtn,
				},
			},
		}

		err = s.pgDB.SendMessage(s.T().Context(), inputMsg)
		s.Require().NoError(err)

		actualMsg, err := s.pgDB.TestGetOutboxMessageByChatID(s.T().Context(), int(inputMsg.ChatID.Int64()))

		s.Require().NoError(err)

		//nolint:lll
		expectedButtonJSON := fmt.Sprintf(`[[{"ID": "%s", "Pay": false, "ChatID": 9, "Caption": "cart button caption", "Payload": "KX8DAQERQ2FydENhbmNlbFBheWxvYWQB/4AAAQEBBkNhcnRJRAEMAAAADP+AAQdjYXJ0IGlkAA==", "Operation": "CartCancel"}]]`,
			cartCancelBtn.ID,
		)

		s.Require().Positive(actualMsg.ID)
		s.Require().Equal(inputMsg.ChatID.Int64(), actualMsg.ChatID)
		s.Require().Equal(int64(inputMsg.ReplyMsgID.Int()), actualMsg.ReplyMessageID.Int64)
		s.Require().True(actualMsg.ReplyMessageID.Valid)
		s.Require().Equal(inputMsg.Text, actualMsg.Text)
		s.Require().Equal(model.MessageTypePlain, actualMsg.Type)
		s.Require().Nil(actualMsg.Payload)
		s.Require().Equal(jsonb.NewString(expectedButtonJSON), actualMsg.Button)
		s.Require().False(actualMsg.IsDispatched)
		s.Require().NotZero(actualMsg.CreatedAt)
		s.Require().Zero(actualMsg.DispatchedAt.Time)
		s.Require().False(actualMsg.DispatchedAt.Valid)
	})

	s.Run("store plain text with two buttons with payload", func() {
		cartCancelBtn, err := button.CartCancel(
			msginfo.ChatIDFromInt(10),
			"cart cancel caption",
			cart.IDFromString("cart cancel id"),
		)
		s.Require().NoError(err)

		cartConfirmBtn, err := button.CartConfirm(
			msginfo.ChatIDFromInt(10),
			"cart confirm caption",
			cart.IDFromString("cart cancel id"),
			currency.IDFromInt(1),
		)
		s.Require().NoError(err)

		inputMsg := messageprocessor.Message{
			ChatID:     msginfo.ChatIDFromInt(10),
			ReplyMsgID: msginfo.MessageIDFromInt(10),
			Text:       "test text",
			Type:       messageprocessor.MessageTypePlain,
			Buttons: []button.ButtonRow{
				{
					cartCancelBtn,
					cartConfirmBtn,
				},
			},
		}

		err = s.pgDB.SendMessage(s.T().Context(), inputMsg)
		s.Require().NoError(err)

		actualMsg, err := s.pgDB.TestGetOutboxMessageByChatID(s.T().Context(), int(inputMsg.ChatID.Int64()))

		s.Require().NoError(err)

		//nolint:lll
		expectedButtonJSON := fmt.Sprintf(
			`[[{"ID": "%s", "Pay": false, "ChatID": 10, "Caption": "cart cancel caption", "Payload": "KX8DAQERQ2FydENhbmNlbFBheWxvYWQB/4AAAQEBBkNhcnRJRAEMAAAAE/+AAQ5jYXJ0IGNhbmNlbCBpZAA=", "Operation": "CartCancel"}, {"ID": "%s", "Pay": false, "ChatID": 10, "Caption": "cart confirm caption", "Payload": "Ov+BAwEBEkNhcnRDb25maXJtUGF5bG9hZAH/ggABAgEGQ2FydElEAQwAAQpDdXJyZW5jeUlEAQQAAAAV/4IBDmNhcnQgY2FuY2VsIGlkAQIA", "Operation": "CartConfirm"}]]`,
			cartCancelBtn.ID,
			cartConfirmBtn.ID,
		)

		s.Require().Positive(actualMsg.ID)
		s.Require().Equal(inputMsg.ChatID.Int64(), actualMsg.ChatID)
		s.Require().Equal(int64(inputMsg.ReplyMsgID.Int()), actualMsg.ReplyMessageID.Int64)
		s.Require().True(actualMsg.ReplyMessageID.Valid)
		s.Require().Equal(inputMsg.Text, actualMsg.Text)
		s.Require().Equal(model.MessageTypePlain, actualMsg.Type)
		s.Require().Nil(actualMsg.Payload)
		s.Require().Equal(jsonb.NewString(expectedButtonJSON), actualMsg.Button)
		s.Require().False(actualMsg.IsDispatched)
		s.Require().NotZero(actualMsg.CreatedAt)
		s.Require().Zero(actualMsg.DispatchedAt.Time)
		s.Require().False(actualMsg.DispatchedAt.Valid)
	})
}
