// Package database implements the database layer of the rental service.
package database

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/shidenkai0/rental/pkg/rental"
)

// DatabaseCarCRUDService is a concrete implementation of the CarCRUDService
// interface using Postgres as a backend.
type DatabaseCarCRUDService struct {
	db *sqlx.DB
}

// Create creates a car in the database, returns id.
func (s *DatabaseCarCRUDService) Create(car rental.Car) (id int, err error) {
	insertStatement := "INSERT INTO cars (make, model, year) VALUES (:make, :model, :year) RETURNING id"
	rows, err := sqlx.NamedQuery(s.db, insertStatement, car)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return 0, err
		}
	}
	return id, err
}

// Get fetches a car from the database.
func (s *DatabaseCarCRUDService) Get(id int) (rental.Car, error) {
	var car rental.Car
	err := sqlx.Get(s.db, &car, "SELECT * FROM cars WHERE id = $1 LIMIT 1", id)
	if err == sql.ErrNoRows {
		return rental.Car{}, rental.ErrCarNotFound
	}
	return car, err
}

// Update updates a car in the database.
func (s *DatabaseCarCRUDService) Update(car rental.Car) error {
	updateStatement := "UPDATE cars SET make = :make, model = :model, year = :year WHERE id = :id"
	_, err := sqlx.NamedExec(s.db, updateStatement, car)
	return err
}

// Delete deletes a car from the database.
func (s *DatabaseCarCRUDService) Delete(carID int) error {
	_, err := s.db.Exec("DELETE FROM cars WHERE id = $1", carID)
	return err
}

// NewDatabaseCarCRUDService returns a new DatabaseCarCRUDService with the provided database as SQL backend.
func NewDatabaseCarCRUDService(db *sqlx.DB) *DatabaseCarCRUDService {
	return &DatabaseCarCRUDService{db: db}
}
