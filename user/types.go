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

type LoginUserData struct {
	Username *string `json:"username,omitempty"`
	Password *string `json:"password,omitempty"`
	Email    *string `json:"email,omitempty"`
}

type VisibleUser struct {
	Id        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
