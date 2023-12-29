package main

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"github.com/umtdemr/go-todo/server"
	"github.com/umtdemr/go-todo/todo"
	"github.com/umtdemr/go-todo/user"
	"os"
)

func main() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	connStr := viper.Get("postgres")

	store, err := NewPostgresStore(connStr.(string))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while connecting to db: %v\n", err)
		os.Exit(1)
	}
	defer store.DB.Close(context.Background())

	apiServer := server.NewAPIServer(":8080")

	todoRepository, err := todo.NewTodoRepository(store.DB)

	if todoRepoInitErr := todoRepository.Init(); todoRepoInitErr != nil {
		fmt.Printf("Error in init: %s\n", todoRepoInitErr)
		os.Exit(1)
	}

	userRepository, err := user.NewUserRepository(store.DB)

	if userRepoInitErr := userRepository.Init(); userRepoInitErr != nil {
		fmt.Printf("Error in init: %s\n", userRepoInitErr)
		os.Exit(1)
	}

	todoAPIRoute := todo.NewTodoAPIRoute(*todoRepository)
	todoAPIRoute.RegisterRoutes(apiServer.Router)

	userService := user.NewUserService(userRepository)
	userAPIRoute := user.NewAPIRoute(*userService)
	userAPIRoute.RegisterAPIRoutes(apiServer.Router)

	apiServer.Run()
}
