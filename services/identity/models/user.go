package models

import (
	pb "ebank/pb/services/identity"
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string
	Email     string
	Password  string
	FirstName string
	LastName  string
	CreatedAt string
}

func (u *User) HashPassword() error {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPass)
	return nil
}

func (u *User) ComparePassword(pass string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pass))
}

func (u *User) ValidateFields() error {
	// Validate email
	if !isEmailValid(u.Email) {
		return errors.New("invalid email")
	}
	// Validate password
	if len(u.Password) < 6 || len(u.Password) > 254 {
		return errors.New("invalid password")
	}
	// Validate first name
	if len(u.FirstName) == 0 {
		return errors.New("invalid first name")
	}
	// Validate last name
	if len(u.LastName) == 0 {
		return errors.New("invalid last name")
	}
	return nil
}

var emailRegex = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

func isEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}

func (u *User) ToProto() *pb.User {
	return &pb.User{
		Id:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}
}
