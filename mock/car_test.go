package mock

import (
	"fmt"
	"testing"

	"github.com/shidenkai0/rental/pkg/rental"
)

func TestMockCarCRUDService_Get(t *testing.T) {
	t.Run("get car", func(t *testing.T) {
		testCarID := 1
		car := rental.Car{ID: testCarID, Make: "Toyota", Model: "Corolla", Year: 2015}
		mockCarCRUDService := &MockCarCRUDService{cars: map[int]*rental.Car{1: &car}}
		got, err := mockCarCRUDService.Get(testCarID)
		if err != nil {
			t.Errorf("got error %v, want nil", err)
		}
		if got != car {
			t.Errorf("got %v, want %v", got, car)
		}
	})
	t.Run("get non-existent car", func(t *testing.T) {
		mockCarCRUDService := &MockCarCRUDService{}
		_, err := mockCarCRUDService.Get(1)
		if err != rental.ErrCarNotFound {
			t.Errorf(fmt.Sprintf("got error %v, want %v", err, rental.ErrCarNotFound))
		}
	})
}

func TestMockCarCRUDService_Create(t *testing.T) {
	testCarID := 1
	t.Run("create car", func(t *testing.T) {
		car := rental.Car{ID: testCarID, Make: "Toyota", Model: "Corolla", Year: 2015}
		mockCarCRUDService := NewMockCarCRUDService()
		_, err := mockCarCRUDService.Create(car)
		if err != nil {
			t.Errorf("got error %v, want nil", err)
		}
		got, err := mockCarCRUDService.Get(testCarID)
		if err != nil {
			t.Errorf("got error %v, want nil", err)
		}
		if got != car {
			t.Errorf("got %v, want %v", got, car)
		}
	})
	t.Run("create duplicate car", func(t *testing.T) {
		car := rental.Car{ID: testCarID, Make: "Toyota", Model: "Corolla", Year: 2015}
		mockCarCRUDService := NewMockCarCRUDService()
		_, err := mockCarCRUDService.Create(car)
		if err != nil {
			t.Errorf("got error %v, want nil", err)
		}
		_, err = mockCarCRUDService.Create(car)
		if err != rental.ErrCarAlreadyExists {
			t.Errorf(fmt.Sprintf("got error %v, want %v", err, rental.ErrCarAlreadyExists))
		}
	})
}

func TestMockCarCRUDService_Update(t *testing.T) {
	testCarID := 1
	car := rental.Car{ID: testCarID, Make: "Toyota", Model: "Corolla", Year: 2015}
	mockCarCRUDService := NewMockCarCRUDService()
	_, err := mockCarCRUDService.Create(car)
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}
	car.Make = "Honda"
	err = mockCarCRUDService.Update(car)
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}
	got, err := mockCarCRUDService.Get(testCarID)
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}
	if got != car {
		t.Errorf("got %v, want %v", got, car)
	}
}

func TestMockCarCRUDService_Delete(t *testing.T) {
	testCarID := 1
	car := rental.Car{ID: testCarID, Make: "Toyota", Model: "Corolla", Year: 2015}
	mockCarCRUDService := NewMockCarCRUDService()
	_, err := mockCarCRUDService.Create(car)
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}
	err = mockCarCRUDService.Delete(testCarID)
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}
	_, err = mockCarCRUDService.Get(testCarID)
	if err != rental.ErrCarNotFound {
		t.Errorf("got error %v, want %v", err, rental.ErrCarNotFound)
	}
}
