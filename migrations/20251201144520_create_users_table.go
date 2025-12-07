package main

import (
	"context"
	"database/sql"
	"log"
	PkgDb "shadowflade/timers/pkg/db"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateUsersTable, downCreateUsersTable)
}

func upCreateUsersTable(ctx context.Context, tx *sql.Tx) error {
	db := PkgDb.Db{}
	db.Connect()
	dbUser := PkgDb.User{}
	userDB := dbUser.Create()

	userTable := userDB.TableName
	rows, err := db.Db.Query("SHOW TABLES LIKE %s;", userTable)

	if err != nil {
		log.Panic(err.Error())
	}

	if !rows.Next() {
		return nil
	}

	query := `
	CREATE TABLE users (
		id INT NOT NULL AUTO_INCREMENT,
		uuid CHAR(36) NOT NULL DEFAULT (UUID()),
		name VARCHAR(255) NULL,
		password TEXT NOT NULL,
		date_inserted DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		date_modified DATETIME NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		
		PRIMARY KEY (id),
		UNIQUE KEY uq_uuid (uuid)
	)
	`

	res, err := db.Db.NamedExec(query, userTable)

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

func downCreateUsersTable(ctx context.Context, tx *sql.Tx) error {
	db := PkgDb.Db{}
	db.Connect()
	dbUser := PkgDb.User{}
	userDB := dbUser.Create()

	userTable := userDB.TableName
	rows, err := db.Db.Query("SHOW TABLES LIKE %s;", userTable)

	if err != nil {
		log.Panic(err.Error())
	}

	if !rows.Next() {
		return nil
	}

	query := `
	drop table if exists %s
	`

	res, err := db.Db.NamedExec(query, userTable)

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
