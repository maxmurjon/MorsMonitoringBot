package models

import "time"

type User struct {
	Id          string    `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	UserName    string   `json:"username"`
	PhoneNumber string    `json:"phone"`
	TelegramId    string    `json:"telegram_id"`
	Role 	  string    `json:"role"`
	IsVerified  bool      `json:"is_verified"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserPrimaryKey struct {
	Id           string `json:"id"`
	Phone_number string `json:"phone_number"`
	TelegramId	string `json:"telegram_id"`
}

type CreateUser struct {
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	UserName    string   `json:"username"`
	PhoneNumber string    `json:"phone"`
	TelegramId    string    `json:"telegram_id"`
	Role 	  string    `json:"role"`
}

type UpdateUser struct {
	Id          string    `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	UserName    string   `json:"username"`
	PhoneNumber string    `json:"phone"`
	TelegramId    string    `json:"telegram_id"`
	Role 	  string    `json:"role"`
	IsVerified  bool      `json:"is_verified"`
}

type GetListUserRequest struct {
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	Search string `json:"search"`
}

type GetListUserResponse struct {
	Count int     `json:"count"`
	Users []*User `json:"users"`
}
