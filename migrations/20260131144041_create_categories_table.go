package main

import (
	"context"
	"database/sql"
	"shadowflade/timers/pkg/db"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddTimerSectionTable, downAddTimerSectionTable)
}

func upAddTimerSectionTable(ctx context.Context, tx *sql.Tx) error {
	timerCategory := db.Category{}
	timerCategoryModel := timerCategory.New()
	var isTableAlreadyExists interface{}
	checkTableExistenceQuery := fmt.Sprintf("show tables like %s", timerCategoryModel.TableName)
	tx.QueryRow(checkTableExistenceQuery).Scan(&isTableAlreadyExists)
	log.Println(isTableAlreadyExists.(bool))

	if isTableAlreadyExists.(bool) {
		return fmt.Errorf("Table already exists")
	}

	query := fmt.Sprintf("create table %s (id int not null auto_increment primary key, name varchar(100) not null, color varchar(100) not null default 'red')", timerCategoryModel.TableName)
	return nil
}

func downAddTimerSectionTable(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
