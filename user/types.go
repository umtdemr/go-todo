package user

import "time"

type DBModel struct {
	ID        int64
	Username  string
	Password  string
	Email     string
	IsActive  bool
	CreatedAt time.Time
}

type CreateUserData struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
