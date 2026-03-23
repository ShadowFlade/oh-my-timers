package main

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddTimerSectionTable, downAddTimerSectionTable)
}

func upAddTimerSectionTable(ctx context.Context, tx *sql.Tx) error {

	query := `
		create table section (
			id int not null auto_increment primary key,
			name varchar(100) not null,
			color varchar(100) not null default 'red'
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
	_, err := tx.Exec(query)

	if err != nil {
		return err
	}

	return nil
}

func downAddTimerSectionTable(ctx context.Context, tx *sql.Tx) error {

	query := "drop table if exists section"
	_, err := tx.Exec(query)

	return err
}
