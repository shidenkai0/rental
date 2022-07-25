// Package rental provides the domain model of the rental service.
package rental

import (
	"fmt"

	"gopkg.in/guregu/null.v4"
)

// Car represents a car.
type Car struct {
	ID         int      `json:"id" db:"id"`
	CustomerID null.Int `json:"renter_id" db:"customer_id"` // Make private (can't be modified through the API)
	Make       string   `json:"make" db:"make"`
	Model      string   `json:"model" db:"model"`
	Year       int      `json:"year" db:"year"`
}

// RenterID returns the ID of the customer who has rented the car, 0 if the car is not rented.
func (car *Car) RenterID() int {
	return int(car.CustomerID.ValueOrZero())
}

// Rented returns true if the car is rented.
func (car *Car) Rented() bool {
	return car.RenterID() != 0
}

// Rent rents the car to a customer.
func (car *Car) Rent(customerID int) error {
	if car.Rented() {
		return ErrCarAlreadyRented
	}
	car.CustomerID = null.IntFrom(int64(customerID))
	return nil
}

// Return returns the car and makes it available for rental.
func (car *Car) Return() error {
	if !car.Rented() {
		return ErrCarNotRented
	}
	car.CustomerID = null.IntFromPtr(nil)
	return nil
}

type CarCRUDService interface {
	Create(car Car) (int, error)
	Get(id int) (Car, error)
	Update(car Car) error
	Delete(carID int) error
}

var (
	ErrCarNotFound      = fmt.Errorf("Car not found")
	ErrCarNotRented     = fmt.Errorf("Car not rented")
	ErrCarAlreadyRented = fmt.Errorf("Car already rented")
	ErrCarAlreadyExists = fmt.Errorf("Car already exists")
)
