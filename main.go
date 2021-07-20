package main

import (
	"errors"
	"fmt"
	"log"
	"os"

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

	controllers.AddIdea(5, "ciao nuova idea")
	controllers.AddIdea(5, "3REFFSF")
	controllers.AddIdea(5, "RSDFSDSS")
	//controllers.RemoveUser("wqdwd")
	//	controllers.RemoveIdeas(5)
	_, id := controllers.GetIDByUsername("wqdwd")
	fmt.Printf("name %d", id)

	err, ideas := controllers.GetUserIdeas(5)
	fmt.Println(ideas)
	log.Printf("saveIT Version:%s ", version)
	log.Print("Started web app on port :" + port)

}
