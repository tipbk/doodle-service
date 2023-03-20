package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.26

import (
	"context"
	"errors"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/tipbk/doodle/model"
	"github.com/tipbk/doodle/util"
	"go.mongodb.org/mongo-driver/mongo"
)

// Post is the resolver for the post field.
func (r *commentResolver) Post(ctx context.Context, obj *model.Comment) (*model.Post, error) {
	post, err := r.PostRepository.GetPostById(obj.PostId)
	if err != nil {
		return nil, err
	}
	return post, nil
}

// ReplyOn is the resolver for the replyOn field.
func (r *commentResolver) ReplyOn(ctx context.Context, obj *model.Comment) (*model.Comment, error) {
	if obj.ReplyOn == nil {
		return nil, nil
	}

	comment, err := r.CommentRepository.GetCommentById(*obj.ReplyOn)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

// User is the resolver for the user field.
func (r *commentResolver) User(ctx context.Context, obj *model.Comment) (*model.User, error) {
	if obj == nil {
		return nil, errors.New("user not found")
	}
	if obj.UserID == "" {
		return nil, errors.New("user not found")
	}

	user, err := r.UserRepository.GetUserByID(obj.UserID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Register is the resolver for the register field.
func (r *mutationResolver) Register(ctx context.Context, input model.RegisterInput) (*model.AuthResponse, error) {
	if input.Username == "" || input.Password == "" || input.ConfirmPassword == "" {
		return nil, errors.New("fields required")
	}
	if input.Password != input.ConfirmPassword {
		return nil, errors.New("password mismatched")
	}
	_, err := r.UserRepository.GetUserByUsername(input.Username)
	if err != mongo.ErrNoDocuments {
		return nil, errors.New("user already exist")
	}
	user, err := r.UserRepository.CreateUser(input.Username, input.Password)
	if err != nil {
		fmt.Println("error when creating new user")
		return nil, errors.New("error occur")
	}
	token, expTime, err := r.JWTService.GenerateUserToken(user)
	if err != nil {
		return nil, err
	}
	response := model.AuthResponse{
		User: user,
		AuthToken: &model.AuthToken{
			AccessToken: token,
			ExpiredAt:   *expTime,
		},
	}
	return &response, nil
}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, input model.LoginInput) (*model.AuthResponse, error) {
	username := input.Username
	password := input.Password
	user, err := r.UserRepository.GetUserByUsernameAndPassword(username, util.HashPassword(password))
	if err != nil {
		return nil, err
	}
	token, expTime, err := r.JWTService.GenerateUserToken(user)
	if err != nil {
		return nil, err
	}
	authToken := model.AuthToken{
		AccessToken: token,
		ExpiredAt:   *expTime,
	}
	response := model.AuthResponse{
		AuthToken: &authToken,
		User:      user,
	}
	return &response, nil
}

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, input *model.CreatePostInput) (*model.Post, error) {
	header := graphql.GetOperationContext(ctx).Headers
	tokenHeader := header.Get("Authorization")
	token := util.GetTokenString(tokenHeader)
	id, err := r.JWTService.ValidateToken(token)
	if err != nil {
		return nil, err
	}
	post := model.Post{
		UserID:      id,
		Title:       input.Title,
		Description: input.Description,
		Hashtag:     input.Hashtag,
	}

	result, err := r.PostRepository.CreatePost(&post)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// CreateComment is the resolver for the createComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, input *model.CreateCommentInput) (*model.Comment, error) {
	header := graphql.GetOperationContext(ctx).Headers
	tokenHeader := header.Get("Authorization")
	token := util.GetTokenString(tokenHeader)
	id, err := r.JWTService.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	_, err = r.PostRepository.GetPostById(input.PostID)
	if err != nil {
		fmt.Println("no result")
		return nil, err
	}

	comment := model.Comment{
		PostId:  input.PostID,
		UserID:  id,
		Comment: input.Comment,
		ReplyOn: input.ReplyToComment,
	}

	result, err := r.CommentRepository.CreateComment(&comment)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Comment is the resolver for the comment field.
func (r *postResolver) User(ctx context.Context, obj *model.Post) (*model.User, error) {
	user, err := r.UserRepository.GetUserByID(obj.UserID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Comment is the resolver for the comment field.
func (r *postResolver) Comment(ctx context.Context, obj *model.Post) ([]*model.Comment, error) {
	postId := obj.ID
	results, err := r.CommentRepository.FindAllCommentsByPostId(postId)
	if err != nil {
		return nil, err
	}
	return results, nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, input model.UserQueryInput) (*model.User, error) {
	if input.ID == nil && input.Username == nil {
		return nil, errors.New("invalid input")
	}
	if input.ID != nil {
		user, err := r.UserRepository.GetUserByID(*input.ID)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, errors.New("not found")
			}
			fmt.Printf("user repository error: %v", err)
			return nil, errors.New("error when getting user")
		}
		return user, nil
	}
	user, err := r.UserRepository.GetUserByUsername(*input.Username)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("not found")
		}
		return nil, errors.New("error when getting user")
	}
	return user, nil
}

// Post is the resolver for the post field.
func (r *queryResolver) Post(ctx context.Context) ([]*model.Post, error) {
	result, err := r.PostRepository.GetAllPostByLimitAndOffset(5, 0)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetAllPostsByFilter is the resolver for the getAllPostsByFilter field.
func (r *queryResolver) GetAllPostsByFilter(ctx context.Context, input *model.PostFilterInput) (*model.PostResponse, error) {
	result, err := r.PostRepository.GetAllPostByLimitAndOffset(input.Limit, input.Offset)
	if err != nil {
		return nil, err
	}

	total, err := r.PostRepository.GetPostsCount()
	if err != nil {
		return nil, err
	}

	totalInt := int(total)
	postResponse := model.PostResponse{
		Post:      result,
		TotalPost: &totalInt,
	}

	return &postResponse, nil

}

// Comment returns CommentResolver implementation.
func (r *Resolver) Comment() CommentResolver { return &commentResolver{r} }

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Post returns PostResolver implementation.
func (r *Resolver) Post() PostResolver { return &postResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type commentResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type postResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
