package db

import (
	"errors"
	"fmt"
	"log"
	"os"
	"shadowflade/timers/pkg/interfaces"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofor-little/env"
	"github.com/jmoiron/sqlx"
)


type Db struct {
	db       *sqlx.DB
	tx       *sqlx.Tx
	dbName   string
	dbHost   string
	login    string
	password string
	cols     []string
}

func (this *Db) connect() (*sqlx.DB, error) {
	if err := env.Load("./.env"); err != nil {
		fmt.Println("error")
		panic(err)
	}

	this.login = env.Get("DB_LOGIN", "i")
	this.password = env.Get("DB_PASS", "fucked")
	this.dbName = env.Get("DB_NAME", "urmom")
	this.dbHost = env.Get("DB_HOST", "host")
	connectStr := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", this.login, this.password, this.dbName)
	db, err := sqlx.Connect("mysql", connectStr)

	if err != nil { //here i guess we can only make an assumption that db is not created yet, mb figure out later what type of error this is
		return db, err
	}

	this.db = db

	return db, nil
}

func (this *Db) createTimer(timer interfaces.Timer) {
	db, err := this.connect()

	if err != nil {
		log.Fatal(err)
	}
	transaction := db.MustBegin()
	res := transaction.MustExec("insert into timers (start

}
