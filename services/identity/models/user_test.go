package models

import "testing"

func TestValidateFields(t *testing.T) {
	tests := []struct {
		user     *User
		errorMsg string
	}{
		{
			user: &User{
				Email:     "test@yahoo.com",
				Password:  "123456",
				FirstName: "John",
				LastName:  "Doe",
			},
		},
		{
			user: &User{
				Password:  "123456",
				FirstName: "John",
				LastName:  "Doe",
			},
			errorMsg: "invalid email",
		},
		{
			user: &User{
				Email:     "test@yahoo.com",
				FirstName: "John",
				LastName:  "Doe",
			},
			errorMsg: "invalid password",
		},
		{
			user: &User{
				Email:     "test@yahoo.com",
				Password:  "213",
				FirstName: "John",
				LastName:  "Doe",
			},
			errorMsg: "invalid password",
		},
		{
			user: &User{
				Email:    "test@yahoo.com",
				Password: "123456",
				LastName: "Doe",
			},
			errorMsg: "invalid first name",
		},
		{
			user: &User{
				Email:     "test@yahoo.com",
				Password:  "123456",
				FirstName: "John",
			},
			errorMsg: "invalid last name",
		},
	}

	for _, test := range tests {
		gotErrorMsg := ""
		if err := test.user.ValidateFields(); err != nil {
			gotErrorMsg = err.Error()
		}
		if gotErrorMsg != test.errorMsg {
			t.Errorf("ValidateFields(%s) is incorrect, got: %s, want: %s.", test.user, gotErrorMsg, test.errorMsg)
		}
	}
}

func TestIsEmailValid(t *testing.T) {
	tests := []struct {
		email string
		want  bool
	}{
		{"test@yahoo.com", true},
		{"a@test.ru", true},
		{"test@yahoo", false},
		{"testyahoo.com", false},
		{"@yahoo.com", false},
		{"@.com", false},
		{"dsaads@.com", false},
	}

	for _, test := range tests {
		got := isEmailValid(test.email)
		if got != test.want {
			t.Errorf("isEmailValid(%s) is incorrect, got: %t, want: %t.", test.email, got, test.want)
		}
	}
}

func TestHashPassword(t *testing.T) {
	pass := "123456"
	u := &User{
		Password: pass,
	}
	if err := u.HashPassword(); err != nil {
		t.Errorf("HashPassword() got error: %t", err)
	}
	if u.Password == pass || u.Password == "" {
		t.Errorf("HashPassword() did not hash password correct: %s => %s", pass, u.Password)
	}
}
