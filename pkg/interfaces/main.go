package interfaces

import (
	"database/sql"
	"shadowflade/timers/pkg/services"
	"time"
)

type Timer struct {
	DateInserted      sql.NullTime `db:"date_inserted"`
	DateModified      sql.NullTime `db:"date_modified"`
	StartTime         time.Time    `db:"start"`
	EndTime           sql.NullTime `db:"end"`
	Duration          int64        `db:"duration"`
	PausedAt          sql.NullTime `db:"paused_at"`
	RunningSince      sql.NullTime `db:"running_since"`
	UserID            int32        `db:"user_id"`
	Title             string       `db:"title"`
	Color             string       `db:"color"`
	Id                int64        `db:"id"`
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

	timer.StartTime = time.Now()
	// timer.StartTime.Valid = true

	timer.DateInserted.Time = time.Now()
	timer.DateInserted.Valid = true

	// fmt.Printf("%#v", timer)

	formattedDuration := services.FormatTimerDuration(timer.Duration)
	timer.FormattedDuration = formattedDuration
	return timer
}

type User struct {
	id       int64  `db:"id"`
	Name     string `db:"name"`
	uuid     string `db:"uuid"`
	Password string `db:"password"`
}

func NewUser(name string, password string) User {
	user := User{}
	user.Name = name
	user.Password = password
	return user
}
