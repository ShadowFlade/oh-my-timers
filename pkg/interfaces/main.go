package interfaces

type Timer struct {
	Start  string `db:"start"`
	End    string `db:"end"`
	UserID int32  `db:"user_id"`
	Title  string `db:"title"`
	Color  string `db:"color"`
}

type TimerTemplate struct {
	Items        []Timer
	IsMoreThan10 bool
	UserID       int
}
