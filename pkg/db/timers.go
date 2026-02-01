package db

import (
	"fmt"
	"log"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"shadowflade/timers/pkg/interfaces"
	"shadowflade/timers/pkg/services"
)

type Timer struct {
	TableName string
}

func (this *Timer) Create() *Timer {
	return &Timer{
		TableName: "timers",
	}
}

func (this *Timer) GetAllUsersTimers(userID int) []interfaces.Timer {
	db := Db{}
	err := db.Connect()
	if err != nil {
		panic(err.Error())
	}
	tx := db.Db.MustBegin()

	userId := strconv.Itoa(userID)
	fmt.Println(userID)
	res, err := db.Db.Queryx(fmt.Sprintf("select * from timers where user_id = %s order by start", userId))

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
		var duration int64
		duration = 0
		if userTimer.RunningSince.Valid {
			duration = time.Now().Unix() - userTimer.RunningSince.Time.Unix()
		} else {
			duration = userTimer.Duration
		}
		userTimer.FormattedDuration = services.FormatTimerDuration(duration)

		userTimers = append(userTimers, userTimer)
	}
	// fmt.Printf("User timers: %v\n", userTimers)

	tx.Commit()

	return userTimers
}

func (this *Db) CreateTimer(timer interfaces.Timer) (int64, error) {
	db := Db{}
	err := db.Connect()
	if err != nil {
		return 0, fmt.Errorf("failed to connect to database: %w", err)
	}

	tx, err := db.Db.Beginx() //why use BeginxX and not MustBegin?
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

func (this *Db) StartTimer(timerId int, startTime int64) (int64, error) {
	db := Db{}
	err := db.Connect()
	if err != nil {
		return 0, fmt.Errorf("failed to connect to database: %w", err)
	}
	sTime := time.Unix(startTime, 0)
	query := `update timers set running_since = ?, date_modified = ? where id = ?;`
	tx, err := db.Db.Beginx()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if tx != nil {
			tx.Rollback()
			fmt.Println("Ya dolboeb")
		}
	}()
	result, err := tx.Exec(query, sTime, sTime, timerId)
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
// Параметры:
//
//   - timerId: id таймера
//
//   - pauseTime - время остановки таймера, которое приходит с фронта, в милисекундах
//
// Returns
//
//	newDuration - int64
func (this *Db) PauseTimer(timerId int, pauseTime int64) (int64, error) {
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
	fmt.Printf(
		"elapsed: #%d, \nnewDuratino: %d, \nrunning since: %s, \nnow: %s",
		elapsedSeconds,
		newDuration,
		userTimer.RunningSince.Time,
		time.Now(),
	)

	newPausedAt := time.Now()

	now := time.Now()
	pausedAt := time.Unix(pauseTime/1000, 0) //TODO[check]:int64 делится на 1000 ?
	updateTimerDurationQuery := `
		UPDATE timers
		SET duration = ?, paused_at = ?, running_since = null, date_modified = ?
		WHERE id = ?`

	tx, err := db.Db.Beginx()
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
		pausedAt,
		now,
		timerId,
	)
	if err != nil {
		log.Panicf(
			"Query: %s \n failed to update timer: %s.\nDuration: %d,\nuserTimer.Duration: %d,\n elapsedSeconds: %d.\n Timer id : %d",
			updateTimerDurationQuery,
			err.Error(),
			newDuration,
			userTimer.Duration,
			elapsedSeconds,
			userTimer.Id,
		)
	}

	if err := tx.Commit(); err != nil {
		log.Panicf("%s", err.Error())
	}
	tx = nil // prevent rollback

	// Optional: get affected rows
	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("Updated %d rows\n", rowsAffected)

	return newDuration, nil

}

func (this *Db) StopTimer(timerId int) (int64, error) {
	db := Db{}
	err := db.Connect()
	if err != nil {
		return 0, fmt.Errorf("failed to connect to database: %w", err)
	}

	tx, err := db.Db.Beginx()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if tx != nil {
			tx.Rollback()
		}
	}()
	stopTimerQuery := `update timers set duration = 0, running_since = null where id = ?`
	result, err := tx.Exec(
		stopTimerQuery,
		timerId,
	)
	if err != nil {
		log.Panicf("Query: %s \n failed to update timer.\n Timer id : %s",
			stopTimerQuery,
			err,
		)
	}

	if err := tx.Commit(); err != nil {
		log.Panicf("failed to commit: %s", err.Error())
	}
	tx = nil // prevent rollback

	// Optional: get affected rows
	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("Updated %d rows\n", rowsAffected)

	return rowsAffected, nil

}

func (this *Db) GetTimerById(timerId int) interfaces.Timer {
	db := Db{}
	err := db.Connect()
	res, err := db.Db.Queryx(
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

func (this *Db) AddOrUpdateTimerColor(timerId int, color string) (int64, error) {
	db := Db{}
	err := db.Connect()
	if err != nil {
		log.Fatalf("Could not connect ot database updating timer color: %s", err.Error())
	}
	tx, err := db.Db.Beginx()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if tx != nil {
			tx.Rollback()
		}
	}()

	updateTitleQuery := `insert into timers (id, color) values (?,?) on duplicate key update color = ?;`

	result, err := tx.Exec(updateTitleQuery, timerId, color, color)
	rowsAffected, _ := result.RowsAffected()

	tx = nil

	return rowsAffected, nil
}

func (this *Db) UpdateTitle(title string, timerId int) (int64, error) {
	db := Db{}
	err := db.Connect()
	if err != nil {
		log.Fatalf("Could not connect ot database updating title: %s", err.Error())
	}
	tx, err := db.Db.Beginx()
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
	lastUpdatedId, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("rows affected %d, id affected %d", rowsAffected, lastUpdatedId)
	tx.Commit()

	tx = nil

	return rowsAffected, nil
}

func (this *Db) RefreshTimer(timerId int) (int64, error) {
	db := Db{}
	err := db.Connect()
	if err != nil {
		log.Fatalf("Could not connect ot database updating title: %s", err.Error())
	}
	tx, err := db.Db.Beginx()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if tx != nil {
			tx.Rollback()
		}
	}()
	now := time.Now()
	refreshTimerQuery := `update timers set running_since = ?, date_modified = ?, start = ?, duration = 0, `
	result, err := tx.Exec(refreshTimerQuery, now, now, now)

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
	tx, err := db.Db.Beginx()
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
		log.Fatalf(
			"Could not delete timer with id: %d. Error: %s",
			timerId,
			err.Error(),
		)
	}
	rowsAffected, err := result.RowsAffected()
	tx.Commit()
	tx = nil

	return int(rowsAffected), nil
}
