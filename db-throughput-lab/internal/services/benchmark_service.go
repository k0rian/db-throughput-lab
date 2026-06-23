package services

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/db-throughput-lab/internal/models"
	"github.com/db-throughput-lab/internal/repository"
	"github.com/db-throughput-lab/pkg/benchmark"
)

type BenchmarkService interface {
	RunWriteTest(dbType string, concurrency int, duration time.Duration) (benchmark.Result, error)
	RunReadTest(dbType string, concurrency int, duration time.Duration, useCache bool) (benchmark.Result, error)
}

type benchmarkService struct {
	mysqlRepo    repository.UserRepository
	postgresRepo repository.UserRepository
	mongoRepo    repository.UserRepository
}

func NewBenchmarkService(mysqlRepo, postgresRepo, mongoRepo repository.UserRepository) BenchmarkService {
	return &benchmarkService{
		mysqlRepo:    mysqlRepo,
		postgresRepo: postgresRepo,
		mongoRepo:    mongoRepo,
	}
}

func (s *benchmarkService) getRepo(dbType string) (repository.UserRepository, error) {
	switch dbType {
	case "mysql":
		return s.mysqlRepo, nil
	case "postgres":
		return s.postgresRepo, nil
	case "mongo":
		return s.mongoRepo, nil
	default:
		return nil, fmt.Errorf("unsupported db type: %s", dbType)
	}
}

func (s *benchmarkService) RunWriteTest(dbType string, concurrency int, duration time.Duration) (benchmark.Result, error) {
	repo, err := s.getRepo(dbType)
	if err != nil {
		return benchmark.Result{}, err
	}

	task := func() error {
		ctx := context.Background()
		if dbType == "mongo" {
			user := &models.MongoUser{
				Name:      fmt.Sprintf("User_%d", rand.Intn(1000000)),
				Age:       rand.Intn(100),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			return repo.CreateMongo(ctx, user)
		}

		user := &models.User{
			Name:      fmt.Sprintf("User_%d", rand.Intn(1000000)),
			Age:       rand.Intn(100),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		return repo.Create(ctx, user)
	}

	testName := fmt.Sprintf("Write Test - %s", dbType)
	return benchmark.RunThroughputTest(testName, duration, concurrency, task), nil
}

func (s *benchmarkService) RunReadTest(dbType string, concurrency int, duration time.Duration, useCache bool) (benchmark.Result, error) {
	repo, err := s.getRepo(dbType)
	if err != nil {
		return benchmark.Result{}, err
	}

	// Prepare data to ensure there's something to read
	var maxID uint = 1000
	var maxOIDs []primitive.ObjectID

	if dbType == "mongo" {
		for i := 0; i < 100; i++ {
			user := &models.MongoUser{Name: "Test", Age: 20}
			repo.CreateMongo(context.Background(), user)
			maxOIDs = append(maxOIDs, user.ID)
		}
	} else {
		for i := 0; i < 100; i++ {
			user := &models.User{Name: "Test", Age: 20}
			repo.Create(context.Background(), user)
			if user.ID > maxID {
				maxID = user.ID
			}
		}
	}

	task := func() error {
		ctx := context.Background()
		if dbType == "mongo" {
			// Randomly pick one from the inserted ones
			id := maxOIDs[rand.Intn(len(maxOIDs))]
			_, err := repo.GetMongoByID(ctx, id, useCache)
			return err
		}

		// For relational DB, pick a random ID up to maxID
		id := uint(rand.Intn(int(maxID))) + 1
		_, err := repo.GetByID(ctx, id, useCache)
		return err
	}

	cacheStatus := "Without Cache"
	if useCache {
		cacheStatus = "With Cache"
	}
	testName := fmt.Sprintf("Read Test - %s (%s)", dbType, cacheStatus)
	return benchmark.RunThroughputTest(testName, duration, concurrency, task), nil
}
