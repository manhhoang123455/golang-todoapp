package services

import (
	"github.com/mashingan/smapping"
	"golang.org/x/crypto/bcrypt"
	"mini_project/entity"
	"mini_project/models"
	"mini_project/request"
)

type UserService interface {
	VerifyCredential(email string, password string) interface{}

	CreateUser(user request.RegisterRequest) models.User

	//FindByEmail(email string, password string) models.User
}

type userService struct {
	userEntity entity.UserEntity
}

func NewUserService(userEntity entity.UserEntity) UserService {
	return &userService{userEntity: userEntity}
}

func (s *userService) VerifyCredential(email string, password string) interface{} {
	res := s.userEntity.VerifyCredential(email)
	if v, ok := res.(models.User); ok {
		comparedPassword := comparePassword(v.Password, []byte(password))
		if v.Email == email && comparedPassword {
			return res
		}
		return false
	}
	return false
}

func (s *userService) CreateUser(user request.RegisterRequest) models.User {
	// create user model
	userToCreate := models.User{}

	// fill user model with data from request model and return error if any error occur during mapping process or return nil if no error occur during mapping process and return user model to caller function to use it
	err := smapping.FillStruct(&userToCreate, smapping.MapFields(&user))
	if err != nil {
		panic(err)
	}

	findByEmail := s.userEntity.FindByEmail(user.Email)

	if findByEmail.ID == 0 {
		res := s.userEntity.InsertUser(userToCreate)
		return res
	}

	return userToCreate
}

func comparePassword(hashedPwd string, plainPassword []byte) bool {
	byteHash := []byte(hashedPwd)

	err := bcrypt.CompareHashAndPassword(byteHash, plainPassword)
	if err != nil {
		return false
	}

	return true
}
