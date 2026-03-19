package db

type Category struct {
	TableName string
}

func (this *Category) New() *Category {
	return &Category{
		TableName: "categories",
	}
}
