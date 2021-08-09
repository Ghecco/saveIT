package controllers

import (
	"errors"
	"fmt"

	"github.com/Ghecco/saveIT/pkg/models"
	util "github.com/Ghecco/saveIT/pkg/util"
	"gorm.io/gorm"
)

//
// User function

// Obtain id from username in the database
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

// Obtain username from userID in the database
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

// Compare the hash password from the database
func LoginUser(username, password string) bool {
	var user models.User
	err := database.Model(&models.User{}).Where("Name = ?", username).Take(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Printf("Username %s not found in database\n", username)
		return false
	}
	if user.Password != "" {
		match := util.CheckPasswordHash(password, user.Password)
		if match == false {
			return false
		}
	}
	fmt.Printf("Username %s logged.\n", username)
	return true
}

// Add user in the database
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

	if len(password) > 24 { // len(password) < 3 || l
		fmt.Printf("Password %s is invalid (lenght)\n", username)
		return false
	}

	user := models.User{Name: username, Password: password}
	fmt.Printf("%s | %s", user.Name, user.Password)

	if password != "" {
		user.Password, _ = util.HashPassword(password)
	}

	database.Create(&user)
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
