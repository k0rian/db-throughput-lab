package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/db-throughput-lab/internal/models"
)

type AppConnections struct {
	MySQL    *gorm.DB
	Postgres *gorm.DB
	Mongo    *mongo.Database
	Redis    *redis.Client
}

func InitConnections() *AppConnections {
	mysqlDB := initMySQL()
	pgDB := initPostgres()
	mongoDB := initMongo()
	redisClient := initRedis()

	// Auto-migrate tables for relational DBs
	if err := mysqlDB.AutoMigrate(&models.User{}); err != nil {
		log.Printf("MySQL migration failed: %v", err)
	}
	if err := pgDB.AutoMigrate(&models.User{}); err != nil {
		log.Printf("Postgres migration failed: %v", err)
	}

	return &AppConnections{
		MySQL:    mysqlDB,
		Postgres: pgDB,
		Mongo:    mongoDB,
		Redis:    redisClient,
	}
}

func initMySQL() *gorm.DB {
	dsn := "root:root@tcp(localhost:3306)/lab_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	setupPool(db)
	return db
}

func initPostgres() *gorm.DB {
	dsn := "host=localhost user=root password=root dbname=lab_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
	}
	setupPool(db)
	return db
}

func setupPool(db *gorm.DB) {
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetMaxOpenConns(500)
	sqlDB.SetConnMaxLifetime(time.Hour)
}

func initMongo() *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	clientOptions := options.Client().ApplyURI("mongodb://root:root@localhost:27017")
	clientOptions.SetMaxPoolSize(500)
	
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	
	return client.Database("lab_db")
}

func initRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
		PoolSize: 500,
	})
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("Warning: Redis not connected: %v\n", err)
	}
	
	return rdb
}
