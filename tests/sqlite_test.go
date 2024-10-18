package tests

import (
	flyway "com.goldstar/goflyway/goflyway"
	"com.goldstar/goflyway/goflyway/command"
	"com.goldstar/goflyway/goflyway/database"
	"com.goldstar/goflyway/goflyway/logger"
	"fmt"
	"testing"
)

func TestSqliteMigrate(t *testing.T) {
	db, err := ConnSqlite()
	if err != nil {
		t.Fatal(err)
	}
	f, err := flyway.Open(database.T_SQLITE, db, &flyway.Config{
		Locations:         []string{"db_migration/sqlite"},
		Logger:            logger.Default.LogMode(logger.Info),
		EnablePlaceholder: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	f.Callbacks().Registry("before:migrate", "test", func(context *command.Context) {
		fmt.Println("test before migrate")
	})
	err = f.Migrate()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("sqlite migrate success")
}
