package app

import "time"

type Photo struct {
	Title    string `json:"title"`
	Caption  string `json:"caption"`
	PhotoUrl string `json:"photo_url"`
}

type Owner struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserData struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Photos   Photo  `json:"photos"`
	Token    string `json:"token"`
}

type UserLogin struct {
	ID       string `json:"id"`
	Username string `json:"users.username"`
	Email    string `json:"users.email"`
	Token    string `json:"users.token"`
	Password string `json:"users.password"`
	Title    string `json:"photos.title"`
	Caption  string `json:"photos.caption"`
	PhotoUrl string `json:"photos.photo_url"`
}

type UserRegister struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}