package todo

import (
	"errors"
)

type IService interface {
	GetAllTodos(userId int64) ([]Todo, error)
	CreateTodo(data *CreateTodoData, userId int64) (*Todo, error)
	UpdateTodo(data *UpdateTodoData, userId int64) (*Todo, error)
	RemoveTodo(todoId int, userId int64) (*Todo, error)
}
type Service struct {
	Repository IRepository
}

func NewTodoService(repo IRepository) *Service {
	return &Service{Repository: repo}
}

func (service *Service) GetAllTodos(userId int64) ([]Todo, error) {
	return service.Repository.GetAllTodos(userId)
}

func (service *Service) CreateTodo(data *CreateTodoData, userId int64) (*Todo, error) {
	if data.Title == "" {
		return nil, errors.New("title need to be sent")
	}

	createTodoData := NewTodo(data.Title)
	return service.Repository.CreateTodo(createTodoData, userId)
}
func (service *Service) UpdateTodo(data *UpdateTodoData, userId int64) (*Todo, error) {
	return service.Repository.UpdateTodo(data, userId)
}

func (service *Service) RemoveTodo(todoId int, userId int64) (*Todo, error) {
	return service.Repository.RemoveTodo(todoId, userId)
}

func (service *Service) GetTodo(todoId int, userId int64) (*Todo, error) {
	return service.Repository.GetTodo(todoId, userId)
}
