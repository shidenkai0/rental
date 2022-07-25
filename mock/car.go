// Package mock provides mock implementations of stateful services.
package mock

import "github.com/shidenkai0/rental/pkg/rental"

type MockCarCRUDService struct {
	cars map[int]*rental.Car
}

// Create creates a car in the Mock state and returns the id.
func (m *MockCarCRUDService) Create(car rental.Car) (id int, err error) {
	if _, ok := m.cars[car.ID]; ok {
		return 0, rental.ErrCarAlreadyExists
	}
	m.cars[car.ID] = &car
	return car.ID, nil
}

// Get fetches a car from the Mock state.
func (m *MockCarCRUDService) Get(id int) (rental.Car, error) {
	car, ok := m.cars[id]
	if !ok {
		return rental.Car{}, rental.ErrCarNotFound
	}
	return *car, nil
}

// Update updates a car in the Mock state.
func (m *MockCarCRUDService) Update(car rental.Car) error {
	m.cars[car.ID] = &car
	return nil
}

// Delete deletes a car from the Mock state.
func (m *MockCarCRUDService) Delete(carID int) error {
	delete(m.cars, carID)
	return nil
}

// NewMockCarCRUDService returns a new MockCarCRUDService.
func NewMockCarCRUDService() *MockCarCRUDService {
	return &MockCarCRUDService{cars: map[int]*rental.Car{}}
}
