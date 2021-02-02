package models

import (
	"ConchaAPI/database"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// TODO user id to UUID
type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
}

func CreateUser(user User) (*User, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("")
	}

	user.Password = string(pass)
	createdUser := database.DBConn.Create(user)

	if createdUser.Error != nil {
		return nil, fmt.Errorf("")
	}
	return &user, nil
}

func GetUsers() (*[]User, error) {
	var users []User
	database.DBConn.Find(&users)
	return &users, nil
}

func GetUserByEmail(Email string) (*User, error) {
	var user User
	database.DBConn.Where("email = ?", Email).First(&user)
	if user.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return &user, nil
}

func GetUserByID(ID int) (*User, error) {
	var user User
	database.DBConn.Find(&user, ID)
	if user.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return &user, nil
}
