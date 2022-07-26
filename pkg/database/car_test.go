package database

import (
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/shidenkai0/rental/pkg/rental"
	"gopkg.in/guregu/null.v4"
)

func TestDatabaseCarCRUDServiceGet(t *testing.T) {
	setupTestDatabase()
	defer teardownTestDatabase()
	db := sqlx.MustConnect("postgres", testDatabaseURL)
	t.Run("get a car", func(t *testing.T) {
		carCRUDService := NewDatabaseCarCRUDService(db)
		car := rental.Car{ID: 1, Make: "Toyota", Model: "Corolla", Year: 2015}
		_, err := carCRUDService.Create(car)
		if err != nil {
			t.Errorf("got error %v, want nil", err)
		}
		got, err := carCRUDService.Get(1)
		if err != nil {
			t.Errorf("got error %v, want nil", err)
		}
		if got != car {
			t.Errorf("got %v, want %v", got, car)
		}
	})
	t.Run("get a non-existent car", func(t *testing.T) {
		carCRUDService := NewDatabaseCarCRUDService(db)
		_, err := carCRUDService.Get(100)
		if err != rental.ErrCarNotFound {
			t.Errorf(fmt.Sprintf("got error %v, want %v", err, rental.ErrCarNotFound))
		}
	})
}

func TestDatabaseCarCRUDServiceCreate(t *testing.T) {
	setupTestDatabase()
	defer teardownTestDatabase()

	db := sqlx.MustConnect("postgres", testDatabaseURL)
	carCRUDService := NewDatabaseCarCRUDService(db)
	car := rental.Car{ID: 1, Make: "Toyota", Model: "Corolla", Year: 2015}
	_, err := carCRUDService.Create(car)
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}

	got, err := carCRUDService.Get(1)
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}
	if got != car {
		t.Errorf("got %v, want %v", got, car)
	}
}

func TestDatabaseCarCRUDServiceUpdate(t *testing.T) {
	setupTestDatabase()
	defer teardownTestDatabase()

	db := sqlx.MustConnect("postgres", testDatabaseURL)

	// Create test customers

	customerCRUDService := NewDatabaseCustomerCRUDService(db)
	testCustomer := rental.Customer{ID: 1, Name: "John Doe"}

	_, err := customerCRUDService.Create(testCustomer)
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}

	carCRUDService := NewDatabaseCarCRUDService(db)
	car := rental.Car{ID: 1, Make: "Toyota", Model: "Corolla", Year: 2015}
	_, err = carCRUDService.Create(car)
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}

	car = rental.Car{ID: 1, Make: "Ford", CustomerID: null.IntFrom(int64(testCustomer.ID)), Model: "Fiesta", Year: 2016}
	err = carCRUDService.Update(car)
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}

	got, err := carCRUDService.Get(1)
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}
	if got != car {
		t.Errorf("got %v, want %v", got, car)
	}
}

func TestDatabaseCarCRUDServiceDelete(t *testing.T) {
	setupTestDatabase()
	defer teardownTestDatabase()

	db := sqlx.MustConnect("postgres", testDatabaseURL)
	carCRUDService := NewDatabaseCarCRUDService(db)
	car := rental.Car{ID: 1, Make: "Toyota", Model: "Corolla", Year: 2015}
	_, err := carCRUDService.Create(car)
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}

	err = carCRUDService.Delete(1)
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}

	car, err = carCRUDService.Get(1)
	if (err != rental.ErrCarNotFound || car != rental.Car{}) {
		t.Errorf(fmt.Sprintf("got error %v, want %v", err, rental.ErrCarNotFound))
	}
}
