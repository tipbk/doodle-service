package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/tipbk/doodle/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type postRepository struct {
	collection *mongo.Collection
}

func NewPostRepository(db *mongo.Database) *postRepository {
	collection := db.Collection("post")
	return &postRepository{
		collection: collection,
	}
}

type PostRepository interface {
	CreatePost(post *model.Post) (*model.Post, error)
	GetPostById(id string) (*model.Post, error)
	GetAllPostByLimitAndOffset(limit, offset int) ([]*model.Post, error)
}

func (r *postRepository) CreatePost(post *model.Post) (*model.Post, error) {
	post.ID = uuid.New().String()
	_, err := r.collection.InsertOne(context.Background(), post)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (r *postRepository) GetPostById(id string) (*model.Post, error) {
	post := model.Post{}
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&post)
	if err != nil {
		return nil, err
	}

	return &post, nil

}

func (r *postRepository) GetAllPostByLimitAndOffset(limit, offset int) ([]*model.Post, error) {
	findOptions := options.Find()
	findOptions.SetSkip(int64(offset))
	findOptions.SetLimit(int64(limit))
	cursor, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	var results []*model.Post
	for cursor.Next(context.Background()) {
		var result model.Post
		err := cursor.Decode(&result)
		if err != nil {
			panic(err)
		}
		results = append(results, &result)
	}
	return results, nil
}
