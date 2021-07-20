package controllers

import (
	"errors"
	"fmt"

	"github.com/Ghecco/saveIT/pkg/config"
	"github.com/Ghecco/saveIT/pkg/models"
	"gorm.io/gorm"
)

const INVALID = -1

var (
	database = config.Database()
)

// User function

// Add User function

func GetIDByUsername(username string) (error, uint64) {
	var user models.User
	err := database.Model(&models.User{}).Where("Name = ?", username).Take(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Printf("Username %s not found in database\n", username)
		return errors.New("username not found\n"), 0
	}
	fmt.Printf("The ID of username:%s is %d\n", username, user.ID)
	return nil, user.ID
}

func GetUsernameByID(UserID uint64) (error, string) {
	var user models.User
	err := database.Model(&models.User{}).Where("ID = ?", UserID).Take(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Printf("ID %d not found in database\n", UserID)
		return errors.New("ID not found\n"), ""
	}
	fmt.Printf("The Username of userID:%d is %s\n", UserID, user.Name)
	return nil, user.Name
}
func AddUser(username, password string) bool {
	var count int64
	database.Model(&models.User{}).Where("Name = ?", username).Count(&count)

	if count != 0 {
		fmt.Printf("The name %s already exists in the database\n", username)
		return false
	}
	if len(username) < 3 || len(username) > 24 {
		fmt.Printf("Username %s is invalid (lenght)\n", username)
		return false
	}

	if len(password) < 3 || len(password) > 24 {
		fmt.Printf("Password %s is invalid (lenght)\n", username)
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

// Remove user function
func RemoveUser(username string) bool {
	var user models.User
	err := database.Where("name = ?", username).Find(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Printf("User %s not found in database\n", username)
		return false
	}
	RemoveIdeas(user.ID)
	database.Where("ID = ?", user.ID).Delete(&user)
	return true
}

// Ideas Function

func GetUserIdeaID(ideaID uint64) (uint64, error) {
	var idea models.Idea
	err := database.Model(&models.Idea{}).Where("ID = ?", ideaID).Take(&idea).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Printf("idea id %d not found in database", idea.ID)
		return 0, errors.New("Idea not found\n")
	}
	err, name := GetUsernameByID(idea.UserID)
	if err != nil {
		fmt.Print("name error")
		return 0, errors.New("name error")
	}
	fmt.Printf("Idea id %d by UserID %s (%s)", ideaID, idea.UserID, name)
	return idea.UserID, nil
}

func AddIdea(userID uint64, content string) bool {
	var count int64
	database.Model(&models.User{}).Where("ID = ?", userID).Count(&count)
	if count == 0 {
		fmt.Printf("UserID %d not exists.\n", userID)
		return false
	}

	if len(content) < 4 || len(content) > 100 {
		fmt.Print("Content lenght error, 4-100\n")
		return false
	}
	idea := models.Idea{UserID: userID, Content: content}
	fmt.Printf("new idea: %d | content: %s\n", idea.UserID, idea.Content)
	database.Create(&idea)
	return true
}

// remove idea function
func RemoveIdea(ideaID uint64) bool {
	var idea models.Idea
	err := database.Where("ID = ?", ideaID).Find(&idea).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Printf("Idea %s not found in database\n", ideaID)
		return false
	}
	database.Where("ID = ?", idea.ID).Delete(&idea)
	fmt.Printf("Idea ID:%d removed.", ideaID)
	return true
}

// remove all ideas  of one user (with userid)
func RemoveIdeas(userID uint64) bool {
	var idea models.Idea
	var count int64
	database.Model(&models.Idea{}).Where("user_id = ?", userID).Count(&count)
	if count == 0 {
		fmt.Printf("No idea related to user id %d\n", userID)
		return false
	}
	database.Where("user_id = ?", userID).Delete(&idea)
	fmt.Printf("Removed %d ideas related to userID %d\n", count, userID)
	return true
}

func GetUserIdeas(userID uint64) (error, []models.Idea) {
	var ideas []models.Idea
	var count int64
	database.Model(&models.Idea{}).Where("user_id = ?", userID).Count(&count)
	if count == 0 {
		fmt.Printf("No idea related to user id %d\n", userID)
		return errors.New("No ideas found in database."), ideas
	}
	database.Table("ideas").Where("user_id = ?", userID).Scan(&ideas)
	//fmt.Println(ideas)
	fmt.Printf("Numbers of ideas of userID %d are %d", userID, count)
	return nil, ideas
}
