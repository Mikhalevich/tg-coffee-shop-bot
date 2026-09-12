package id

import (
	"fmt"
	"strconv"
)

type IntID int

func (i IntID) Int() int {
	return int(i)
}

func (i IntID) String() string {
	return strconv.FormatInt(int64(i), 10)
}

func IntIDFromInt(i int) IntID {
	return IntID(i)
}

func IntIDFromString(s string) (IntID, error) {
	intID, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse int: %w", err)
	}

	return IntID(intID), nil
}
