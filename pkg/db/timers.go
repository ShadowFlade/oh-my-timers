package db

import (
	"fmt"
	"log"
	"shadowflade/timers/pkg/interfaces"

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

func (this *Db) CreateTimer(timer interfaces.Timer) (int64, error) {
	db, err := this.connect()

	if err != nil {
		log.Fatal(err)
	}
	transaction := db.MustBegin()
	dbHelper := Helper{}
	query, err := dbHelper.GenerateInsertQuery("timers", timer)
	if err != nil {
		fmt.Printf("Error generating insert query")
		panic(err)
	}
	res := transaction.MustExec(query)

	newId, err := res.LastInsertId()

	if err != nil {
		fmt.Printf(err.Error())
		panic(err)
	}

	if newId > 0 {
		commitErr := transaction.Commit()
		if commitErr != nil {
			fmt.Printf(commitErr.Error())
			panic(commitErr.Error())
		}

		return newId, nil
	}

	return 0, nil

}
