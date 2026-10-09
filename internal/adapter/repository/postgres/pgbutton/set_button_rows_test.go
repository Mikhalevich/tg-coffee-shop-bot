package pgbutton_test

import (
	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/button"
)

func (s *ButtonSuit) TestSetButtonRows() {
	s.Run("multiple rows", func() {
		var (
			ctx    = s.T().Context()
			first  = button.MustCreateButton("First", button.OperationCartConfirm, button.WithStyle(button.StyleSuccess))
			second = button.MustCreateButton("Second", button.OperationCartCancel)
			share  = button.MustShareButton("Share", "https://example.com")
		)

		err := s.pgButton.SetButtonRows(ctx, button.Row(first, second), button.Row(share))

		s.Require().NoError(err)
		s.Require().Equal(3, s.buttonsCount())

		for _, expected := range []button.Button{first, second, share} {
			actual, err := s.pgButton.GetButton(ctx, expected.ID)
			s.Require().NoError(err)
			s.Require().Equal(&expected, actual)
		}
	})

	s.Run("no rows", func() {
		ctx := s.T().Context()

		s.Require().NoError(s.pgButton.SetButtonRows(ctx))
		s.Require().NoError(s.pgButton.SetButtonRows(ctx, button.Row()))
		s.Require().Equal(0, s.buttonsCount())
	})

	s.Run("pay button skipped", func() {
		var (
			ctx = s.T().Context()
			btn = button.MustCreateButton("Cancel", button.OperationOrderCancel)
		)

		s.Require().NoError(s.pgButton.SetButtonRows(ctx, button.Row(button.Pay("Pay")), button.Row(btn)))
		s.Require().NoError(s.pgButton.SetButtonRows(ctx, button.Row(button.Pay("Pay"))))
		s.Require().Equal(1, s.buttonsCount())
	})

	s.Run("overwrite existing", func() {
		var (
			ctx     = s.T().Context()
			btn     = button.MustCreateButton("Old", button.OperationCartConfirm)
			updated = btn
		)

		updated.Caption = "New"

		s.Require().NoError(s.pgButton.SetButtonRows(ctx, button.Row(btn)))
		s.Require().NoError(s.pgButton.SetButtonRows(ctx, button.Row(updated)))

		actual, err := s.pgButton.GetButton(ctx, btn.ID)

		s.Require().NoError(err)
		s.Require().Equal(&updated, actual)
	})
}

func (s *ButtonSuit) buttonsCount() int {
	var (
		ctx   = s.T().Context()
		count int
	)

	s.Require().NoError(sqlx.GetContext(ctx, s.transactor.ExtContext(ctx), &count, "SELECT COUNT(*) FROM buttons"))

	return count
}
