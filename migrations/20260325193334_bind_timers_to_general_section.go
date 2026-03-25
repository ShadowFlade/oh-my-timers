package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upBindTimersToGeneralSection, downBindTimersToGeneralSection)
}

func upBindTimersToGeneralSection(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
        INSERT INTO section (name, color)
        VALUES ('general', 'red')
        ON DUPLICATE KEY UPDATE name = name
    `)
	if err != nil {
		return fmt.Errorf("inserting general section: %w", err)
	}

	// _, err = tx.ExecContext(ctx, `
	//        INSERT INTO timer_section (user_id, timer_id, section_id, color)
	//        SELECT
	//            t.user_id,
	//            t.id,
	//            s.id,
	//            COALESCE(t.color, 'red')
	//        FROM timers t
	//        CROSS JOIN section s
	//        WHERE t.user_id IS NOT NULL
	//        AND s.name = 'general'
	//        AND NOT EXISTS (
	//            SELECT 1 FROM timer_section ts WHERE ts.timer_id = t.id
	//        )
	//    `)
	_, err = tx.ExecContext(ctx, `
        INSERT INTO timer_section (user_id, timer_id, section_id, color)
        SELECT
            t.user_id,
            t.id,
            (SELECT id FROM section WHERE name = 'general' LIMIT 1),
            COALESCE(t.color, 'red')
        FROM timers t
        WHERE t.user_id IS NOT NULL
        AND NOT EXISTS (
            SELECT 1 FROM timer_section ts WHERE ts.timer_id = t.id
        )
    `)
	if err != nil {
		return fmt.Errorf("binding timers to section: %w", err)
	}

	return nil
}
func downBindTimersToGeneralSection(ctx context.Context, tx *sql.Tx) error {
	//query := fmt.Sprintf("drop table if exists %s;", "timer_section")
	// _, err := tx.ExecContext(ctx, query)
	//
	// return err
	return nil
}
