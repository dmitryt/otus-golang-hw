package calendar

import (
	"errors"
)

var ErrDateBusy = errors.New("another event starts on the same date")
