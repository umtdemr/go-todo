package main

import (
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
	fmt.Println(store)
	server := NewAPIServer(":8080", store)
	server.Run()
}
