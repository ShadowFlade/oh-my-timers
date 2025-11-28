package db

import (
	"log"

	"shadowflade/timers/pkg/interfaces"
)

type User struct {
}

func (this *User) CreateUser(user interfaces.User) int64 {
	db := Db{}
	db.Connect()
	tx := db.db.MustBegin()
	query := `
	insert into users (name, password) values (:name, :password)
	`

	res, err := db.db.NamedExec(query, user)
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
