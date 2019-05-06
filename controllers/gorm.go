package controllers

import "github.com/jinzhu/gorm"

// InDB structs for DB
type InDB struct {
	DB *gorm.DB
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
