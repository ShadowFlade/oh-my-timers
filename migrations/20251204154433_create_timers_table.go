package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	PkgDb "shadowflade/timers/pkg/db"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateTimersTable, downCreateTimersTable)
}

func upCreateTimersTable(ctx context.Context, tx *sql.Tx) error {
	dbUser := PkgDb.Timer{}
	timerDb := dbUser.Create()

	timersTableName := timerDb.TableName
	var isTimersTableAlreadyExists bool

	err := tx.QueryRow(fmt.Sprintf("SHOW TABLES LIKE %s;", timersTableName)).Scan(&isTimersTableAlreadyExists)

	if err == nil {
		fmt.Println("Table timers already exists")
		return nil
	}

	query := fmt.Sprintf(`
	CREATE TABLE %s (
		id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		user_id INT DEFAULT NULL,
		start DATETIME DEFAULT NULL,
		end DATETIME DEFAULT NULL,
		date_inserted DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		date_modified DATETIME DEFAULT NULL,
		duration INT UNSIGNED NOT NULL DEFAULT 0,
		paused_at DATETIME DEFAULT NULL,
		running_since DATETIME DEFAULT NULL,
		title VARCHAR(255) DEFAULT NULL,
		color VARCHAR(255) DEFAULT NULL
	)
	`, timersTableName)

	res, createTableError := tx.Exec(query)

	if createTableError != nil {
		tx.Rollback()
		log.Panic(err.Error())
	}
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		tx.Rollback()
		log.Panic(err.Error())
	}

	if rowsAffected > 0 {
		tx.Commit()
		return nil
	}

	return nil
}

func downCreateTimersTable(ctx context.Context, tx *sql.Tx) error {
	dbTimer := PkgDb.Timer{}
	timerDb := dbTimer.Create()

	timersTableName := timerDb.TableName
	fmt.Print(timersTableName, "table name")

	rows, err := tx.Exec(
		"SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?;",
		timersTableName,
	)

	if err != nil {
		tx.Rollback()
		log.Panic(err.Error())
	}

	rowsAffected, _ := rows.RowsAffected()

	if rowsAffected == 0 {
		return nil
	}

	query := `
	drop table if exists ?;
	`

	res, err := tx.Exec(query, timersTableName)

	if err != nil {
		tx.Rollback()
		log.Panic(err.Error())
	}

	rowsAffected, _ = res.RowsAffected()

	if rowsAffected > 0 {
		tx.Commit()
	}

	return nil
}
