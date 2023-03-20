package model

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type User struct {
	ID          string  `bson:"_id" json:"id"`
	DisplayName *string `bson:"displayName" json:"displayName"`
	Email       string  `bson:"email" json:"email"`
	Password    string  `bson:"password" json:"-"`
}

func (u *User) HashPassword(password string) string {
	salt := ""
	passwordInput := fmt.Sprintf("%v%v", password, salt)
	hash := sha256.Sum256([]byte(passwordInput))
	hashString := hex.EncodeToString(hash[:])
	return hashString
}
