// Package api implements the API of the rental service as
// defined in the API specification at api/rental-v1.0.yml
package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shidenkai0/rental/pkg/api/gen"
	"github.com/shidenkai0/rental/pkg/rental"
)

func NewServer(carCRUDService rental.CarCRUDService, customerCRUDService rental.CustomerCRUDService) *Server {
	return &Server{
		CarCRUDService:      carCRUDService,
		CustomerCRUDService: customerCRUDService,
	}
}

type Server struct {
	CarCRUDService      rental.CarCRUDService
	CustomerCRUDService rental.CustomerCRUDService
}

// toAPICar converts a rental.Car to an api.Car and deals with nullable fields.
func toAPICar(car rental.Car) gen.Car {
	var renterID int64
	if car.CustomerID.Valid {
		renterID = car.CustomerID.Int64
	}

	return gen.Car{
		Id:       int64(car.ID),
		Make:     car.Make,
		RenterId: int(renterID),
		Model:    car.Model,
		Year:     car.Year,
	}
}

// toAPICustomer converts a rental.Customer to an api.Customer.
func toAPICustomer(customer rental.Customer) gen.Customer {
	return gen.Customer{
		Id:   int64(customer.ID),
		Name: customer.Name,
	}
}

// Create a new car
// (POST /car)
func (s *Server) CreateCar(ctx echo.Context) error {
	createCar := gen.CreateUpdateCarRequest{}
	if err := ctx.Bind(&createCar); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	car := rental.Car{Make: createCar.Make, Model: createCar.Model, Year: createCar.Year}
	var err error
	car.ID, err = s.CarCRUDService.Create(car)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	apiCar := toAPICar(car)
	return ctx.JSON(http.StatusCreated, apiCar)
}

// Deletes a car
// (DELETE /car/{carId})
func (s *Server) DeleteCar(ctx echo.Context, carId int64) error {
	err := s.CarCRUDService.Delete(int(carId))
	if err == rental.ErrCarNotFound {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.NoContent(http.StatusNoContent)
}

// Find car by ID
// (GET /car/{carId})
func (s *Server) GetCarById(ctx echo.Context, carId int64) error {
	car, err := s.CarCRUDService.Get(int(carId))
	if err == rental.ErrCarNotFound {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	apiCar := toAPICar(car)
	return ctx.JSON(http.StatusOK, apiCar)
}

// Updates a car
// (PUT /car/{carId})
func (s *Server) UpdateCar(ctx echo.Context, carId int64) error {
	CreateCar := gen.CreateUpdateCarRequest{}
	if err := ctx.Bind(&CreateCar); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	car := rental.Car{ID: int(carId), Make: CreateCar.Make, Model: CreateCar.Model, Year: CreateCar.Year}
	err := s.CarCRUDService.Update(car)
	if err == rental.ErrCarNotFound {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	apiCar := toAPICar(car)
	return ctx.JSON(http.StatusOK, apiCar)
}

func (s *Server) rent(carID int, customerID int) error {
	car, err := s.CarCRUDService.Get(carID)
	if err != nil {
		return err
	}

	customer, err := s.CustomerCRUDService.Get(customerID)
	if err != nil {
		return err
	}

	if err := car.Rent(customer.ID); err != nil {
		return err
	}
	return s.CarCRUDService.Update(car)
}

// Rent a car
// (GET /car/{carId}/rent)
func (s *Server) RentCar(ctx echo.Context, carId int64, params gen.RentCarParams) error {
	err := s.rent(int(carId), int(params.CustomerId))
	if err == rental.ErrCarNotFound {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	if err == rental.ErrCustomerNotFound {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err == rental.ErrCarAlreadyRented {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.NoContent(http.StatusNoContent)
}

// Create a new customer
// (POST /customer)
func (s *Server) CreateCustomer(ctx echo.Context) error {
	createCustomer := gen.CreateUpdateCustomerRequest{}
	if err := ctx.Bind(&createCustomer); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	customer := rental.Customer{Name: createCustomer.Name}
	var err error
	customer.ID, err = s.CustomerCRUDService.Create(customer)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	apiCustomer := toAPICustomer(customer)
	return ctx.JSON(http.StatusCreated, apiCustomer)
}

// Deletes a customer
// (DELETE /customer/{customerId})
func (s *Server) DeleteCustomer(ctx echo.Context, customerId int64) error {
	err := s.CustomerCRUDService.Delete(int(customerId))
	if err == rental.ErrCustomerNotFound {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.NoContent(http.StatusNoContent)
}

// Find customer by ID
// (GET /customer/{customerId})
func (s *Server) GetCustomerById(ctx echo.Context, customerId int64) error {
	customer, err := s.CustomerCRUDService.Get(int(customerId))
	if err == rental.ErrCustomerNotFound {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	apiCustomer := toAPICustomer(customer)
	return ctx.JSON(http.StatusOK, apiCustomer)
}

// Updates a customer
// (PUT /customer/{customerId})
func (s *Server) UpdateCustomer(ctx echo.Context, customerId int64) error {
	CreateCustomer := gen.CreateUpdateCustomerRequest{}
	if err := ctx.Bind(&CreateCustomer); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	customer := rental.Customer{ID: int(customerId), Name: CreateCustomer.Name}
	err := s.CustomerCRUDService.Update(customer)
	if err == rental.ErrCustomerNotFound {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	apiCustomer := toAPICustomer(customer)
	return ctx.JSON(http.StatusOK, apiCustomer)
}
