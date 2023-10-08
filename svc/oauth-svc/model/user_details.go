package model

import (
	"strings"
)

type UserDetails struct {
	UserId      int64
	Username    string
	Password    string
	Authorities []string
}

func (userDetails *UserDetails) IsMatch(username string, password string) bool {
	return strings.EqualFold(username, userDetails.Username) && strings.EqualFold(password, userDetails.Password)
}
