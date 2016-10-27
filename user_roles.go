package main

type UserRole struct {
	user_id int `json:"user_id"`
	role_id int `json:"role_id"`
}

type UserRoles []UserRole
