// Package auth/store.go
package auth

import (
	"errors"

	"github.com/zalando/go-keyring"
)

const (
	service     = "mycli"
	usernameKey = "username"
	tokenKey    = "token"
)

type Credentials struct {
	Username string
	Token    string
}

func Save(creds Credentials) error {
	if err := keyring.Set(service, tokenKey, creds.Token); err != nil {
		return err
	}
	return keyring.Set(service, usernameKey, creds.Username)
}

func Load() (Credentials, error) {
	token, err := keyring.Get(service, tokenKey)
	if err != nil {
		return Credentials{}, errors.New("not logged in, run login")
	}

	username, err := keyring.Get(service, usernameKey)
	if err != nil {
		return Credentials{}, errors.New("not logged in, run login")
	}

	return Credentials{Username: username, Token: token}, nil
}

func Delete() error {
	if err := keyring.Delete(service, tokenKey); err != nil {
		return err
	}
	return keyring.Delete(service, usernameKey)
}
