package models

import "gorm.io/gorm"

type Todo struct {
	gorm.Model
	Text   string
	Status uint64
	User   string
}

type User struct {
	gorm.Model
	Account string
	Pass    []byte
	Admin   bool
}
