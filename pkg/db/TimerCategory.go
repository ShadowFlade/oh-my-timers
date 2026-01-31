package db

type TimerCategory struct {
	TableName string
}

func (this *TimerCategory) New() *TimerCategory {
	return &TimerCategory{
		TableName: "timer_categories",
	}
}
