package db

import (
	"fmt"

	"github.com/gofor-little/env"
	"github.com/jmoiron/sqlx"
)

type Db struct {
	db          *sqlx.DB
	dbName      string
	dbHost      string
	login       string
	password    string
	UsersTable  string
	TimersTable string
}

func (this *Db) Connect() error {
	if err := env.Load("./.env"); err != nil {
		fmt.Println("error")
		panic(err)
	}

	this.login = env.Get("DB_LOGIN", "i")
	this.password = env.Get("DB_PASS", "fucked")
	this.dbName = env.Get("DB_NAME", "urmom")
	this.dbHost = env.Get("DB_HOST", "host")
	connectStr := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?parseTime=true&loc=UTC", this.login, this.password, this.dbName)
	db, err := sqlx.Connect("mysql", connectStr)
	this.UsersTable = "users"
	this.TimersTable = "timers"
// 	db.Query("SET time_zone = '+03:00'")

	if err != nil {
		return err
	}

	this.db = db

	return nil
}
