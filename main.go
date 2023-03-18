package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	_ "github.com/go-sql-driver/mysql"
)

var router *gin.Engine

type Recipe struct {
	Id         int
	Name       string
	Definition string
	Author     string
}

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
	router.GET("/index", handlerIndex)
	router.GET("/index/more", handlerIndexMore)
	router.GET("/registration", handlerRegistration)
	router.GET("/authorization", handlerAuthorization)
	router.POST("/user/reg", handlerUserRegistration)
	router.POST("/user/auth", handlerUserAuthorization)
	_ = router.Run(":8080")
}

// pkg.go.dev/text/template
func handlerIndexMore(c *gin.Context) {
	db, err := sql.Open("mysql", "root:password@(localhost:3306)/world?parseTime=true")
	result, err := db.Query("SELECT  * from Recipe where id_r=1")
	if err != nil {
		log.Println(err)
	}
	result.Next()
	recipe := Recipe{}
	var id int
	var name string
	var definition string
	var author string

	err = result.Scan(&id, &name, &definition, &author)

	recipe.Id = id
	recipe.Name = name
	recipe.Definition = definition
	recipe.Author = author
	if err != nil {
		panic(err)
	}
	c.HTML(200, "layout.html", gin.H{
		"search":  true,
		"content": recipe,
	})
}
func handlerIndex(c *gin.Context) {
	db, err := sql.Open("mysql", "root:password@(localhost:3306)/world?parseTime=true")
	result, err := db.Query("SELECT  * from Recipe ")
	if err != nil {
		log.Println(err)
	}
	recipe := Recipe{}
	recipes := []Recipe{}

	for result.Next() {
		var id int
		var name string
		var definition string
		var author string

		err = result.Scan(&id, &name, &definition, &author)

		recipe.Id = id
		recipe.Name = name
		recipe.Definition = definition
		recipe.Author = author

		recipes = append(recipes, recipe)

		if err != nil {
			panic(err)
		}
	}

	// var tmpl = template.Must(template.ParseFiles("./templates/layout.html"))
	// nerr := tmpl.Execute(w, recipes)

	c.HTML(200, "layout.html", gin.H{
		"search":  false,
		"content": recipes,
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

	}

	e = user.Create()
	if e != nil {
		c.JSON(200, gin.H{
			"Error": "Не удалось зарегистрировать пользователя",
		})
		c.Redirect(http.StatusFound, "/")
		// c.Redirect(http.StatusFound, "/foo")

	}

	// c.HTML(200, "layout.html", gin.H{})
	// c.JSON(200, gin.H{
	// 	"Error": nil,
	// })
}

func handlerUserAuthorization(c *gin.Context) {
	var user User
	e := c.BindJSON(&user)
	if e != nil {
		c.JSON(200, gin.H{
			"Error": e.Error(),
		})
		return
	}

	e = user.Select()
	fmt.Println("sign up ")
	if e != nil {
		c.HTML(200, "layout.html", gin.H{
			"Role": "manager",
		})
		return
	}

	// c.JSON(200, gin.H{
	// 	"Error": nil,
	// })
	c.HTML(200, "layout.html", gin.H{
		"Role": "manager",
	})

}

type User struct {
	Login     string `json:"Login"`
	Password  string `json:"Password"`
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
}

// Create создание нового пользователя в базе
func (u User) Create() error {
	{ // Insert a new user

		username := u.Login
		password := u.Password
		firstname := u.FirstName
		lastname := u.LastName
		createdAt := time.Now()

		db, err := sql.Open("mysql", "root:password@(localhost:3306)/world?parseTime=true")
		result, err := db.Exec(`INSERT INTO users (username, password, firstname, lastname, created_at) VALUES (?, ?, ?, ?, ?)`, username, password, firstname, lastname, createdAt)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(result)

		fmt.Println(username, password, firstname, lastname, createdAt)
	}
	return nil
}

func (u User) Select() error {
	{ // Insert a new user

		username := u.Login
		// password := u.Password
		var (
			getusername string
			getpassword string
		)
		db, err := sql.Open("mysql", "root:password@(localhost:3306)/world?parseTime=true")
		// result, err := db.Exec(`select  VALUES (?, ?, ?)`, username, password, createdAt)
		query := "SELECT  username, password FROM users WHERE username = ?"
		if err := db.QueryRow(query, username).Scan(&getusername, &getpassword); err != nil {
			log.Fatal(err)
		}

		fmt.Println(getusername, getpassword)
		if err != nil {
			log.Fatal(err)
		}

		// fmt.Println(result)

		// fmt.Println(username, password, createdAt)
	}
	return nil
}
