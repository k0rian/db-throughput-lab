package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/db-throughput-lab/internal/config"
	"github.com/db-throughput-lab/internal/controllers"
	"github.com/db-throughput-lab/internal/repository"
	"github.com/db-throughput-lab/internal/services"
)

func main() {
	// Initialize database and redis connections
	conns := config.InitConnections()

	// Initialize repositories
	mysqlRepo := repository.NewUserRepository(conns.MySQL, nil, nil, conns.Redis, "mysql")
	pgRepo := repository.NewUserRepository(nil, conns.Postgres, nil, conns.Redis, "postgres")
	mongoRepo := repository.NewUserRepository(nil, nil, conns.Mongo, conns.Redis, "mongo")

	// Initialize services
	benchService := services.NewBenchmarkService(mysqlRepo, pgRepo, mongoRepo)

	// Initialize controllers
	benchController := controllers.NewBenchmarkController(benchService)

	// Set up Gin router
	router := gin.Default()

	// Routes
	api := router.Group("/api/v1")
	{
		api.POST("/benchmark/write", benchController.RunWriteTest)
		api.POST("/benchmark/read", benchController.RunReadTest)
	}

	log.Println("Server is starting on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
