package main

import (
	"Golang/config"
	"Golang/models"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"

	_ "github.com/go-sql-driver/mysql"
)

var router *gin.Engine
var store = sessions.NewCookieStore([]byte("super-secret"))

func main() {
	var e error

	if e != nil {
		fmt.Println(e)
		return
	}
	router := gin.Default()

	router.Static("/assets/", "assets/")
	router.LoadHTMLGlob("templates/*.html")
	// authRouter := router.Group("/user", auth)

	router.GET("/", handlerIndex)
	router.GET("/products/search", handleSearch)
	router.GET("/products/breakfast", handleProductBreakfast)
	router.GET("/products/salads", handleProductSalad)
	router.GET("/registration", handlerRegistration)
	router.GET("/authorization", handlerAuthorization)
	router.POST("/user/reg", handlerUserRegistration)
	router.POST("/user/auth", handlerUserAuthorization)
	// authRouter.GET("/profile", profileHandler)

	_ = router.Run(":8080")
}

var RecipesList []models.Recipe

func handleProductBreakfast(c *gin.Context) {
	var searchItem models.SearchItem
	db, _ := config.LoadDB()
	result, err := db.Query("SELECT id_r,name,description,categories FROM recipe WHERE  categories= 'Breakfast' ")
	if err != nil {
		fmt.Print(err)
	}
	RecipesList, e := searchItem.Search(result)
	if e != nil {
		fmt.Print(e)
	}

	c.HTML(200, "layout.html", gin.H{
		"search":  false,
		"content": RecipesList,
	})
}
func handleProductSalad(c *gin.Context) {
	var searchItem models.SearchItem
	db, _ := config.LoadDB()
	result, err := db.Query("SELECT id_r,name,description,categories FROM recipe WHERE  categories='Salads' ")
	if err != nil {
		fmt.Print(err)
	}
	RecipesList, e := searchItem.Search(result)
	if e != nil {
		fmt.Print(e)
	}

	c.HTML(200, "layout.html", gin.H{
		"search":  false,
		"content": RecipesList,
	})
}
func handleSearch(c *gin.Context) {
	query := c.Query("searchItem")
	fmt.Println(query)
	var searchItem models.SearchItem
	var text = "%" + query + "%"
	db, _ := config.LoadDB()
	result, err := db.Query("SELECT id_r,name,description,categories FROM recipe WHERE name LIKE ? ", text)
	if err != nil {
		fmt.Print(err)
	}
	RecipesList, e := searchItem.Search(result)
	if e != nil {
		fmt.Print(e)
	}

	c.HTML(200, "layout.html", gin.H{
		"search":  false,
		"content": RecipesList,
	})
}

func handlerIndex(c *gin.Context) {

	db, _ := config.LoadDB()
	var searchItem models.SearchItem

	var text = "%%"
	result, err := db.Query("SELECT id_r,name,description,categories FROM recipe WHERE name LIKE ? ", text)
	if err != nil {
		fmt.Print(err)
	}
	RecipesList, e := searchItem.Search(result)

	if e != nil {
		fmt.Println(e)
	}

	c.HTML(200, "layout.html", gin.H{
		"search":  false,
		"content": RecipesList,
	})

}
func handlerRegistration(c *gin.Context) {
	c.HTML(200, "registration.html", gin.H{})
}

func handlerAuthorization(c *gin.Context) {
	c.HTML(200, "authorization.html", gin.H{})
}

func handlerUserRegistration(c *gin.Context) {

	var user models.User
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
		// c.Redirect(http.StatusFound, "/authorization")

	} else {
		fmt.Println("Signed up!")
		c.JSON(200, gin.H{
			"response": "signed up!",
		})
		return
	}

}

func handlerUserAuthorization(c *gin.Context) {
	var user models.User
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
