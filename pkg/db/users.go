package db

import (
	"log"

	"shadowflade/timers/pkg/interfaces"
)

type User struct {
	TableName string
}

func (this *User) Create() *User {
	return &User{
		TableName: "users",
	}
}

func (this *User) CreateUser(user interfaces.User) int64 {
	db := Db{}
	db.Connect()
	tx := db.Db.MustBegin()
	query := `
	insert into users (name, password) values (:name, :password)
	`

	res, err := db.Db.NamedExec(query, user)
	if err != nil {
		log.Fatalf("Error on creating user. Query: %s. Err: %s", query, err.Error())
	}
	newId, err := res.LastInsertId()
	if err != nil {
		panic(err.Error())
	}
	tx.Commit()
	return newId
}
