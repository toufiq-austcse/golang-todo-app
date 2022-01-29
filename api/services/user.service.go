package services

import (
	"gihub.com/toufiq-austcse/todo-app-go/api/dto"
	"gihub.com/toufiq-austcse/todo-app-go/api/models"
	"gihub.com/toufiq-austcse/todo-app-go/api/repositories"
)

type UserService struct {
	repository repositories.UserRepository
}

func NewUserService(repository repositories.UserRepository) UserService {
	return UserService{repository: repository}
}

func (service UserService) InsertUser(dto dto.RegisterUserDto) (models.User, error) {
	return service.repository.InsertUser(dto)

}

func (service UserService) findUserByEmail(email string) (models.User, error) {
	return service.repository.FindUserByEmail(email)

}
func (service UserService) FindUserById(id string) (models.User, error) {
	return service.repository.FindUserByUserId(id)
}
