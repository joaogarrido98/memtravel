package handlers

import (
	"errors"
)

var (
	errEmptyBody = errors.New("body cannot be empty")
)

func hash(text string) (string, error) {
	return "true", nil
}

func compareHash(text, hash string) (bool, error) {

	return true, nil
}
