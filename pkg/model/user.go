package model

import (
	"errors"
	"fmt"

	passwordvalidator "github.com/wagslane/go-password-validator"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username       string
	HashedPassword string
	Todos          []Todo
}

const (
	MIN_USERNAME_LEN = 5
)

func CreateUser(db *gorm.DB, username string, password string, passwordConfirm string) (*User, error) {

	MINIMUM_ENTROPY := passwordvalidator.GetEntropy("Cemdppue57")

	// check if username is lengthy enough
	if len(username) < MIN_USERNAME_LEN {
		return nil, fmt.Errorf("username must be at of at least length %v", MIN_USERNAME_LEN)
	}

	// check if username is taken

	var users []User
	db.Where(&User{Username: username}).Find(&users)
	if len(users) != 0 {
		return nil, fmt.Errorf("username already taken")
	}

	// check if passwords match
	if password != passwordConfirm {
		return nil, fmt.Errorf("passwords do no match")
	}

	errLocal := passwordvalidator.Validate(password, MINIMUM_ENTROPY)
	// check if password is strong enough
	if errLocal != nil {
		return nil, fmt.Errorf("password is not strong enough")
	}

	// No errors
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	user := &User{
		Username:       username,
		HashedPassword: string(hash),
	}

	result := db.Create(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func LoginUser(db *gorm.DB, username string, password string) (*User, error) {

	var user User
	result := db.Where(&User{Username: username}).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("login failed")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))

	if err != nil {
		return nil, fmt.Errorf("login failed")
	}
	// username/password matches, guy's legit
	return &user, nil
}

func GetUser(db *gorm.DB, userId uint) (*User, error) {
	var user User
	result := db.Where(&User{Model: gorm.Model{
		ID: userId,
	}}).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("user not found")
	}

	return &user, nil

}
