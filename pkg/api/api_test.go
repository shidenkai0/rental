package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/shidenkai0/rental/mock"
	"github.com/shidenkai0/rental/pkg/api/gen"
	"github.com/shidenkai0/rental/pkg/rental"
	"gopkg.in/guregu/null.v4"
)

func TestServerCreateCar(t *testing.T) {
	// Setup
	e := echo.New()
	s := &Server{CarCRUDService: mock.NewMockCarCRUDService()}

	createCar := gen.CreateUpdateCarRequest{Make: "Toyota", Model: "Corolla", Year: 2018}

	createCarJSON, _ := json.Marshal(createCar)

	req := httptest.NewRequest(http.MethodPost, "/car", bytes.NewBuffer(createCarJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resp := httptest.NewRecorder()
	ctx := e.NewContext(req, resp)
	ctx.SetPath("/car")

	// Test

	if err := s.CreateCar(ctx); err != nil {
		t.Errorf("got error %v, want nil", err)
	}
	if resp.Code != http.StatusCreated {
		t.Errorf("got %d status code, want %d", resp.Code, http.StatusCreated)
	}

	want := rental.Car{ID: 0, Make: createCar.Make, Model: createCar.Model, Year: createCar.Year}

	got, err := s.CarCRUDService.Get(0) // The ID can't be specified through the API, so it defaults to 0.
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestServerDeleteCar(t *testing.T) {
	// Setup

	e := echo.New()
	s := &Server{CarCRUDService: mock.NewMockCarCRUDService()}

	testCarID := 1

	car := rental.Car{ID: testCarID, Make: "Toyota", Model: "Corolla", Year: 2015}
	if _, err := s.CarCRUDService.Create(car); err != nil {
		t.Errorf("got error %v, want nil", err)
	}

	path := fmt.Sprintf("/car/%d", testCarID)
	req := httptest.NewRequest(http.MethodDelete, path, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resp := httptest.NewRecorder()
	ctx := e.NewContext(req, resp)
	ctx.SetPath(path)

	// Test

	if err := s.DeleteCar(ctx, int64(testCarID)); err != nil {
		t.Errorf("got error %v, want nil", err)
	}

	if resp.Code != http.StatusNoContent {
		t.Errorf("got %d status code, want %d", resp.Code, http.StatusNoContent)
	}

	_, err := s.CarCRUDService.Get(testCarID)
	if err != rental.ErrCarNotFound {
		t.Errorf("got error %v, want %v", err, rental.ErrCarNotFound)
	}
}

func TestServerGetCarById(t *testing.T) {
	// Setup

	e := echo.New()
	s := &Server{CarCRUDService: mock.NewMockCarCRUDService()}

	testCarID := 1

	car := rental.Car{ID: testCarID, Make: "Toyota", Model: "Corolla", Year: 2015}
	if _, err := s.CarCRUDService.Create(car); err != nil {
		t.Errorf("got error %v, want nil", err)
	}

	path := fmt.Sprintf("/car/%d", testCarID)
	req := httptest.NewRequest(http.MethodGet, path, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resp := httptest.NewRecorder()
	ctx := e.NewContext(req, resp)
	ctx.SetPath(path)

	// Test

	if err := s.GetCarById(ctx, int64(testCarID)); err != nil {
		t.Errorf("got error %v, want nil", err)
	}
	if resp.Code != http.StatusOK {
		t.Errorf("got %d status code, want %d", resp.Code, http.StatusOK)
	}

	// We add the trailing newline separately because it can't be specified in a multi-line string.
	wantBody := `{"id":1,"make":"Toyota","model":"Corolla","renter_id":0,"year":2015}` + "\n"
	got := resp.Body.String()

	if got != wantBody {
		t.Errorf("got %s, want %s", resp.Body.String(), wantBody)
	}
}

func TestServerUpdateCar(t *testing.T) {
	// Setup

	e := echo.New()
	s := &Server{CarCRUDService: mock.NewMockCarCRUDService()}

	testCarID := 1
	rentedToID := 2

	// set CustomerID to make sure CustomerID is kept as-is when updating
	car := rental.Car{ID: testCarID, Make: "Toyota", CustomerID: null.NewInt(int64(rentedToID), true), Model: "Corolla", Year: 2015}
	if _, err := s.CarCRUDService.Create(car); err != nil {
		t.Errorf("got error %v, want nil", err)
	}

	updateCar := gen.CreateUpdateCarRequest{Make: "Honda", Model: "Civic", Year: 2017}
	updateCarJSON, _ := json.Marshal(updateCar)

	path := fmt.Sprintf("/car/%d", testCarID)
	req := httptest.NewRequest(http.MethodPut, path, bytes.NewBuffer(updateCarJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resp := httptest.NewRecorder()
	ctx := e.NewContext(req, resp)
	ctx.SetPath(path)

	// Test

	if err := s.UpdateCar(ctx, int64(testCarID)); err != nil {
		t.Errorf("got error %v, want nil", err)
	}

	if resp.Code != http.StatusOK {
		t.Errorf("got %d status code, want %d", resp.Code, http.StatusOK)
	}

	wantBody := `{"id":1,"make":"Honda","model":"Civic","renter_id":2,"year":2017}` + "\n"
	if resp.Body.String() != wantBody {
		t.Errorf("got %s, want %s", resp.Body.String(), wantBody)
	}

	got, err := s.CarCRUDService.Get(testCarID)
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}

	want := rental.Car{ID: testCarID, Make: updateCar.Make, Model: updateCar.Model, CustomerID: null.NewInt(int64(rentedToID), true), Year: updateCar.Year}
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestServerRentCar(t *testing.T) {
	t.Run("rent an available car", func(t *testing.T) {
		// Setup

		e := echo.New()
		carCRUDService := mock.NewMockCarCRUDService()
		customerCRUDService := mock.NewMockCustomerCRUDService()
		s := &Server{CarCRUDService: carCRUDService, CustomerCRUDService: customerCRUDService}

		testCarID := 1

		testCustomerID := 1

		car := rental.Car{ID: testCarID, Make: "Toyota", Model: "Corolla", Year: 2015}
		if _, err := s.CarCRUDService.Create(car); err != nil {
			t.Errorf("got error %v, want nil", err)
		}
		customer := rental.Customer{ID: testCustomerID, Name: "John Doe"}
		if _, err := s.CustomerCRUDService.Create(customer); err != nil {
			t.Errorf("got error %v, want nil", err)
		}

		path := fmt.Sprintf("/car/%d/rent", testCarID)
		req := httptest.NewRequest(http.MethodGet, path, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		resp := httptest.NewRecorder()
		ctx := e.NewContext(req, resp)
		ctx.SetPath(path)

		// Test

		if err := s.RentCar(ctx, int64(testCarID), gen.RentCarParams{CustomerId: int64(testCustomerID)}); err != nil {
			t.Errorf("got error %v, want nil", err)
		}
		if resp.Code != http.StatusNoContent {
			t.Errorf("got %d status code, want %d", resp.Code, http.StatusNoContent)
		}

		rentedCar, err := s.CarCRUDService.Get(testCarID)
		if err != nil {
			t.Errorf("got error %v, want nil", err)
		}
		if rentedCar.RenterID() == 0 {
			t.Errorf("car %d not rented to customer %d", testCarID, testCustomerID)
		}
	})
	t.Run("rent a rented car", func(t *testing.T) {
		// Setup
		e := echo.New()
		carCRUDService := mock.NewMockCarCRUDService()
		customerCRUDService := mock.NewMockCustomerCRUDService()
		s := &Server{CarCRUDService: carCRUDService, CustomerCRUDService: customerCRUDService}

		testCarID := 1
		testCustomerID := 1

		car := rental.Car{ID: testCarID, Make: "Toyota", Model: "Corolla", Year: 2015, CustomerID: null.IntFrom(int64(testCustomerID))}
		if _, err := s.CarCRUDService.Create(car); err != nil {
			t.Errorf("got error %v, want nil", err)
		}
		customer := rental.Customer{ID: testCustomerID, Name: "John Doe"}
		if _, err := s.CustomerCRUDService.Create(customer); err != nil {
			t.Errorf("got error %v, want nil", err)
		}

		path := fmt.Sprintf("/car/%d/rent", testCarID)
		req := httptest.NewRequest(http.MethodGet, path, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		resp := httptest.NewRecorder()
		ctx := e.NewContext(req, resp)
		ctx.SetPath(path)

		// Test

		err := s.RentCar(ctx, int64(testCarID), gen.RentCarParams{CustomerId: int64(testCustomerID)})
		if err == nil {
			t.Errorf("got nil, want error")
		}

		he, _ := err.(*echo.HTTPError)

		got, want := he.Code, http.StatusForbidden
		if got != want {
			t.Errorf("got %d status code, want %d", got, want)
		}
	})
	t.Run("rent a car to non-existent customer", func(t *testing.T) {
		// Setup
		e := echo.New()
		carCRUDService := mock.NewMockCarCRUDService()
		customerCRUDService := mock.NewMockCustomerCRUDService()
		s := &Server{CarCRUDService: carCRUDService, CustomerCRUDService: customerCRUDService}

		testCarID := 1

		car := rental.Car{ID: testCarID, Make: "Toyota", Model: "Corolla", Year: 2015}
		if _, err := s.CarCRUDService.Create(car); err != nil {
			t.Errorf("got error %v, want nil", err)
		}

		path := fmt.Sprintf("/car/%d/rent", testCarID)
		req := httptest.NewRequest(http.MethodGet, path, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		resp := httptest.NewRecorder()
		ctx := e.NewContext(req, resp)
		ctx.SetPath(path)

		// Test

		err := s.RentCar(ctx, int64(testCarID), gen.RentCarParams{CustomerId: int64(1)})
		if err == nil {
			t.Errorf("got nil, want error")
		}

		he, _ := err.(*echo.HTTPError)

		got, want := he.Code, http.StatusBadRequest
		if got != want {
			t.Errorf("got %d status code, want %d", got, want)
		}
	})
}

func TestServerCreateCustomer(t *testing.T) {
	// Setup
	e := echo.New()
	s := &Server{CustomerCRUDService: mock.NewMockCustomerCRUDService()}

	createCustomer := gen.CreateUpdateCustomerRequest{Name: "John Doe"}
	createCustomerJSON, _ := json.Marshal(createCustomer)

	req := httptest.NewRequest(http.MethodPost, "/customer", bytes.NewBuffer(createCustomerJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resp := httptest.NewRecorder()
	ctx := e.NewContext(req, resp)
	ctx.SetPath("/customer")

	// Test

	if err := s.CreateCustomer(ctx); err != nil {
		t.Errorf("got error %v, want nil", err)
	}
	if resp.Code != http.StatusCreated {
		t.Errorf("got %d status code, want %d", resp.Code, http.StatusCreated)
	}

	want := rental.Customer{ID: 0, Name: createCustomer.Name}
	got, err := s.CustomerCRUDService.Get(0) // The ID can't be specified through the API, so it defaults to 0.
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestServerDeleteCustomer(t *testing.T) {
	// Setup

	e := echo.New()
	s := &Server{CustomerCRUDService: mock.NewMockCustomerCRUDService()}

	testCustomerID := 1

	customer := rental.Customer{ID: testCustomerID, Name: "John Doe"}
	if _, err := s.CustomerCRUDService.Create(customer); err != nil {
		t.Errorf("got error %v, want nil", err)
	}

	path := fmt.Sprintf("/customer/%d", customer.ID)
	req := httptest.NewRequest(http.MethodDelete, path, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resp := httptest.NewRecorder()
	ctx := e.NewContext(req, resp)
	ctx.SetPath(path)

	// Test

	if err := s.DeleteCustomer(ctx, int64(testCustomerID)); err != nil {
		t.Errorf("got error %v, want nil", err)
	}
	if resp.Code != http.StatusNoContent {
		t.Errorf("got %d status code, want %d", resp.Code, http.StatusNoContent)
	}

	_, err := s.CustomerCRUDService.Get(testCustomerID)
	if err != rental.ErrCustomerNotFound {
		t.Errorf("got error %v, want %v", err, rental.ErrCustomerNotFound)
	}
}

func TestServerGetCustomerById(t *testing.T) {
	// Setup
	e := echo.New()
	s := &Server{CustomerCRUDService: mock.NewMockCustomerCRUDService()}

	testCustomerID := 1

	customer := rental.Customer{ID: testCustomerID, Name: "John Doe"}
	if _, err := s.CustomerCRUDService.Create(customer); err != nil {
		t.Errorf("got error %v, want nil", err)
	}

	path := fmt.Sprintf("/customer/%d", testCustomerID)
	req := httptest.NewRequest(http.MethodGet, path, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resp := httptest.NewRecorder()
	ctx := e.NewContext(req, resp)
	ctx.SetPath(path)

	// Test

	if err := s.GetCustomerById(ctx, int64(testCustomerID)); err != nil {
		t.Errorf("got error %v, want nil", err)
	}
	if resp.Code != http.StatusOK {
		t.Errorf("got %d status code, want %d", resp.Code, http.StatusOK)
	}

	got, err := s.CustomerCRUDService.Get(testCustomerID)
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}
	if got != customer {
		t.Errorf("got %v, want %v", got, customer)
	}
}

func TestServerUpdateCustomer(t *testing.T) {
	// Setup
	e := echo.New()
	s := &Server{CustomerCRUDService: mock.NewMockCustomerCRUDService()}

	testCustomerID := 1

	customer := rental.Customer{ID: testCustomerID, Name: "John Doe"}
	if _, err := s.CustomerCRUDService.Create(customer); err != nil {
		t.Errorf("got error %v, want nil", err)
	}

	updateCustomer := gen.CreateUpdateCustomerRequest{Name: "Jane Doe"}
	updateCustomerJSON, _ := json.Marshal(updateCustomer)

	path := fmt.Sprintf("/customer/%d", testCustomerID)
	req := httptest.NewRequest(http.MethodPut, path, bytes.NewBuffer(updateCustomerJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resp := httptest.NewRecorder()
	ctx := e.NewContext(req, resp)
	ctx.SetPath(path)

	// Test

	if err := s.UpdateCustomer(ctx, int64(testCustomerID)); err != nil {
		t.Errorf("got error %v, want nil", err)
	}

	if resp.Code != http.StatusOK {
		t.Errorf("got %d status code, want %d", resp.Code, http.StatusOK)
	}

	got, err := s.CustomerCRUDService.Get(testCustomerID)
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}

	want := rental.Customer{ID: testCustomerID, Name: updateCustomer.Name}
	if got != want {
		t.Errorf("got %v, want %v", got, customer)
	}
}
