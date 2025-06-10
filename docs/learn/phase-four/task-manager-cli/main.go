package main

import (
	"fmt"
	"github.com/spf13/viper"
)

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s\n", err)
	} else {
		fmt.Println("Config file loaded successfully")
	}
}

func init() {
	initConfig()
}

func main() {
	databaseEngine := viper.GetString("database.engine")
	databaseUser := viper.GetString("database.user")
	databasePassword := viper.GetString("database.password")
	databaseHost := viper.GetString("database.host")
	databasePort := viper.GetInt("database.port")

	if databaseEngine == "" || databaseUser == "" || databasePassword == "" || databaseHost == "" || databasePort == 0 {
		fmt.Println("Missing required database configuration")
		return
	} else {
		fmt.Println("Database URL:", fmt.Sprintf("%s://%s:%s@%s:%d", databaseEngine, databaseUser, databasePassword, databaseHost, databasePort))
	}

}

// go build -o task-manager-cli main.go
// ./task-manager-cli
