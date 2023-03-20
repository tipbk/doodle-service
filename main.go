package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/tipbk/doodle/graph"
	"github.com/tipbk/doodle/repository"
	"github.com/tipbk/doodle/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const defaultPort = "8080"

func main() {
	godotenv.Load(".env")
	db := initiateMongoClient()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	userRepository := repository.NewUserRepository(db)
	commentRepository := repository.NewCommentRepository(db)
	postRepository := repository.NewPostRepository(db)
	jwtService := service.NewJWTService(os.Getenv("JWT_SECRET"))
	resolver := &graph.Resolver{
		UserRepository:    userRepository,
		CommentRepository: commentRepository,
		PostRepository:    postRepository,
		JWTService:        jwtService,
	}
	r.POST("/query", graphqlHandler(db, resolver))
	r.GET("/", playgroundHandler())
	r.Run()
}

func graphqlHandler(db *mongo.Database, gr *graph.Resolver) gin.HandlerFunc {

	h := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: gr}))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func initiateMongoClient() *mongo.Database {
	user := os.Getenv("MONGO_USER")
	password := os.Getenv("MONGO_PASSWORD")
	url := fmt.Sprintf("mongodb+srv://%v:%v@cluster0.laosc.mongodb.net/?retryWrites=true&w=majority", user, password)
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(url))
	if err != nil {
		log.Fatal(err)
	}
	db := client.Database(os.Getenv("MONGO_DATABASE"))
	return db
}
