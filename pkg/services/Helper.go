package services

import (
	"fmt"
)

type Helper struct {
}

type FormattedDuration struct {
	Hours   int64
	Minutes int64
	Seconds int
}

func FormatTimerDuration(duration int64) string {
	hours := duration / 3600
	remainder := duration % 3600
	minutes := remainder / 60
	seconds := remainder % 60

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}
