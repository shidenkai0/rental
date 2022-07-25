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
	err = m.Up()
	if err == migrate.ErrNoChange {
		return
	}
	if err != nil {
		panic(err)
	}
}

// teardownTestDatabase tears down the test database. It should be deferred after setupTestDatabase in
// test cases as it allows to cleanly reset the database state in between tests to provide isolation.
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
