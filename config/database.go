package config

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ichtrojan/thoth"
)

func Database() *sql.DB {
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

	credentials := fmt.Sprintf("%s:%s@(%s:3306)/?charset=utf8&parseTime=true", user, pass, host)

	database, err := sql.Open("mysql", credentials)

	if err != nil {
		logger.Log(err)
		log.Fatal(err)
	} else {
		fmt.Println("MySQL database connection successful.")
	}
	return database
}
