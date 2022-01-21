package service

import (
	"errors"
	"fitgoapi/model"
	"fitgoapi/utils"
	"fmt"
)

type IUserService interface {
	Login(username string, password string) bool
	CreateUser(user model.User) (model.User, error)
	ValidateEmail(email string) (bool, error)
	FindAllInstructorsByCityCode(userId int64, cityCode string, offset int, size int) (model.PageableUser, error)
	MakeConnection(userID int64, contactID int64) error
	LoadUserConnections(userID int64, page int, size int) (model.PageableUserPlain, error)
	UpdateUser(userID int64, updatedUser model.User) error
	UpdateFcmToken(username string, fcmToken string) error
}

type UserService struct {
}

func (service UserService) Login(username string, password string) bool {
	user, err := userRepository.FetchUserByEmailAndPassword(username, password)

	if err != nil {
		fmt.Printf("%s", err.Error())
		return false
	}

	return username == user.Username
}

func (service UserService) CreateUser(user model.User) (model.User, error) {

	isValid, err := service.ValidateEmail(user.Username)
	if err != nil {
		return model.User{}, err
	}

	if isValid {
		_, user, err := userRepository.SaveUser(user)
		return user, err
	}

	return user, errors.New("email inválido ou já utilizado")
}

func (service UserService) ValidateEmail(email string) (bool, error) {

	if !utils.IsEmailValid(email) {
		return false, nil
	}

	valid, err := userRepository.VerifyEmail(email)

	if err != nil {
		return false, err
	}

	return valid, nil
}

func (service UserService) FindAllInstructorsByCityCode(userId int64, cityCode string, offset int, size int) (model.PageableUser, error) {
	return userRepository.LoadInstructorsByCityCode(userId, cityCode, offset, size)
}

func (service UserService) MakeConnection(userID int64, contactID int64) error {
	return userRepository.MakeConnection(userID, contactID)
}

func (service UserService) LoadUserConnections(userID int64, page int, size int) (model.PageableUserPlain, error) {
	return userRepository.LoadUserConnections(userID, page, size)
}

func (service UserService) UpdateUser(userID int64, updatedUser model.User) error {
	return userRepository.UpdateUser(userID, updatedUser)
}

func (service UserService) UpdateFcmToken(username string, fcmToken string) error {
	return userRepository.UpdateFcmToken(username, fcmToken)
}
