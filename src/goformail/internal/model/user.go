package model

type UserRequest struct {
	Email       string   `json:"email"`
	Password    string   `json:"password"`
	Permissions []string `json:"permissions"`
}

type UserResponse struct {
	Id          int      `json:"id"`
	Email       string   `json:"email"`
	Permissions []string `json:"permissions"`
}

var Permissions = []string{
	"ADMIN",
	"CRT_LIST",
	"MOD_LIST",
	"CRT_USER",
	"MOD_USER",
}
