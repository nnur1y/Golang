package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"

	_ "github.com/go-sql-driver/mysql"
)

var router *gin.Engine
var store = sessions.NewCookieStore([]byte("super-secret"))

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
	authRouter := router.Group("/user", auth)

	router.GET("/", handlerIndex)
	router.GET("/index", handlerIndex)
	router.GET("/index/more", handlerIndexMore)
	router.GET("/registration", handlerRegistration)
	router.GET("/authorization", handlerAuthorization)
	router.POST("/user/reg", handlerUserRegistration)
	router.POST("/user/auth", handlerUserAuthorization)
	authRouter.GET("/profile", profileHandler)

	_ = router.Run(":8080")
}

func auth(c *gin.Context) {
	fmt.Println("auth middleware running")
	session, _ := store.Get(c.Request, "session")
	fmt.Println("session:", session)
	_, ok := session.Values["user"]
	if !ok {
		c.HTML(http.StatusForbidden, "registration.html", nil)
		c.Abort()
		return
	}
	fmt.Println("middleware done")
	c.Next()
}

func profileHandler(c *gin.Context) {
	session, _ := store.Get(c.Request, "session")
	var user = &User{}
	val := session.Values["user"]
	var ok bool
	if user, ok = val.(*User); !ok {
		fmt.Println("was not of type *User")
		c.HTML(http.StatusForbidden, "authorization.html", nil)
		return
	}
	c.HTML(http.StatusOK, "profile.html", gin.H{"user": user})
}

// pkg.go.dev/text/template
func handlerIndexMore(c *gin.Context) {
	db, _ := sql.Open("mysql", "root:password@(localhost:3306)/world?parseTime=true")
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
	db, _ := sql.Open("mysql", "root:password@(localhost:3306)/world?parseTime=true")
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
	fmt.Println(e)
	if e != nil {
		c.JSON(200, gin.H{
			"Error": "Не удалось зарегистрировать пользователя",
		})
		c.Redirect(http.StatusFound, "/authorization")

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
			//incorrect email or password
			c.HTML(200, "authorization.html", gin.H{"message": "incorrect username or password"})
		} else {
			session, _ := store.Get(c.Request, "session")
			// session struct has field Values map[interface{}]interface{}
			session.Values["user"] = user
			// save before writing to response/return from handler
			session.Save(c.Request, c.Writer)

			// correct email and password
			fmt.Println(session)
			c.HTML(200, "loggedin.html", gin.H{"email": user.Email})

		}
	}
	c.HTML(200, "authorization.html", gin.H{
		"login":   true,
		"message": "incorrect username or password",
	})
}

type User struct {
	Email     string `json:"Email"`
	Password  string `json:"Password"`
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
}

// Create создание нового пользователя в базе
func (u User) Create() error {
	{ // Insert a new user

		email := u.Email
		password := u.Password
		firstname := u.FirstName
		lastname := u.LastName
		createdAt := time.Now()

		db, _ := sql.Open("mysql", "root:password@(localhost:3306)/world?parseTime=true")
		result, err := db.Exec(`INSERT INTO users (email, password, firstname, lastname, created_at) VALUES (?, ?, ?, ?, ?)`, email, password, firstname, lastname, createdAt)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(result)

		fmt.Println(email, password, firstname, lastname, createdAt)
	}
	return nil
}

func (u User) Select() error {
	{ // Insert a new user

		email := u.Email
		password := u.Password

		var (
			getemail    string
			getpassword string
		)
		db, err := sql.Open("mysql", "root:password@(localhost:3306)/world?parseTime=true")
		query := "SELECT email, password FROM users WHERE email = ? and password = ?"
		if err := db.QueryRow(query, email, password).Scan(&getemail, &getpassword); err != nil {
			fmt.Println("Incorrect email or password")
			return err
		}

		fmt.Println("Correct, Logged in")
		return err
	}
}
