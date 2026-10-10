package pgdailyposition_test

import (
	"sync"
	"time"
)

func (s *DailyPositionSuit) TestPosition() {
	s.Run("increments within a day", func() {
		var (
			ctx = s.T().Context()
			day = time.Date(2026, time.October, 10, 9, 0, 0, 0, time.UTC)
		)

		for expected := 1; expected <= 3; expected++ {
			actual, err := s.pgDailyPosition.Position(ctx, day)

			s.Require().NoError(err)
			s.Require().Equal(expected, actual)
		}
	})

	s.Run("same date different time shares counter", func() {
		var (
			ctx     = s.T().Context()
			morning = time.Date(2026, time.October, 10, 0, 0, 1, 0, time.UTC)
			evening = time.Date(2026, time.October, 10, 23, 59, 59, 0, time.UTC)
		)

		first, err := s.pgDailyPosition.Position(ctx, morning)
		s.Require().NoError(err)
		s.Require().Equal(1, first)

		second, err := s.pgDailyPosition.Position(ctx, evening)
		s.Require().NoError(err)
		s.Require().Equal(2, second)
	})

	s.Run("different days have independent counters", func() {
		var (
			ctx   = s.T().Context()
			today = time.Date(2026, time.October, 10, 12, 0, 0, 0, time.UTC)
			next  = today.AddDate(0, 0, 1)
		)

		_, err := s.pgDailyPosition.Position(ctx, today)
		s.Require().NoError(err)

		pos, err := s.pgDailyPosition.Position(ctx, today)
		s.Require().NoError(err)
		s.Require().Equal(2, pos)

		pos, err = s.pgDailyPosition.Position(ctx, next)
		s.Require().NoError(err)
		s.Require().Equal(1, pos)
	})

	s.Run("concurrent calls return unique positions", func() {
		const callsCount = 20

		var (
			ctx       = s.T().Context()
			day       = time.Date(2026, time.October, 10, 12, 0, 0, 0, time.UTC)
			waitGroup sync.WaitGroup
			positions = make([]int, callsCount)
			errs      = make([]error, callsCount)
		)

		for i := range callsCount {
			waitGroup.Go(func() {
				positions[i], errs[i] = s.pgDailyPosition.Position(ctx, day)
			})
		}

		waitGroup.Wait()

		seen := make(map[int]struct{}, callsCount)

		for i := range callsCount {
			s.Require().NoError(errs[i])
			s.Require().GreaterOrEqual(positions[i], 1)
			s.Require().LessOrEqual(positions[i], callsCount)

			seen[positions[i]] = struct{}{}
		}

		s.Require().Len(seen, callsCount)
	})
}
