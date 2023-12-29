package user

type Service struct {
	repository IRepository
}

func NewUserService(repo IRepository) *Service {
	return &Service{repository: repo}
}

func (service *Service) CreateUser(data *CreateUserData) error {
	if len(data.Username) < 3 || len(data.Username) > 20 {
		return ErrorUsernameLength
	}
	return service.repository.CreateUser(data)
}

func (service *Service) Login(data *LoginUserData) bool {
	return service.repository.Login(data)
}
