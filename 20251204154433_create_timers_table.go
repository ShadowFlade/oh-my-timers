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
	db := PkgDb.Db{}
	dbUser := PkgDb.Timer{}
	timerDb := dbUser.Create()

	timersTableName := timerDb.TableName
	rows, err := db.Db.Query("SHOW TABLES LIKE %s;", timersTableName)

	if err != nil {
		log.Panic(err.Error())
	}

	if !rows.Next() {
		return nil
	}

	query := `
	CREATE TABLE timers (
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
	`

	res, err := db.Db.NamedExec(query, timersTableName)

	if err != nil {
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
	db := PkgDb.Db{}
	dbTimer := PkgDb.Timer{}
	timerDb := dbTimer.Create()

	timersTableName := timerDb.TableName
	fmt.Print(timersTableName, "table name")
	rows, err := db.Db.Query("SHOW TABLES LIKE %s;", timersTableName)

	if err != nil {
		log.Panic(err.Error())
	}

	if !rows.Next() {
		return nil
	}

	query := `
	drop table if exists %s
	`

	res, err := db.Db.NamedExec(query, timersTableName)

	if err != nil {
		tx.Rollback()
		log.Panic(err.Error())
	}

	rowsAffected, err := res.RowsAffected()

	if rowsAffected > 0 {
		tx.Commit()
	}

	return nil
}
