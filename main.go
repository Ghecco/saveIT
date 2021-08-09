package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Ghecco/saveIT/pkg/telegram"
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
		log.Fatal("VERSION not set in .env")
	}
	fmt.Printf("saveIT Version:%s ", version)

	// Telegram Server

	telegramToken, tokenExist := os.LookupEnv("TELEGRAM_TOKEN")

	if !tokenExist {
		logger.Log(errors.New("TELEGRAM_TOKEN not set in .env"))
		log.Fatal("TELEGRAM_TOKEN not set in .env")
	}

	telegram.Init(telegramToken)

}
