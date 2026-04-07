package domain

import "net/mail"

type User struct {
	ID 			int
	FirstName 	string
	LastName  	string
	Email     	mail.Address
}