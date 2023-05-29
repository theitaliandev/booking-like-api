package types

import (
	"fmt"
	"net/mail"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost      = 12
	minFirstNameLen = 2
	minLastNameLen  = 2
	minPasswordLen  = 7
)

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirtsName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"encryptedPassword" json:"-"`
}

type CreateUserParams struct {
	FirtsName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UpdateUserParams struct {
	FirtsName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

func (params CreateUserParams) Validate() map[string]error {
	errors := map[string]error{}
	if len(params.FirtsName) < minFirstNameLen {
		errors["firstName"] = fmt.Errorf("first name length should be at least %d characters", minFirstNameLen)
	}
	if len(params.LastName) < minLastNameLen {
		errors["lastName"] = fmt.Errorf("last name length should be at least %d characters", minLastNameLen)
	}
	if len(params.Password) < minPasswordLen {
		errors["password"] = fmt.Errorf("password length should be at least %d characters", minPasswordLen)
	}
	_, err := mail.ParseAddress(params.Email)
	if err != nil {
		errors["email"] = fmt.Errorf("%s", err.Error())
	}
	return errors
}

func (params UpdateUserParams) ValidateUpdateUser() map[string]error {
	errors := map[string]error{}
	if len(params.FirtsName) < minFirstNameLen {
		errors["firstName"] = fmt.Errorf("first name length should be at least %d characters", minFirstNameLen)
	}
	if len(params.LastName) < minLastNameLen {
		errors["lastName"] = fmt.Errorf("last name length should be at least %d characters", minLastNameLen)
	}
	_, err := mail.ParseAddress(params.Email)
	if err != nil {
		errors["email"] = fmt.Errorf("%s", err.Error())
	}
	return errors
}

func NewUserFromParams(userParams *CreateUserParams) (*User, map[string]error, error) {
	errors := userParams.Validate()
	if len(errors) > 0 {
		return nil, errors, nil
	}

	epw, err := bcrypt.GenerateFromPassword([]byte(userParams.Password), bcryptCost)
	if err != nil {
		return nil, nil, err
	}
	newUser := &User{
		FirtsName:         userParams.FirtsName,
		LastName:          userParams.LastName,
		Email:             userParams.Email,
		EncryptedPassword: string(epw),
	}
	return newUser, nil, nil
}
