/*
Contains functions to set up the server and run it.
*/
package main

import (
	"context"
	"log"
	"net/http"

	"cloud.google.com/go/storage"
	ginSess "github.com/gin-contrib/sessions"
	redisSess "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/ryanozx/skillnet-milestone2-backend/database"
	"github.com/ryanozx/skillnet-milestone2-backend/helpers"
	"google.golang.org/api/option"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load("./../.env")
	if err != nil {
		panic(err)
	}
	serverConfig := initialiseProdServer()
	serverConfig.setupRoutes()
	serverConfig.runRouter()
	log.Println("Setup complete!")
}

// serverConfig contains the essentials to run the backend - a router,
// a Redis database for fast reads, and a database for persistent data
type serverConfig struct {
	db          *gorm.DB
	store       redisSess.Store
	router      *gin.Engine
	GoogleCloud *storage.Client
}

// Returns a server configuration with the production database (as defined
// by environmental variables) set as the database
func initialiseProdServer() *serverConfig {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	db := database.ConnectProdDatabase()
	store := setupRedis()
	server := serverConfig{
		db:     db,
		router: router,
		store:  store,
	}
	server.setupGoogleCloud()
	return &server
}

// Sets up the Redis store from environmental variables
func setupRedis() redisSess.Store {
	env := helpers.RetrieveRedisEnv()
	redisAddress := env.Address()
	const redisNetwork = "tcp"
	store, err := redisSess.NewStore(env.MaxConn, redisNetwork, redisAddress, "", []byte(env.Secret))
	if err != nil {
		log.Fatal(err.Error())
	}
	store.Options(ginSess.Options{
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		MaxAge:   86400,
	})
	return store
}

func (s *serverConfig) setupGoogleCloud() {
	ctx := context.Background()
	env := helpers.RetrieveGoogleCloudEnv()

	client, err := storage.NewClient(ctx, option.WithCredentialsFile(env.Filepath))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	s.GoogleCloud = client
}

func (server *serverConfig) runRouter() {
	env := helpers.RetrieveWebAppEnv()
	routerAddress := env.Address()
	err := server.router.Run(routerAddress)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
}
