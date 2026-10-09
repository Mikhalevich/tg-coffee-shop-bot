package pgbutton_test

import (
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/button"
)

func (s *ButtonSuit) TestGetButton() {
	s.Run("success", func() {
		var (
			ctx      = s.T().Context()
			expected = button.MustCreateButton(
				"Cancel",
				button.OperationOrderCancel,
				button.WithStyle(button.StyleDanger),
				button.WithDeleteAfterProcess(),
				button.WithPayload(42),
			)
		)

		s.Require().NoError(s.pgButton.SetButtonRows(ctx, button.Row(expected)))

		actual, err := s.pgButton.GetButton(ctx, expected.ID)

		s.Require().NoError(err)
		s.Require().Equal(&expected, actual)
	})

	s.Run("deleted after get", func() {
		var (
			ctx = s.T().Context()
			btn = button.MustCreateButton("Categories", button.OperationCartViewCategories)
		)

		s.Require().NoError(s.pgButton.SetButtonRows(ctx, button.Row(btn)))

		_, err := s.pgButton.GetButton(ctx, btn.ID)
		s.Require().NoError(err)

		actual, err := s.pgButton.GetButton(ctx, btn.ID)

		s.Require().Nil(actual)
		s.Require().EqualError(err, "button not found")
		s.Require().True(s.pgButton.IsNotFoundError(err))
	})

	s.Run("not found", func() {
		ctx := s.T().Context()

		actual, err := s.pgButton.GetButton(ctx, button.IDFromString("unknown"))

		s.Require().Nil(actual)
		s.Require().EqualError(err, "button not found")
		s.Require().True(s.pgButton.IsNotFoundError(err))
	})
}
