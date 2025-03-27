package db

import (
	"fmt"
	"log"
	"shadowflade/timers/pkg/interfaces"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type Timer struct {
}

func (this *Timer) GetAllUsersTimers(userID int) []interfaces.Timer {
	db := Db{}
	err := db.Connect()
	if err != nil {
		panic(err.Error())
	}
	tx := db.db.MustBegin()

	userId := strconv.Itoa(userID)
	res, err := db.db.Queryx(fmt.Sprintf("select * from timers where user_id = %s", userId))

	if err != nil {
		panic(err.Error())
	}

	var userTimers []interfaces.Timer
	for res.Next() {
		var userTimer interfaces.Timer
		err := res.StructScan(userTimer)
		if err != nil {
			log.Fatal(err.Error())
		}
		userTimers = append(userTimers, userTimer)
	}
	tx.Commit()

	return userTimers
}

func (this *Db) CreateTimer(timer interfaces.Timer) (int64, error) {
	db := Db{}
	err := db.Connect()

	if err != nil {
		log.Fatal(err)
	}
	transaction := db.db.MustBegin()
	dbHelper := Helper{}
	query, err := dbHelper.GenerateInsertQuery(db.TimersTable, timer)
	fmt.Print(query, " QUERY!!!!")
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
