package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateCategoriesTable, downCreateCategoriesTable)
}

func upCreateCategoriesTable(ctx context.Context, tx *sql.Tx) error {
	createSectionTableQuery := `create table if not exists section
		(
			id int not null auto_increment primary key,
			name varchar(100) not null,
			color varchar(100) not null default 'red'
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`

	if _, err := tx.Exec(createSectionTableQuery); err != nil {
		return err
	}

	query := `create table if not exists timer_section
		(
			id int not null auto_increment primary key,
			user_id int not null,
			timer_id int not null,
			section_id int not null,
			color varchar(100) not null default 'red',
			foreign key (user_id) references users(id),
			foreign key (section_id) references section(id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`

	_, err := tx.Exec(query)

	return err
}

func downCreateCategoriesTable(ctx context.Context, tx *sql.Tx) error {
	query := fmt.Sprintf("drop table if exists %s;", "timer_section")
	_, err := tx.Exec(query)

	return err
}
