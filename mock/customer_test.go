package mock

import (
	"fmt"
	"testing"

	"github.com/shidenkai0/rental/pkg/rental"
)

func TestMockCustomerCRUDService_Get(t *testing.T) {
	t.Run("get customer", func(t *testing.T) {
		testCustomerID := 1
		customer := rental.Customer{ID: testCustomerID, Name: "John Doe"}
		mockCustomerCRUDService := &MockCustomerCRUDService{customers: map[int]*rental.Customer{1: &customer}}
		got, err := mockCustomerCRUDService.Get(testCustomerID)
		if err != nil {
			t.Errorf("got error %v, want nil", err)
		}
		if got != customer {
			t.Errorf("got %v, want %v", got, customer)
		}
	})
	t.Run("get non-existent customer", func(t *testing.T) {
		testCustomerID := 1
		mockCustomerCRUDService := &MockCustomerCRUDService{}
		want := rental.ErrCustomerNotFound
		_, err := mockCustomerCRUDService.Get(testCustomerID)
		if err != want {
			t.Errorf(fmt.Sprintf("got error %v, want %v", err, want))
		}
	})
}

func TestMockCustomerCRUDService_Create(t *testing.T) {
	testCustomerID := 1
	t.Run("create customer", func(t *testing.T) {
		customer := rental.Customer{ID: testCustomerID, Name: "John Doe"}
		mockCustomerCRUDService := NewMockCustomerCRUDService()
		_, err := mockCustomerCRUDService.Create(customer)
		if err != nil {
			t.Errorf("got error %v, want nil", err)
		}
		got, err := mockCustomerCRUDService.Get(testCustomerID)
		if err != nil {
			t.Errorf("got error %v, want nil", err)
		}
		if got != customer {
			t.Errorf("got %v, want %v", got, customer)
		}
	})
	t.Run("create duplicate customer", func(t *testing.T) {
		customer := rental.Customer{ID: testCustomerID, Name: "John Doe"}
		mockCustomerCRUDService := NewMockCustomerCRUDService()
		_, err := mockCustomerCRUDService.Create(customer)
		if err != nil {
			t.Errorf("got error %v, want nil", err)
		}
		_, err = mockCustomerCRUDService.Create(customer)
		if err != rental.ErrCustomerAlreadyExists {
			t.Errorf(fmt.Sprintf("got error %v, want %v", err, rental.ErrCustomerAlreadyExists))
		}
	})
}

func TestMockCustomerCRUDService_Update(t *testing.T) {

	testCustomerID := 1
	customer := rental.Customer{ID: testCustomerID, Name: "John Doe"}
	mockCustomerCRUDService := NewMockCustomerCRUDService()
	_, err := mockCustomerCRUDService.Create(customer)
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}
	err = mockCustomerCRUDService.Update(customer)
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}
	got, err := mockCustomerCRUDService.Get(testCustomerID)
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}
	if got != customer {
		t.Errorf("got %v, want %v", got, customer)
	}

}

func TestMockCustomerCRUDService_Delete(t *testing.T) {
	testCustomerID := 1
	customer := rental.Customer{ID: testCustomerID, Name: "John Doe"}
	mockCustomerCRUDService := NewMockCustomerCRUDService()
	_, err := mockCustomerCRUDService.Create(customer)
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}
	err = mockCustomerCRUDService.Delete(testCustomerID)
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}
	want := rental.ErrCustomerNotFound
	_, err = mockCustomerCRUDService.Get(testCustomerID)
	if err != want {
		t.Errorf("got error %v, want %v", err, rental.ErrCustomerNotFound)
	}
}
