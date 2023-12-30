package todo

type IService interface {
	GetAllTodos() ([]Todo, error)
}
type Service struct {
	repository IRepository
}

func NewTodoService(repo IRepository) IService {
	return &Service{repository: repo}
}

func (service *Service) GetAllTodos() ([]Todo, error) {
	return service.repository.GetAllTodos()
}
