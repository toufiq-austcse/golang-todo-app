package services

import (
	"gihub.com/toufiq-austcse/todo-app-go/api/dto"
	"gihub.com/toufiq-austcse/todo-app-go/api/models"
	"gihub.com/toufiq-austcse/todo-app-go/api/repositories"
)

type TodoService struct {
	repository repositories.TodoRepository
}

func NewTodoService(repository repositories.TodoRepository) TodoService {
	return TodoService{repository: repository}
}

func (service TodoService) Create(dto dto.CreateTodoDto, user models.AuthUser) (models.Todo, error) {
	return service.repository.CreateTodo(dto, user)
}
func (service TodoService) FineOneTodoById(id string, user models.AuthUser) (models.Todo, error) {
	return service.repository.FindOneById(id, user)
}
func (service TodoService) FindAll(user models.AuthUser) ([]models.Todo, error) {
	return service.repository.FindAll(user)
}
func (service TodoService) DeleteOneById(id string, user models.AuthUser) (bool, error) {
	return service.repository.DeleteOneById(id, user)
}
func (service TodoService) UpdateOneById(id string, user models.AuthUser, updatedTodo dto.UpdateTodoDto) (models.Todo, error) {
	return service.repository.UpdateOneById(id, user, updatedTodo)
}
