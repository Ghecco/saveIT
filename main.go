package main

import (
	"errors"
	"log"
	"os"

	"github.com/Ghecco/saveIT/pkg/config"
	"github.com/Ghecco/saveIT/pkg/controllers"
	"github.com/ichtrojan/thoth"
	"github.com/joho/godotenv"
)

func main() {
	logger, err := thoth.Init("json")
	if err != nil {
		log.Fatal(err)
	}

	if err := godotenv.Load(); err != nil {
		logger.Log(errors.New(".env is not found"))
		log.Fatal(".env is not found.")
	}

	version, versionExist := os.LookupEnv("VERSION")

	if !versionExist {
		logger.Log(errors.New("VERSION not set in .env"))
		log.Fatal("PORT not set in .env")
	}

	port, exist := os.LookupEnv("PORT")

	if !exist {
		logger.Log(errors.New("PORT not set in .env"))
		log.Fatal("PORT not set in .env")
	}

	// Testing

	config.Database()
	controllers.AddUser("wqdwd", "ciao")
	controllers.AddIdea(1, "ciao nuova idea")

	log.Printf("saveIT Version:%s ", version)
	log.Print("Started web app on port :" + port)

}
