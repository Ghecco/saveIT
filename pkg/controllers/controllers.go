package controllers

import (
	"fmt"

	"github.com/Ghecco/saveIT/pkg/config"
	"github.com/Ghecco/saveIT/pkg/models"
)

var (
	database = config.Database()
)

func AddUser(username, password string) bool {
	var count int64
	database.Model(&models.User{}).Where("Name = ?", username).Count(&count)

	if count != 0 {
		fmt.Printf("The name %s already exists in the database\n", username)
		return false
	}
	if len(username) < 3 || len(username) > 24 {
		fmt.Printf("Username %s is invalid (lenght)", username)
		return false
	}

	if len(password) < 3 || len(password) > 24 {
		fmt.Printf("Password %s is invalid (lenght)", username)
		return false
	}

	user := models.User{Name: username, Password: password}
	fmt.Printf("%s | %s", user.Name, user.Password)

	err := database.Create(&user)
	if err != nil {
		fmt.Print("errore")
		return false
	}
	return true
}

func AddIdea(userID uint64, content string) bool {
	var count int64
	database.Model(&models.User{}).Where("ID = ?", userID).Count(&count)
	if count == 0 {
		fmt.Println("UserID %d not exists.", userID)
		return false
	}

	if len(content) < 4 || len(content) > 100 {
		fmt.Print("Content lenght error, 4-100")
		return false
	}
	idea := models.Idea{UserID: userID, Content: content}
	fmt.Printf("new idea: %d | content: %s", idea.UserID, idea.Content)
	err := database.Create(&idea)
	if err != nil {
		fmt.Print("errore")
		return false
	}
	return true
}
