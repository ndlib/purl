package main

type User struct {
	Id       int    `json:"id"`
	username string `json:"username"`
}

type Users []User

type UserRole struct {
	user_id int `json:"user_id"`
	role_id int `json:"role_id"`
}

type UserRoles []UserRole
