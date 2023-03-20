package util

import (
	"context"
	"errors"
	"strings"

	"github.com/tipbk/doodle/model"
)

func GetUserFromContext(ctx context.Context) (model.User, error) {
	user, ok := ctx.Value("user").(model.User)
	if !ok {
		return model.User{}, errors.New("auth failed")
	}
	return user, nil
}

func GetTokenString(header string) string {
	if header == "" {
		return ""
	}
	BEARER := "Bearer "
	if len(header) <= len(BEARER) {
		return ""
	}

	result := strings.Split(header, BEARER)
	return result[1]
}
