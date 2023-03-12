package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	_ "github.com/go-sql-driver/mysql"
)

var router *gin.Engine

var database *sql.DB

func main() {
	var e error

	if e != nil {
		fmt.Println(e)
		return
	}

	router = gin.Default()
	router.Static("/assets/", "assets/")
	router.LoadHTMLGlob("templates/*.html")
	router.GET("/", handlerIndex)
	router.GET("/registration", handlerRegistration)
	router.GET("/authorization", handlerAuthorization)
	router.POST("/user/reg", handlerUserRegistration)
	router.POST("/user/auth", handlerUserAuthorization)
	_ = router.Run(":8080")
}

// pkg.go.dev/text/template
func handlerIndex(c *gin.Context) {
	c.HTML(200, "layout.html", gin.H{
		"Role": "manager",
	})
}
func handlerRegistration(c *gin.Context) {
	c.HTML(200, "registration.html", gin.H{})
}

func handlerAuthorization(c *gin.Context) {
	c.HTML(200, "authorization.html", gin.H{})
}

func handlerUserRegistration(c *gin.Context) {

	var user User
	// fmt.Println(c.)
	e := c.BindJSON(&user)
	if e != nil {
		c.JSON(200, gin.H{
			"Error": e.Error(),
		})
		return
	}

	e = user.Create()
	if e != nil {
		c.JSON(200, gin.H{
			"Error": "Не удалось зарегистрировать пользователя",
		})
		return
	}

	// c.JSON(200, gin.H{
	// 	"Error": nil,
	// })
}

func handlerUserAuthorization(c *gin.Context) {
	// var user User
	// e := c.BindJSON(&user)
	// if e != nil {
	// 	c.JSON(200, gin.H{
	// 		"Error": e.Error(),
	// 	})
	// 	return
	// }

	// e = user.Select()
	// if e != nil {
	// 	c.JSON(200, gin.H{
	// 		"Error": "Не удалось авторизоваться",
	// 	})
	// 	return
	// }

	// c.JSON(200, gin.H{
	// 	"Error": nil,
	// })
}

type User struct {
	Login     string `json:"Login"`
	Password  string `json:"Password"`
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	Role      string `json:"Role,-"`
}

// Create создание нового пользователя в базе
func (u User) Create() error {
	{ // Insert a new user

		username := u.FirstName
		password := u.Password
		createdAt := time.Now()

		// http.Redirect(w, result, "/", 301)
		db, err := sql.Open("mysql", "root:password@(localhost:3306)/world?parseTime=true")
		result, err := db.Exec(`INSERT INTO users (username, password, created_at) VALUES (?, ?, ?)`, username, password, createdAt)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(result)

		fmt.Println(username, password, createdAt)
	}
	return nil
}
