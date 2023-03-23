package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/tipbk/doodle/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type commentRepository struct {
	collection *mongo.Collection
}

func NewCommentRepository(db *mongo.Database) *commentRepository {
	collection := db.Collection("comment")
	return &commentRepository{
		collection: collection,
	}
}

type CommentRepository interface {
	CreateComment(comment *model.Comment) (*model.Comment, error)
	FindAllCommentsByPostId(postId string) ([]*model.Comment, error)
	GetCommentById(commentId string) (*model.Comment, error)
}

func (r *commentRepository) CreateComment(comment *model.Comment) (*model.Comment, error) {
	id := uuid.New().String()
	comment.ID = id
	now := time.Now()
	comment.CreatedAt = &now
	_, err := r.collection.InsertOne(context.Background(), comment)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (r *commentRepository) FindAllCommentsByPostId(postId string) ([]*model.Comment, error) {
	filter := bson.M{"postId": postId}
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"createdAt", 1}})
	cur, err := r.collection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	var results []*model.Comment
	for cur.Next(context.Background()) {
		comment := model.Comment{}
		err := cur.Decode(&comment)
		if err != nil {
			return nil, err
		}
		results = append(results, &comment)
	}

	return results, nil
}

func (r *commentRepository) GetCommentById(commentId string) (*model.Comment, error) {
	filter := bson.M{"_id": commentId}
	comment := model.Comment{}
	err := r.collection.FindOne(context.Background(), filter).Decode(&comment)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}
