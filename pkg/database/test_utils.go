package database

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

const (
	testDatabaseURL        = "postgres://rental:rental@localhost:5432/rental?sslmode=disable"
	databaseMigrationsPath = "file://../../db/migrations"
)

// setupTestDatabase sets up a test database with the migrations defined in the migrations folder, and returns a sqlx.DB instance for tests to use to change database state.
func setupTestDatabase() {
	db, err := sql.Open("postgres", testDatabaseURL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance(databaseMigrationsPath, "postgres", driver)
	if err != nil {
		panic(err)
	}
	if err := m.Up(); err != nil {
		panic(err)
	}
}

func teardownTestDatabase() {
	db, err := sql.Open("postgres", testDatabaseURL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance(databaseMigrationsPath, "postgres", driver)
	if err != nil {
		panic(err)
	}
	if err := m.Down(); err != nil {
		panic(err)
	}
}
