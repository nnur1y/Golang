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

type Recipe struct {
	Id         int
	Name       string
	Definition string
	Author     string
}

var searchText string

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
	// router.GET("/search", handleSearch)
	router.POST("/index/more", handlerIndexMore)
	router.GET("/index/more", handleSearch)
	router.GET("/registration", handlerRegistration)
	router.GET("/authorization", handlerAuthorization)
	router.POST("/user/reg", handlerUserRegistration)
	router.POST("/user/auth", handlerUserAuthorization)
	_ = router.Run(":8080")
}

var RecipesList []Recipe

func handleSearch(c *gin.Context) {
	c.HTML(200, "layout.html", gin.H{
		"error":   false,
		"content": RecipesList,
	})
}

// pkg.go.dev/text/template
func handlerIndexMore(c *gin.Context) {

	var searchItem SearchItem
	e := c.BindJSON(&searchItem)
	if e != nil {
		fmt.Print(e)
	}
	e = searchItem.Search()
	if e != nil {
		fmt.Println("Logged in!")
		c.HTML(200, "layout.html", gin.H{
			"error":   false,
			"content": RecipesList,
		})
		c.Request.Response.Location()
		return
	} else {

		c.HTML(200, "layout.html", gin.H{
			"error":   true,
			"content": RecipesList,
		})
		return
	}
}
func handlerIndex(c *gin.Context) {

	db, _ := sql.Open("mysql", "root:password@(localhost:3306)/world?parseTime=true")
	result, err := db.Query("SELECT  * from Recipe")
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
	// fmt.Println(recipes)
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

		// c.Redirect(http.StatusFound, "/foo")

	} else {
		fmt.Println("Signed up!")
	}

}

func handlerUserAuthorization(c *gin.Context) {
	var user User
	e := c.BindJSON(&user)

	if e != nil {
		c.JSON(200, gin.H{
			"Error": e.Error(),
		})
		fmt.Println("Некорректные данные")

	} else {
		e = user.Select()

		if e != nil {
			fmt.Println("Logged in!")
			c.HTML(200, "layout.html", gin.H{
				"error": false,
			})
			return
		} else {

			c.HTML(200, "layout.html", gin.H{
				"error": true,
			})
			return
		}
	}

}

type User struct {
	Login     string `json:"Login"`
	Password  string `json:"Password"`
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
}

type SearchItem struct {
	SearchItem string `json:"searchItem"`
}

func (s SearchItem) Search() error {

	searchText = "%" + s.SearchItem + "%"
	db, _ := sql.Open("mysql", "root:password@(localhost:3306)/world?parseTime=true")

	result, err := db.Query("SELECT * FROM recipe WHERE name LIKE ? ", searchText)
	// result, err2 := db.Query("SELECT  username, password FROM users WHERE username = ?", username)
	if err != nil {
		panic(err)
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
	fmt.Println(recipes)
	RecipesList = recipes
	return nil

}

// Create создание нового пользователя в базе
func (u User) Create() error {
	{ // Insert a new user

		username := u.Login
		password := u.Password
		firstname := u.FirstName
		lastname := u.LastName
		createdAt := time.Now()

		db, _ := sql.Open("mysql", "root:password@(localhost:3306)/world?parseTime=true")
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
		var err2 error
		username := u.Login
		// password := u.Password
		var (
			getusername string
			getpassword string
		)
		db, err := sql.Open("mysql", "root:password@(localhost:3306)/world?parseTime=true")
		if err != nil {
			fmt.Println(err)
		}
		// result, err := db.Exec(`select  VALUES (?, ?, ?)`, username, password, createdAt)
		// query := "SELECT  username, password FROM users WHERE username = ?"
		result, err2 := db.Query("SELECT  username, password FROM users WHERE username = ?", username)
		result.Next()
		err2 = result.Scan(&getusername, &getpassword)
		// if err := db.QueryRow(query, username).Scan(&getusername, &getpassword); err != nil {
		// log.Fatal(err)
		// }
		if err2 != nil {
			return err2
			fmt.Print("error")
		} else {
			fmt.Println(result, getusername, getpassword)
		}

		// fmt.Println(result)

		// fmt.Println(username, password, createdAt)
	}
	return nil
}
