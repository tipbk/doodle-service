package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/tipbk/doodle/model"
	"github.com/tipbk/doodle/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *userRepository {
	collection := db.Collection("user")
	return &userRepository{
		collection: collection,
	}
}

type UserRepository interface {
	GetUserByID(id string) (*model.User, error)
	GetUserByUsername(username string) (*model.User, error)
	GetUserByUsernameAndPassword(username, password string) (*model.User, error)
	CreateUser(email string, password string) (*model.User, error)
}

func (r *userRepository) GetUserByID(id string) (*model.User, error) {
	filter := bson.M{"_id": id}
	user := model.User{}
	err := r.collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		fmt.Printf("get user by id error: %v\n", err)
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByUsername(username string) (*model.User, error) {
	filter := bson.M{"username": username}
	user := model.User{}
	err := r.collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		fmt.Printf("get user by id error: %v\n", err)
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByUsernameAndPassword(username, password string) (*model.User, error) {
	filter := bson.M{"username": username, "password": password}
	user := model.User{}
	err := r.collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		fmt.Printf("get user by id error: %v\n", err)
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) CreateUser(username string, password string) (*model.User, error) {
	user := model.User{
		ID:       uuid.New().String(),
		Username: username,
	}
	user.Password = util.HashPassword(password)
	_, err := r.collection.InsertOne(context.Background(), user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
