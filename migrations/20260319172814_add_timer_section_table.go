package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"shadowflade/timers/pkg/db"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateCategoriesTable, downCreateCategoriesTable)
}

func upCreateCategoriesTable(ctx context.Context, tx *sql.Tx) error {
	timerCategory := db.Category{}
	timerCategoryModel := timerCategory.New()
	var isTableAlreadyExists interface{}
	checkTableExistenceQuery := fmt.Sprintf("show tables like %s", timerCategoryModel.TableName)
	tx.QueryRow(checkTableExistenceQuery).Scan(&isTableAlreadyExists)
	log.Println(isTableAlreadyExists.(bool))

	if isTableAlreadyExists.(bool) {
		return fmt.Errorf("Table already exists")
	}

	query := fmt.Sprintf("create table %s (id int not null auto_increment primary key, user_id bigint not null  foreign key (user_id) references users(id), color varchar(100) not null default 'red', category_id bigint not null foreign key (category_id) references categories(id))", timerCategoryModel.TableName)

	_, err := tx.Exec(query)

	if err != nil {
		panic(err.Error())
	}

	return nil
}

func downCreateCategoriesTable(ctx context.Context, tx *sql.Tx) error {
	query := fmt.Sprintf("drop table %s")
	_, err := tx.Exec(query)

	if err != nil {
		panic(err.Error())
	}

	return nil
}
