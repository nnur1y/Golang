package main

import (
	"Golang/controller"
	"Golang/models"
	"Golang/view"
	"encoding/gob"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"

	_ "github.com/go-sql-driver/mysql"
)

var store = sessions.NewCookieStore([]byte("super-secret"))

func init() {
	store.Options.HttpOnly = true //since we are not accessing ANY COOKIES W/ JavaScript, set to true
	store.Options.Secure = true   // requires secure HTTPS connection
	gob.Register(&models.User{})
}

type RecipePageData struct {
	RecipeData  []models.Recipe
	CommentList []models.CommentData
}

func main() {
	var e error

	if e != nil {
		fmt.Println(e)
		return
	}
	router := gin.Default()

	router.Static("/assets/", "assets/")
	router.LoadHTMLGlob("templates/*.html")
	authRouter := router.Group("/user", auth)

	router.GET("/", handlerIndex)
	router.GET("/products/search", handleSearch)
	router.GET("/products/breakfast", handleProductBreakfast)
	router.GET("/products/salads", handleProductSalad)
	router.GET("/registration", handlerRegistration)
	router.GET("/authorization", handlerAuthorization)
	router.GET("/recipe/:id", handlerRecipe)
	router.POST("/user/reg", handlerUserRegistration)
	router.POST("/user/auth", handlerUserAuthorization)
	router.POST("/sendComment", handleSendComment)
	authRouter.GET("/profile", profileHandler)
	authRouter.GET("/logout", logoutHandler)

	_ = router.Run(":8080")
}

func auth(c *gin.Context) {
	controller.Auth(c, store)
}

var RecipesList []models.Recipe

func handleSendComment(c *gin.Context) {
	controller.SendComment(c)
}

func handlerRecipe(c *gin.Context) {
	view.Recipe(c)
}
func handleProductBreakfast(c *gin.Context) {
	view.GetCategory(c, "Breakfast")
}
func handleProductSalad(c *gin.Context) {
	view.GetCategory(c, "salads")
}
func handleSearch(c *gin.Context) {
	view.Search(c)
}

func handlerIndex(c *gin.Context) {
	view.MainPage(c)
}
func handlerRegistration(c *gin.Context) {
	c.HTML(200, "registration.html", gin.H{})
}

func handlerAuthorization(c *gin.Context) {
	c.HTML(200, "authorization.html", gin.H{})
}

func handlerUserRegistration(c *gin.Context) {
	controller.UserRegistration(c)
}

func handlerUserAuthorization(c *gin.Context) {
	controller.UserAuthorization(c, store)
}

func profileHandler(c *gin.Context) {
	controller.Profile(c, store)
}

func logoutHandler(c *gin.Context) {
	controller.Logout(c, store)
}
