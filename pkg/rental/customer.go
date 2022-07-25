// Package rental provides the domain model of the rental service.
package rental

import "fmt"

// Customer represents a customer.
type Customer struct {
	ID   int    `json:"id" sql:"id"`
	Name string `json:"name" sql:"name"`
}

type CustomerCRUDService interface {
	Create(customer Customer) (int, error)
	Get(id int) (Customer, error)
	Update(customer Customer) error
	Delete(customerID int) error
}

var (
	ErrCustomerNotFound      = fmt.Errorf("Customer not found")
	ErrCustomerAlreadyExists = fmt.Errorf("Customer already exists")
)
