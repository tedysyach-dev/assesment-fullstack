package main

import (
	"context"
	"fmt"
	"os"
	"wms/core/config"
	"wms/migrations"

	"github.com/uptrace/bun/migrate"
)

func main() {
	v := config.NewViper()
	logger := config.NewLogger(v)
	db := config.NewBun(v, logger)

	migrator := migrate.NewMigrator(db, migrations.Migrations)

	cmd := ""
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	ctx := context.Background()

	switch cmd {
	case "init":
		if err := migrator.Init(ctx); err != nil {
			logger.Fatalf("init failed: %v", err)
		}
		fmt.Println("Migration table created")

	case "up":
		if err := migrator.Lock(ctx); err != nil {
			logger.Fatalf("lock failed: %v", err)
		}
		defer migrator.Unlock(ctx)

		group, err := migrator.Migrate(ctx)
		if err != nil {
			logger.Fatalf("migrate failed: %v", err)
		}
		if group.IsZero() {
			fmt.Println("No new migrations to run")
			return
		}
		fmt.Printf("Migrated: %s\n", group)

	case "down":
		if err := migrator.Lock(ctx); err != nil {
			logger.Fatalf("lock failed: %v", err)
		}
		defer migrator.Unlock(ctx)

		group, err := migrator.Rollback(ctx)
		if err != nil {
			logger.Fatalf("rollback failed: %v", err)
		}
		if group.IsZero() {
			fmt.Println("Nothing to rollback")
			return
		}
		fmt.Printf("Rolled back: %s\n", group)

	case "status":
		ms, err := migrator.MigrationsWithStatus(ctx)
		if err != nil {
			logger.Fatalf("status failed: %v", err)
		}
		for _, m := range ms {
			fmt.Printf("%-40s %s\n", m.Name, m.MigratedAt)
		}

	default:
		fmt.Println("Usage: migrate [init|up|down|status]")
	}
}
