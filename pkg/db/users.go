package db

import "fmt"

type User struct {
}

func (this *User) CreateUser(name string) int64 {
	db := Db{}
	db.Connect()
	tx := db.db.MustBegin()
	query := fmt.Sprintf("insert into %s (name) values ('%s');", db.UsersTable, name)

	res := db.db.MustExec(query)
	fmt.Printf(query)
	newId, err := res.LastInsertId()
	if err != nil {
		panic(err.Error())
	}
	tx.Commit()
	return newId
}
