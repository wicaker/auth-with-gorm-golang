package controllers

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/wicaker/go_auth/structs"
	"golang.org/x/crypto/bcrypt"
)

// ServerName function to be create server name
func (s *SMPTServer) ServerName() string {
	return s.host + ":" + s.port
}

// BuildMessage function to be build messagae thats will be sending to user when register
func (mail *Mail) BuildMessage() string {
	message := ""
	message += fmt.Sprintf("From: %s\r\n", mail.senderID)
	if len(mail.toIds) > 0 {
		message += fmt.Sprintf("To: %s\r\n", strings.Join(mail.toIds, ";"))
	}

	message += fmt.Sprintf("Subject: %s\r\n", mail.subject)
	message += "\r\n" + mail.body

	return message
}

// SendEmail function to send email
func SendEmail(email string) {
	mail := Mail{}
	mail.senderID = os.Getenv("email_account")
	mail.toIds = []string{email}
	mail.subject = "Wellcome " + email
	mail.body = "Hello " + email + ". Thank you for registering. Please wait for next steps"

	messageBody := mail.BuildMessage()

	smtpServer := SMPTServer{host: "smtp.gmail.com", port: "465"}

	log.Println(smtpServer.host)
	//build an auth
	auth := smtp.PlainAuth("", mail.senderID, os.Getenv("email_password"), smtpServer.host)

	// Gmail will reject connection if it's not secure
	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpServer.host,
	}

	conn, err := tls.Dial("tcp", smtpServer.ServerName(), tlsconfig)
	if err != nil {
		log.Panic(err)
	}

	client, err := smtp.NewClient(conn, smtpServer.host)
	if err != nil {
		log.Panic(err)
	}

	// step 1: Use Auth
	if err = client.Auth(auth); err != nil {
		log.Panic(err)
	}

	// step 2: add all from and to
	if err = client.Mail(mail.senderID); err != nil {
		log.Panic(err)
	}
	for _, k := range mail.toIds {
		if err = client.Rcpt(k); err != nil {
			log.Panic(err)
		}
	}

	// Data
	w, err := client.Data()
	if err != nil {
		log.Panic(err)
	}

	_, err = w.Write([]byte(messageBody))
	if err != nil {
		log.Panic(err)
	}

	err = w.Close()
	if err != nil {
		log.Panic(err)
	}

	client.Quit()

	log.Println("Mail sent successfully")
}

//Validate incoming user details...
func Validate(email, password string) string {
	if !strings.Contains(email, "@") {
		return "Email invalid"
	}

	if len(password) < 6 {
		return "Password must more than 6 characters"
	}

	return "pass"
}

// RegisterUser for registering new user
func (idb *InDB) RegisterUser(c *gin.Context) {
	var (
		user   structs.User
		result gin.H
	)

	email := c.PostForm("email")
	password := c.PostForm("password")

	if Validate(email, password) == "pass" {
		err := idb.DB.Where("email = ?", email).First(&user).Error
		if err != nil {
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

			user.Email = email
			user.Password = string(hashedPassword)

			SendEmail(email)
			idb.DB.Create(&user)
			result = gin.H{
				"result": user,
			}
		} else {
			result = gin.H{
				"result": "Email already exist",
				"count":  0,
			}
		}
	} else {
		result = gin.H{
			"result": Validate(email, password),
			"count":  0,
		}
	}
	c.JSON(http.StatusOK, result)
}

// LoginUser for log controllers of user
func (idb *InDB) LoginUser(c *gin.Context) {
	var (
		user   structs.User
		result gin.H
	)

	email := c.PostForm("email")
	password := c.PostForm("password")

	checkEmail := idb.DB.Where("email = ?", email).First(&user).Error
	if checkEmail != nil {
		result = gin.H{
			"result": "email not found",
		}
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		result = gin.H{
			"result": "Invalid login credentials. Please try again",
		}
	} else {
		sign := jwt.New(jwt.GetSigningMethod("HS256"))
		token, err := sign.SignedString([]byte("secret"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			c.Abort()
		}
		result = gin.H{
			"token": token,
		}
	}
	fmt.Println(user)
	c.JSON(http.StatusOK, result)
}
