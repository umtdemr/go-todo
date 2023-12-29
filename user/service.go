package user

type Service struct {
	repository IRepository
}

func NewUserService(repo IRepository) *Service {
	return &Service{repository: repo}
}

func (service *Service) CreateUser(data *CreateUserData) error {
	return service.repository.CreateUser(data)
}

func (service *Service) Login(data *LoginUserData) bool {
	return service.repository.Login(data)
}
