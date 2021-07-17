package config

import (
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ichtrojan/thoth"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Database() *gorm.DB {
	logger, _ := thoth.Init("json")

	user, exist := os.LookupEnv("DB_USER")

	if !exist {
		logger.Log(errors.New("DB_USER not set in .env"))
		log.Fatal("DB_USER not set in .env")
	}

	pass, exist := os.LookupEnv("DB_PASS")

	if !exist {
		logger.Log(errors.New("DB_USER not set in .env"))
		log.Fatal("DB_PASS not set in .env")
	}

	host, exist := os.LookupEnv("DB_HOST")

	if !exist {
		logger.Log(errors.New("DB_HOST not set in .env"))
		log.Fatal("DB_HOST not set in .env")
	}

	dbname, exist := os.LookupEnv("DB_NAME")

	if !exist {
		logger.Log(errors.New("DB_NAME not set in .env"))
		log.Fatal("DB_NAME not set in .env")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8&parseTime=True&loc=Local", user, pass, host, dbname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("errore, connessione col database MySQL non andata a buon fine.")
	} else {
		fmt.Println("nessun errore, connessione col database MySQL con successo")
	}
	return db
}
