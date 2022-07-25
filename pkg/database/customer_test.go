package database

import (
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/shidenkai0/rental/pkg/rental"
)

func TestDatabaseCustomerCRUDServiceGet(t *testing.T) {
	setupTestDatabase()
	defer teardownTestDatabase()
	db := sqlx.MustConnect("postgres", testDatabaseURL)
	t.Run("get a customer", func(t *testing.T) {
		customerCRUDService := NewDatabaseCustomerCRUDService(db)
		customer := rental.Customer{ID: 1, Name: "John Doe"}
		_, err := customerCRUDService.Create(customer)
		if err != nil {
			t.Errorf("got error %v, want nil", err)
		}
		got, err := customerCRUDService.Get(1)
		if err != nil {
			t.Errorf("got error %v, want nil", err)
		}
		if got != customer {
			t.Errorf("got %v, want %v", got, customer)
		}
	})
	t.Run("get a non-existent customer", func(t *testing.T) {
		customerCRUDService := NewDatabaseCustomerCRUDService(db)
		_, err := customerCRUDService.Get(100)
		if err != rental.ErrCustomerNotFound {
			t.Errorf(fmt.Sprintf("got error %v, want %v", err, rental.ErrCustomerNotFound))
		}
	})
}

func TestDatabaseCustomerCRUDServiceCreate(t *testing.T) {
	setupTestDatabase()
	defer teardownTestDatabase()

	db := sqlx.MustConnect("postgres", testDatabaseURL)
	customerCRUDService := NewDatabaseCustomerCRUDService(db)
	customer := rental.Customer{ID: 1, Name: "John Doe"}
	_, err := customerCRUDService.Create(customer)
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}

	got, err := customerCRUDService.Get(1)
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}
	if got != customer {
		t.Errorf("got %v, want %v", got, customer)
	}
}

func TestDatabaseCustomerCRUDServiceUpdate(t *testing.T) {
	setupTestDatabase()
	defer teardownTestDatabase()

	db := sqlx.MustConnect("postgres", testDatabaseURL)
	customerCRUDService := NewDatabaseCustomerCRUDService(db)
	customer := rental.Customer{ID: 1, Name: "John Doe"}
	_, err := customerCRUDService.Create(customer)
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}

	customer.Name = "Jane Doe"
	err = customerCRUDService.Update(customer)
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}

	got, err := customerCRUDService.Get(1)
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}
	if got != customer {
		t.Errorf("got %v, want %v", got, customer)
	}
}

func TestDatabaseCustomerCRUDServiceDelete(t *testing.T) {
	setupTestDatabase()
	defer teardownTestDatabase()

	db := sqlx.MustConnect("postgres", testDatabaseURL)
	customerCRUDService := NewDatabaseCustomerCRUDService(db)
	customer := rental.Customer{ID: 1, Name: "John Doe"}
	_, err := customerCRUDService.Create(customer)
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}

	err = customerCRUDService.Delete(1)
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}

	_, err = customerCRUDService.Get(1)
	if err != rental.ErrCustomerNotFound {
		t.Errorf(fmt.Sprintf("got error %v, want %v", err, rental.ErrCustomerNotFound))
	}
}
