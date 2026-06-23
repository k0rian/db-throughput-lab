package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"

	"github.com/db-throughput-lab/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uint, useCache bool) (*models.User, error)
	CreateMongo(ctx context.Context, user *models.MongoUser) error
	GetMongoByID(ctx context.Context, id primitive.ObjectID, useCache bool) (*models.MongoUser, error)
}

type userRepository struct {
	mysqlDB  *gorm.DB
	pgDB     *gorm.DB
	mongoDB  *mongo.Database
	redisCli *redis.Client
	dbType   string // "mysql" or "postgres"
}

func NewUserRepository(mysqlDB, pgDB *gorm.DB, mongoDB *mongo.Database, redisCli *redis.Client, dbType string) UserRepository {
	return &userRepository{
		mysqlDB:  mysqlDB,
		pgDB:     pgDB,
		mongoDB:  mongoDB,
		redisCli: redisCli,
		dbType:   dbType,
	}
}

func (r *userRepository) getDB() *gorm.DB {
	if r.dbType == "postgres" {
		return r.pgDB
	}
	return r.mysqlDB
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	return r.getDB().WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetByID(ctx context.Context, id uint, useCache bool) (*models.User, error) {
	cacheKey := fmt.Sprintf("user:%s:%d", r.dbType, id)

	if useCache && r.redisCli != nil {
		val, err := r.redisCli.Get(ctx, cacheKey).Result()
		if err == nil {
			var user models.User
			if json.Unmarshal([]byte(val), &user) == nil {
				return &user, nil
			}
		}
	}

	var user models.User
	if err := r.getDB().WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}

	if useCache && r.redisCli != nil {
		if data, err := json.Marshal(user); err == nil {
			r.redisCli.Set(ctx, cacheKey, data, 5*time.Minute)
		}
	}

	return &user, nil
}

func (r *userRepository) CreateMongo(ctx context.Context, user *models.MongoUser) error {
	collection := r.mongoDB.Collection("users")
	res, err := collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		user.ID = oid
	}
	return nil
}

func (r *userRepository) GetMongoByID(ctx context.Context, id primitive.ObjectID, useCache bool) (*models.MongoUser, error) {
	cacheKey := fmt.Sprintf("user:mongo:%s", id.Hex())

	if useCache && r.redisCli != nil {
		val, err := r.redisCli.Get(ctx, cacheKey).Result()
		if err == nil {
			var user models.MongoUser
			if json.Unmarshal([]byte(val), &user) == nil {
				return &user, nil
			}
		}
	}

	collection := r.mongoDB.Collection("users")
	var user models.MongoUser
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}

	if useCache && r.redisCli != nil {
		if data, err := json.Marshal(user); err == nil {
			r.redisCli.Set(ctx, cacheKey, data, 5*time.Minute)
		}
	}

	return &user, nil
}
