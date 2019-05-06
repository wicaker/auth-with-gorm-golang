package structs

import "github.com/jinzhu/gorm"

// User struct contain email and password
type User struct {
	gorm.Model
	Email    string `json:"email"`
	Password string `json:"password"`
}
