package main

import (
	"github.com/a8m/documentdb"
)

type User struct {
	documentdb.Document
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

type DB interface {
	Get(id string) *User
	GetAll() []*User
	Add(u *User) *User
	Update(u *User) *User
	Delete(id string) error
}
