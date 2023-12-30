package todo

type Service struct {
	repository IRepository
}

func NewTodoService(repo IRepository) *Service {
	return &Service{repository: repo}
}
