package rental

import "testing"

func TestCarRent(t *testing.T) {
	t.Run("rent an available car", func(t *testing.T) {
		car := Car{ID: 1, Make: "Toyota", Model: "Corolla", Year: 2015}
		customer := Customer{ID: 1, Name: "John Doe"}
		err := car.Rent(customer.ID)
		if err != nil {
			t.Errorf("got error %v, want nil", err)
		}
		if car.RenterID() != customer.ID {
			t.Errorf("got %v, want %v", car.RenterID(), customer)
		}
	})
	t.Run("rent a rented car", func(t *testing.T) {
		car := Car{ID: 1, Make: "Toyota", Model: "Corolla", Year: 2015}
		customer := Customer{ID: 1, Name: "John Doe"}
		car.Rent(customer.ID)
		err := car.Rent(customer.ID)
		if err != ErrCarAlreadyRented {
			t.Errorf("got error %v, want %v", err, ErrCarAlreadyRented)
		}
	})
}

func TestCarReturn(t *testing.T) {
	t.Run("return a rented car", func(t *testing.T) {
		car := Car{ID: 1, Make: "Toyota", Model: "Corolla", Year: 2015}
		customer := Customer{ID: 1, Name: "John Doe"}
		car.Rent(customer.ID)
		err := car.Return()
		if err != nil {
			t.Errorf("got error %v, want nil", err)
		}
		if car.RenterID() != 0 {
			t.Errorf("got %v, want %v", car.RenterID(), 0)
		}
	})
	t.Run("return a not rented car", func(t *testing.T) {
		car := Car{ID: 1, Make: "Toyota", Model: "Corolla", Year: 2015}
		err := car.Return()
		if err != ErrCarNotRented {
			t.Errorf("got error %v, want %v", err, ErrCarNotRented)
		}
	})
}

func TestCarRented(t *testing.T) {
	t.Run("check if a car is rented", func(t *testing.T) {
		car := Car{ID: 1, Make: "Toyota", Model: "Corolla", Year: 2015}
		customer := Customer{ID: 1, Name: "John Doe"}
		car.Rent(customer.ID)
		if !car.Rented() {
			t.Errorf("got %v, want %v", car.Rented(), true)
		}
	})
	t.Run("check if a car is not rented", func(t *testing.T) {
		car := Car{ID: 1, Make: "Toyota", Model: "Corolla", Year: 2015}
		if car.Rented() {
			t.Errorf("got %v, want %v", car.Rented(), false)
		}
	})
}
