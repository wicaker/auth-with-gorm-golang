package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/wicaker/go_auth/config"
	"github.com/wicaker/go_auth/controllers"
)

func main() {
	db := config.DBInit()
	inDB := &controllers.InDB{DB: db}

	router := gin.Default()

	router.POST("/register", inDB.RegisterUser)
	router.POST("/login", inDB.LoginUser)
	router.Run(":4000")
}
