// Package mock provides mock implementations of stateful services.
package mock

import (
	"github.com/shidenkai0/rental/pkg/rental"
)

type MockCustomerCRUDService struct {
	customers map[int]*rental.Customer
}

// Create creates a customer in the Mock state and returns the id.
func (m *MockCustomerCRUDService) Create(customer rental.Customer) (id int, err error) {
	if _, ok := m.customers[customer.ID]; ok {
		return 0, rental.ErrCustomerAlreadyExists
	}
	m.customers[customer.ID] = &customer
	return customer.ID, nil
}

// Get fetches a customer from the Mock state.
func (m *MockCustomerCRUDService) Get(id int) (rental.Customer, error) {
	customer, ok := m.customers[id]
	if !ok {
		return rental.Customer{}, rental.ErrCustomerNotFound
	}
	return *customer, nil
}

// Update updates a customer in the Mock state.
func (m *MockCustomerCRUDService) Update(customer rental.Customer) error {
	m.customers[customer.ID] = &customer
	return nil
}

// Delete deletes a customer from the Mock state.
func (m *MockCustomerCRUDService) Delete(customerID int) error {
	delete(m.customers, customerID)
	return nil
}

// NewMockCustomerCRUDService returns a new MockCustomerCRUDService.
func NewMockCustomerCRUDService() *MockCustomerCRUDService {
	return &MockCustomerCRUDService{customers: map[int]*rental.Customer{}}
}
