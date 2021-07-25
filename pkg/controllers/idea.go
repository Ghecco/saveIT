package controllers

import (
	"errors"
	"fmt"
	"time"

	"github.com/Ghecco/saveIT/pkg/models"
	"gorm.io/gorm"
)

// Ideas Function

// Get precise idea
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

// Add idea in the database
func AddIdea(userID uint64, title string, content string) bool {
	var count int64
	database.Model(&models.User{}).Where("ID = ?", userID).Count(&count)
	if count == 0 {
		fmt.Printf("UserID %d not exists.\n", userID)
		return false
	}

	if len(title) < 4 || len(title) > 50 {
		fmt.Print("Title lenght error, 4-50\n")
		return false
	}
	dt := time.Now()

	idea := models.Idea{UserID: userID, Title: title, Content: content, Time: dt.Format("01-02-2OO6 15:04")}
	fmt.Printf("new idea: %d | Title: %s Time:%s\n", idea.UserID, idea.Title, idea.Time)
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

// Get all user's ideas
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
