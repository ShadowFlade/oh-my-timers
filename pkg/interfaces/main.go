package interfaces

import (
	"database/sql"
	"shadowflade/timers/pkg/services"
)

type Timer struct {
	DateInserted      string         `db:"date_inserted"`
	DateModified      sql.NullString `db:"date_modified"`
	StartTime         sql.NullTime   `db:"start"`
	EndTime           sql.NullString `db:"end"`
	Duration          int64          `db:"duration"`
	PausedAt          sql.NullTime   `db:"paused_at"`
	RunningSince      sql.NullTime   `db:"running_since"`
	UserID            int32          `db:"user_id"`
	Title             string         `db:"title"`
	Color             string         `db:"color"`
	Id                int64          `db:"id"`
	FormattedDuration string
}

type TimerTemplate struct {
	Items        []Timer
	IsMoreThan10 bool
	UserID       int
}

func NewTimer(userId int32, title string, color string) Timer {
	timer := Timer{}
	if userId != 0 {
		timer.UserID = userId
	}
	timer.Title = title
	timer.Color = color

	formattedDuration := services.FormatTimerDuration(timer.Duration)
	timer.FormattedDuration = formattedDuration
	return timer
}
