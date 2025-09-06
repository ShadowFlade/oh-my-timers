package db

import (
	"fmt"
	"log"
	"shadowflade/timers/pkg/global"
	"shadowflade/timers/pkg/interfaces"
	"strconv"
	"time"

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
		err := res.StructScan(&userTimer)
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
		return 0, fmt.Errorf("failed to connect to database: %w", err)
	}

	tx, err := db.db.Beginx() // Use Beginx instead of MustBegin for better error handling
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if tx != nil {
			tx.Rollback()
		}
	}()

	query := fmt.Sprintf("INSERT INTO timers (start, end, user_id, title, color) VALUES ('%s', '%s', %d, '%s', '%s');", timer.Start, timer.End, timer.UserID, timer.Title, timer.Color)
	result, err := tx.Exec(query)

	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %w", err)
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return 0, fmt.Errorf("no rows were inserted")
	}

	newId, err := result.LastInsertId()

	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Set tx to nil to prevent defer from rolling back
	tx = nil
	return newId, nil
}

func (this *Db) SetDuration(timerId int) (int64, error) {
	db := Db{}
	err := db.Connect()
	if err != nil {
		return 0, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err != nil {
		panic(err.Error())
	}

	userTimer := this.GetTimerById(timerId)

	runningSince, _ := time.Parse(
		global.MYSQL_DATETIME_FORMAT,
		userTimer.RunningSince,
	)

	newDuration := runningSince.Add(time.Since(time.Now())).Unix() + userTimer.Duration

	updateTimerDurationQuery := fmt.Sprintf("update %s where id = %s set duration = %s", this.TimersTable, timerId, newDuration)

	tx, err := db.db.Beginx() // Use Beginx instead of MustBegin for better error handling
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if tx != nil {
			tx.Rollback()
		}
	}()

	result, err := tx.Exec(updateTimerDurationQuery)

	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %w", err)
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return 0, fmt.Errorf("no rows were inserted")
	}

	tx = nil

	return rowsAffected, nil

}

func (this *Db) GetTimerById(timerId int) interfaces.Timer {
	db := Db{}
	err := db.Connect()
	res, err := db.db.Queryx(
		fmt.Sprintf("select * from timers where id = %s", timerId),
	)

	if err != nil {
		panic(err.Error())
	}
	var userTimer interfaces.Timer

	err = res.StructScan(&userTimer)

	if err != nil {
		log.Fatal(err.Error())
	}

	return userTimer
}
