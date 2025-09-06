package db

import (
	"fmt"
	"log"
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

	query := fmt.Sprintf("INSERT INTO timers (start, end, user_id, title, color) VALUES ('%s', '%s', %d, '%s', '%s');", timer.StartTime, timer.EndTime, timer.UserID, timer.Title, timer.Color)
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
	fmt.Println(userTimer, " user timer")

	var start time.Time

	if userTimer.RunningSince.Valid {
		start = userTimer.RunningSince.Time
	} else if userTimer.StartTime.Valid {
		start = userTimer.StartTime.Time
	}

	// Calculate elapsed time
	elapsedSeconds := int(time.Since(start).Seconds())
	newDuration := userTimer.Duration + int64(elapsedSeconds)

	newPausedAt := time.Now()

	updateTimerDurationQuery := `
		UPDATE timers
		SET duration = ?, paused_at = ?
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

	fmt.Println(newDuration, newPausedAt, timerId)
	result, err := tx.Exec(
		updateTimerDurationQuery,
		newDuration,
		newPausedAt,
		timerId,
	)
	if err != nil {
		log.Panicf("Query: %s failed to update timer: %w", updateTimerDurationQuery, err)
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
