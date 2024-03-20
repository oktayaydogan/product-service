package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	. "github.com/swaggo/files"
	. "github.com/swaggo/gin-swagger"
	"product-service/api"
	"product-service/config"
	"product-service/database"
	_ "product-service/docs"
)

// @title Product Service API
// @version 1.0
// @description This API provides endpoints for managing products and health check.

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io
// @BasePath /v2

// @Router /swagger/*any [get]
// @Router /health [get]
// @Router /graphql [post]

func main() {
	if err := runServer(); err != nil {
		log.Fatal(err)
	}
}

func runServer() error {
	if err := config.LoadConfig(); err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		config.AppConfig.DBUsername,
		config.AppConfig.DBPassword,
		config.AppConfig.DBHost,
		config.AppConfig.DBPort,
		config.AppConfig.DBName,
	)

	if err := database.InitDB(dsn); err != nil {
		return fmt.Errorf("failed to initialize database: %v", err)
	}

	r := setupRouter()

	port := ":8080"

	if err := r.Run(port); err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}

	return nil
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	// @Router /swagger/*any [get]
	r.GET("/swagger/*any", WrapHandler(Handler))

	// @Router /health [get]
	r.GET("/health", api.HealthHandler)

	// @Router /graphql [post]
	r.POST("/graphql", api.GraphQLHandler)

	return r
}
