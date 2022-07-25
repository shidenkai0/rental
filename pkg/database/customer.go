// Package database implements the database layer of the rental service.
package database

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/shidenkai0/rental/pkg/rental"
)

// DatabaseCustomerCRUDService is a concrete implementation of the CustomerCRUDService
// interface using Postgres as a backend.
type DatabaseCustomerCRUDService struct {
	db *sqlx.DB
}

// Create creates a customer in the database, returns id.
func (s *DatabaseCustomerCRUDService) Create(customer rental.Customer) (id int, err error) {
	insertStatement := "INSERT INTO customers (name) VALUES (:name) RETURNING id"
	rows, err := sqlx.NamedQuery(s.db, insertStatement, customer)
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

// Get fetches a customer from the database.
func (s *DatabaseCustomerCRUDService) Get(id int) (rental.Customer, error) {
	var customer rental.Customer
	err := sqlx.Get(s.db, &customer, "SELECT * FROM customers WHERE id = $1 LIMIT 1", id)
	if err == sql.ErrNoRows {
		return rental.Customer{}, rental.ErrCustomerNotFound
	}
	return customer, err
}

// Update updates a customer in the database.
func (s *DatabaseCustomerCRUDService) Update(customer rental.Customer) error {
	updateStatement := "UPDATE customers SET name = :name WHERE id = :id"
	_, err := sqlx.NamedExec(s.db, updateStatement, customer)
	return err
}

// Delete deletes a customer from the database.
func (s *DatabaseCustomerCRUDService) Delete(customerID int) error {
	_, err := s.db.Exec("DELETE FROM customers WHERE id = $1", customerID)
	return err
}

// NewDatabaseCustomerCRUDService returns a new DatabaseCustomerCRUDService with the provided database as SQL backend.
func NewDatabaseCustomerCRUDService(db *sqlx.DB) *DatabaseCustomerCRUDService {
	return &DatabaseCustomerCRUDService{db: db}
}
