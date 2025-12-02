package migrations

import (
	"context"
	"database/sql"
	"github.com/pressly/goose/v3"
	"shadowflade/timers/pkg/db/"
)

func init() {
	goose.AddMigrationContext(upCreateUsersTable, downCreateUsersTable)
}

func upCreateUsersTable(ctx context.Context, tx *sql.Tx) error {
	db := db.Db{}
	db = db.Connect
	db.
	// This code is executed when the migration is applied.
	return nil
}

func downCreateUsersTable(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
