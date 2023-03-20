package graph

import (
	"github.com/tipbk/doodle/repository"
	"github.com/tipbk/doodle/service"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	UserRepository    repository.UserRepository
	CommentRepository repository.CommentRepository
	PostRepository    repository.PostRepository
	JWTService        service.JWTService
}
