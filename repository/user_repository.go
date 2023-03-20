package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/tipbk/doodle/model"
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
	GetUserByEmail(email string) (*model.User, error)
	GetUserByEmailAndPassword(email, password string) (*model.User, error)
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

func (r *userRepository) GetUserByEmail(email string) (*model.User, error) {
	filter := bson.M{"email": email}
	user := model.User{}
	err := r.collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		fmt.Printf("get user by id error: %v\n", err)
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByEmailAndPassword(email, password string) (*model.User, error) {
	filter := bson.M{"email": email, "password": password}
	user := model.User{}
	err := r.collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		fmt.Printf("get user by id error: %v\n", err)
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) CreateUser(email string, password string) (*model.User, error) {
	user := model.User{
		ID:    uuid.New().String(),
		Email: email,
	}
	user.Password = user.HashPassword(password)
	_, err := r.collection.InsertOne(context.Background(), user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
