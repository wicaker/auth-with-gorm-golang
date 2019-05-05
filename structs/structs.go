package structs

import "github.com/jinzhu/gorm"

// User struct contain email and password
type User struct {
	gorm.Model
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Mail struct to send email
type Mail struct {
	senderID string
	toIds    []string
	subject  string
	body     string
}

// SMPTServer to set host
type SMPTServer struct {
	host string
	port string
}
