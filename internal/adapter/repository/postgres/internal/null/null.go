package null

import (
	"database/sql"
)

func NullString(value string) sql.NullString {
	if value == "" {
		return sql.NullString{}
	}

	return sql.NullString{
		String: value,
		Valid:  true,
	}
}

func NullIntPositive(value int32) sql.NullInt32 {
	if value <= 0 {
		return sql.NullInt32{}
	}

	return sql.NullInt32{
		Int32: value,
		Valid: true,
	}
}
