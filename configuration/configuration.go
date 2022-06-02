package configuration

import (
	"encoding/json"
	"log"
	"os"

	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type configuration struct {
	Repository repositoryConfig
}

type repositoryConfig struct {
	Server                string
	Port                  string
	User                  string
	Password              string
	DatabaseConfiguration string
}

func getConfiguration() *configuration {
	var c configuration
	file, err := os.Open("./config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&c)
	if err != nil {
		log.Fatal(err)
	}

	userDB := os.Getenv("userDB")
	passwordDB := os.Getenv("passwordDB")

	if userDB == "" || passwordDB == "" {
		log.Fatal("Environment Variables userDB or passwordDB is empty")
	}

	c.Repository.User = userDB
	c.Repository.Password = passwordDB

	return &c
}

// GetConnection obtiene una conexión a la bd
func GetConnectionStringCalculator() string {
	c := getConfiguration()
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", c.Repository.User, c.Repository.Password, c.Repository.Server, c.Repository.Port, c.Repository.DatabaseConfiguration)
}
