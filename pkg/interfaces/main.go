package interfaces

type Timer struct {
	StartTime    string `db:"start"`
	EndTime      string `db:"end"`
	Duration     int64  `db:"duration"`
	PausedAt     string `db:"paused_at"`
	RunningSince string `db:"running_since"`
	UserID       int32  `db:"user_id"`
	Title        string `db:"title"`
	Color        string `db:"color"`
	Id           int64  `db:"id"`
}

type TimerTemplate struct {
	Items        []Timer
	IsMoreThan10 bool
	UserID       int
}
