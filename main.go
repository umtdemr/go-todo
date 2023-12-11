package main

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
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
	defer store.db.Close(context.Background())

	if initErr := store.Init(); initErr != nil {
		fmt.Printf("Error in init: %s\n", initErr)
		os.Exit(1)
	}

	server := NewAPIServer(":8080", store)
	server.Run()
}
