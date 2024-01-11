package main

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"github.com/umtdemr/go-todo/email"
	"github.com/umtdemr/go-todo/logger"
	"github.com/umtdemr/go-todo/server"
	"github.com/umtdemr/go-todo/todo"
	"github.com/umtdemr/go-todo/user"
	"os"
)

func main() {
	viper.SetConfigFile(".env")
	errReadConfig := viper.ReadInConfig()

	if errReadConfig != nil {
		fmt.Fprintf(os.Stderr, "Error while getting env settings: %s\n", errReadConfig)
		os.Exit(1)
	}

	appEnv := viper.Get("APP_ENV")
	if appEnv == nil {
		appEnv = "dev"
	}

	setEnvErr := os.Setenv("APP_ENV", appEnv.(string))

	if setEnvErr != nil {
		fmt.Fprintf(os.Stderr, "Error while setting app env setting: %s\n", setEnvErr)
		os.Exit(1)
	}

	log := logger.Get()

	connStr := viper.Get("postgres")

	store, err := NewPostgresStore(connStr.(string))
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't connect to the db")
	}
	defer store.DB.Close(context.Background())

	email.Init()

	userRepository, err := user.NewUserRepository(store.DB)

	if userRepoInitErr := userRepository.Init(); userRepoInitErr != nil {
		log.Fatal().Msg("Couldn't create user table")
	}

	todoRepository, err := todo.NewTodoRepository(store.DB)

	if todoRepoInitErr := todoRepository.Init(); todoRepoInitErr != nil {
		log.Fatal().Msg("Couldn't create todo table")
	}

	apiServer := server.NewAPIServer(":8080")

	userService := user.NewUserService(userRepository)
	userAPIRoute := user.NewAPIRoute(*userService)
	userAPIRoute.RegisterAPIRoutes(apiServer.Router)

	todoService := todo.NewTodoService(todoRepository)
	todoAPIRoute := todo.NewTodoAPIRoute(todoService)
	todoAPIRoute.RegisterRoutes(apiServer.Router, *userService)

	log.Info().Msg("Server is running")
	apiServer.Run()
}
