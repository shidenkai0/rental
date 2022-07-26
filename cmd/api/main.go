package main

import (
	"crypto/subtle"
	"fmt"
	"net/http"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/shidenkai0/rental/pkg/api"
	"github.com/shidenkai0/rental/pkg/api/gen"
	"github.com/shidenkai0/rental/pkg/database"
	"github.com/spf13/viper"
)

func main() {
	// Set default config
	viper.SetDefault("port", "8080")
	viper.SetDefault("database_url", "postgres://rental:rental@localhost:5432/rental")
	viper.SetDefault("database_max_pool_size", "5")
	viper.SetDefault("database_max_idle_connections", "2")
	viper.SetDefault("basic_auth_username", "rental")
	viper.SetDefault("basic_auth_password", "rental")

	// Read config from env
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	port := viper.GetString("port")
	databaseURL := viper.GetString("database_url")
	databaseMaxOpenConns := viper.GetInt("database_max_open_conns")
	databaseMaxIdleConns := viper.GetInt("database_max_idle_conns")
	basicAuthUsername := viper.GetString("basic_auth_user")
	basicAuthPassword := viper.GetString("basic_auth_password")

	// Log config
	fmt.Printf("port: %s\n", port)
	fmt.Printf("database_url: %s\n", databaseURL)
	fmt.Printf("database_max_open_conns: %d\n", databaseMaxOpenConns)
	fmt.Printf("database_max_idle_conns: %d\n", databaseMaxIdleConns)
	fmt.Printf("basic_auth_user: %s\n", basicAuthUsername)
	fmt.Printf("basic_auth_password: %s\n", basicAuthPassword)

	// Setup echo middleware

	e := echo.New()

	// Setup healthcheck for Kubernetes probes
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	// Connect to database
	db := sqlx.MustConnect("postgres", databaseURL)
	db.SetMaxOpenConns(databaseMaxOpenConns)
	db.SetMaxIdleConns(databaseMaxIdleConns)
	defer db.Close()

	// Setup API server
	carCRUDService := database.NewDatabaseCarCRUDService(db)
	customerCRUDService := database.NewDatabaseCustomerCRUDService(db)
	server := api.NewServer(carCRUDService, customerCRUDService)

	// Setup API middleware
	v1APIGroup := e.Group("/v1")
	v1APIGroup.Use(middleware.Logger())
	v1APIGroup.Use(middleware.Recover())
	v1APIGroup.Use(middleware.CORS())
	v1APIGroup.Use(middleware.Gzip())
	v1APIGroup.Use(middleware.Secure())
	v1APIGroup.Use(middleware.BodyLimit("1M"))
	v1APIGroup.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		// Use constant time comparison to prevent timing attacks
		if subtle.ConstantTimeCompare([]byte(username), []byte(basicAuthUsername)) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte(basicAuthPassword)) == 1 {
			return true, nil
		}
		return false, nil
	}))

	gen.RegisterHandlers(v1APIGroup, server)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
