package main

type User struct {
	Id       int    `json:"id"`
	username string `json:"username"`
}

type Users []User
