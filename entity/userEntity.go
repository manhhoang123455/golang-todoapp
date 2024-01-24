package entity

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"mini_project/models"
)

type UserEntity interface {
	InsertUser(user models.User) models.User

	VerifyCredential(email string) interface{}

	FindByEmail(email string) models.User
}

type userConnection struct {
	connection *gorm.DB
}

func NewUserEntity(db *gorm.DB) UserEntity {
	return &userConnection{
		connection: db,
	}
}

func (db *userConnection) InsertUser(user models.User) models.User {
	user.Password = hashAndSalt([]byte(user.Password))
	if len(user.Password) > 1 {
		db.connection.Save(&user)
	}

	return user
}

func (db *userConnection) VerifyCredential(email string) interface{} {
	var user models.User
	res := db.connection.Where("email = ?", email).Take(&user)

	if res.Error == nil {
		return user
	}

	return nil
}

func (db *userConnection) FindByEmail(email string) models.User {
	var user models.User
	db.connection.Where("email = ?", email).Take(&user)
	return user
}

func hashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		// if failed to hash password, return empty string
		return ""
	}

	// return hashed password
	return string(hash)
}
