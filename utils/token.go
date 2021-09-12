package utils

import "github.com/google/uuid"

func NewToken() string {
	return uuid.New().String()
}
