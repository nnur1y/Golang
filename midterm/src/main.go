package main

import (
	"Golang/config"
	"Golang/controller"
	"Golang/models"
	"Golang/view"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"

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
	router.POST("/addRecipe", addRecipe)
	router.POST("/user/reg", handlerUserRegistration)
	router.POST("/user/auth", handlerUserAuthorization)
	router.POST("/sendComment", handleSendComment)
	// authRouter.GET("/profile", profileHandler)
	authRouter.GET("/logout", logoutHandler)
	router.GET("/profile", profileHandler)
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
	view.Search(c, store)

}

func handlerIndex(c *gin.Context) {
	view.MainPage(c, store)
}
func handlerRegistration(c *gin.Context) {
	c.HTML(200, "registration.html", gin.H{})
}

func handlerAuthorization(c *gin.Context) {
	c.HTML(200, "authorization.html", gin.H{})
}

func handlerUserRegistration(c *gin.Context) {
	store = controller.UserRegistration(c, store)
}

func handlerUserAuthorization(c *gin.Context) {
	store = controller.UserAuthorization(c, store)
}

func profileHandler(c *gin.Context) {
	controller.Profile(c, store)
}

func logoutHandler(c *gin.Context) {
	controller.Logout(c, store)
}
func addRecipe(c *gin.Context) {
	// name,
	// description,
	// cooking_time,
	// category
	// Get recipe details from form data
	name := c.PostForm("name")
	description := c.PostForm("description")
	coooking_time := c.PostForm("cooking_time")
	categories := c.PostForm("category")
	db, _ := config.LoadDB()
	// Insert recipe into database
	result, err := db.Exec(
		"INSERT INTO `recipe`( `name`, `description`,  `coooking_time`, `categories`) VALUES (?,?,?,?)", name, description, coooking_time, categories)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add recipe"})
		return
	}
	// Get ID of newly inserted recipe
	recipeID, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get recipe ID"})
		return
	}

	var user models.User

	session, err := store.Get(c.Request, "session")
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println("session:", session)
	userVal, ok := session.Values["user"]
	user = *userVal.(*models.User)
	if !ok {
		fmt.Println("no user", userVal)
	}

	userrecipe, err := db.Exec("INSERT INTO `user_recipes`( `recipe_id`, `user_id`) VALUES (?,?)", recipeID, user.Id)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add recipe"})
		return
	}
	fmt.Println(userrecipe)
	// Return success response with ID of newly inserted recipe
	c.JSON(http.StatusOK, gin.H{"message": "Recipe added successfully", "recipe_id": recipeID})

}
