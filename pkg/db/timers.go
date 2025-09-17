package db

import (
	"fmt"
	"log"
	"shadowflade/timers/pkg/interfaces"
	"shadowflade/timers/pkg/services"
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
	fmt.Println(userID)
	res, err := db.db.Queryx(fmt.Sprintf("select * from timers where user_id = %s order by start", userId))

	if err != nil {
		panic(err.Error())
	}

	var userTimers []interfaces.Timer
	for res.Next() {
		userTimer := interfaces.NewTimer(0, "", "")
		err := res.StructScan(&userTimer)
		if err != nil {
			log.Fatal(err.Error())
		}
		userTimer.FormattedDuration = services.FormatTimerDuration(userTimer.Duration)
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

	query := `
	insert into timers (start, end, user_id, title, color, date_inserted, duration) values (:start, :end, :user_id, :title, :color, :date_inserted, :duration)
	`

	result, err := tx.NamedExec(query, timer)

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

func (this *Db) StartTimer(timerId int) (int64, error) {
	db := Db{}
	err := db.Connect()
	now := time.Now()
	if err != nil {
		return 0, fmt.Errorf("failed to connect to database: %w", err)
	}
	query := `update timers set running_since = ?, date_modified = ? where id = ?;`
	tx, err := db.db.Beginx()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if tx != nil {
			tx.Rollback()
			fmt.Println("Ya dolboeb")
		}
	}()
	result, err := tx.Exec(query, now, now, timerId)
	if err != nil {
		log.Panic(err.Error())
	}
	affectedRows, _ := result.RowsAffected()
	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}
	tx = nil
	return affectedRows, nil

}

// Ставит таймер на паузу (обновляет paused_at, running_since)
// Returns
//
//	newDuration - int64
func (this *Db) PauseTimer(timerId int) (int64, error) {
	db := Db{}
	err := db.Connect()
	if err != nil {
		return 0, fmt.Errorf("failed to connect to database: %w", err)
	}

	userTimer := this.GetTimerById(timerId)

	var start time.Time

	if userTimer.RunningSince.Valid {
		start = userTimer.RunningSince.Time
		fmt.Println("Using running since")
	} else if !userTimer.StartTime.IsZero() {
		start = userTimer.StartTime
		fmt.Println("Start time is not zero")
	} else {
		log.Panicln("You fucked up with time bro")
	}

	// Calculate elapsed time
	elapsedSeconds := int(time.Since(start).Seconds())
	newDuration := userTimer.Duration + int64(elapsedSeconds)
	fmt.Printf("elapsed: #%s, \nnewDuratino: %s, \nrunning since: %s, \nnow: %s", elapsedSeconds, newDuration, userTimer.RunningSince.Time, time.Now())

	newPausedAt := time.Now()

	now := time.Now()
	updateTimerDurationQuery := `
		UPDATE timers
		SET duration = ?, paused_at = ?, running_since = null, date_modified = ?
		WHERE id = ?`

	tx, err := db.db.Beginx()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if tx != nil {
			tx.Rollback()
		}
	}()

	fmt.Println(newDuration, newPausedAt, now, now, timerId)
	result, err := tx.Exec(
		updateTimerDurationQuery,
		newDuration,
		now,
		now,
		timerId,
	)
	if err != nil {
		log.Panicf("Query: %s \n failed to update timer: %w.\nDuration: %s,\nuserTimer.Duration: %s,\n elapsedSeconds: %s.\n Timer id : %s", updateTimerDurationQuery, err, newDuration, userTimer.Duration, elapsedSeconds, userTimer.Id)
		// log.Panicf(err.Error())
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Panicf("failed to commit: %w", err)
	}
	tx = nil // prevent rollback

	// Optional: get affected rows
	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("Updated %d rows\n", rowsAffected)

	return newDuration, nil

}

func (this *Db) GetTimerById(timerId int) interfaces.Timer {
	db := Db{}
	err := db.Connect()
	res, err := db.db.Queryx(
		fmt.Sprintf("select * from timers where id = %d;", int(timerId)),
	)

	if err != nil {
		panic(err.Error())
	}
	userTimer := interfaces.NewTimer(0, "", "")
	for res.Next() {
		err := res.StructScan(&userTimer)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	if err != nil {
		log.Fatal(err.Error())
	}

	return userTimer

}

func (this *Db) UpdateTitle(title string, timerId int) (int64, error) {
	db := Db{}
	err := db.Connect()
	if err != nil {
		log.Fatalf("Could not connect ot database updating title: %s", err.Error())
	}
	tx, err := db.db.Beginx()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if tx != nil {
			tx.Rollback()
		}
	}()
	updateTitleQuery := `update timers set title = ? where id = ?;`
	result, err := tx.Exec(updateTitleQuery, title, timerId)
	rowsAffected, _ := result.RowsAffected()

	tx = nil

	return rowsAffected, nil
}

func (this *Db) DeleteTimer(timerId int64) (int, error) {
	db := Db{}
	err := db.Connect()
	if err != nil {
		log.Fatalf("Could not connect ot database updating title: %s", err.Error())
	}
	tx, err := db.db.Beginx()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if tx != nil {
			tx.Rollback()
		}
	}()

	deleteTimerQuery := `delete from timers where id = ?`

	result, err := tx.Exec(deleteTimerQuery, timerId)
	if err != nil {
		log.Fatalln("Could not delete timer with id: %d. Error: %s", timerId, err.Error())
	}
	rowsAffected, err := result.RowsAffected()
	tx.Commit()
	tx = nil

	return int(rowsAffected), nil
}
