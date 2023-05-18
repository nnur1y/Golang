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
	router.GET("/", welcomePage)
	router.GET("/recipes", handlerIndex)
	router.GET("/products/search", handleSearch)
	router.GET("/signup", signupView)
	router.GET("/category", category)
	router.GET("/category/:title", categoryRecipes)
	router.GET("/login", loginView)
	router.GET("/recipe/:id", handlerRecipe)
	router.POST("/addRecipe", addRecipe)
	router.POST("/user/reg", handlerUserRegistration)
	router.POST("/user/auth", handlerUserAuthorization)
	router.POST("/sendFeedback", sendFeedback)
	// authRouter.GET("/profile", profileHandler)
	router.GET("/user/logout", logoutHandler)
	router.GET("/profile", profileHandler)
	_ = router.Run(":8080")
}

var RecipesList []models.Recipe

func categoryRecipes(c *gin.Context) {
	title := c.Param("title")
	var recipeList []models.Recipe
	recipeList = view.GetCategory(c, title)
	session, _ := store.Get(c.Request, "session")
	user := session.Values["user"].(*models.User)

	c.HTML(200, "recipes.html", gin.H{
		"content":      recipeList,
		"categoryName": title,
		"username":     user.Username,
	})

}
func category(c *gin.Context) {

	db, err := config.LoadDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Query the "categories_all" table
	rows, err := db.Query("SELECT Categories, imgurl FROM categories_all")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Create an empty array to store the retrieved categories
	categories := []models.Category{}

	// Iterate over the query results and populate the categories array
	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.Title, &category.ImgURL)
		if err != nil {
			log.Fatal(err)
		}
		categories = append(categories, category)
	}

	// Check for any errors during the iteration
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	session, _ := store.Get(c.Request, "session")
	user := session.Values["user"].(*models.User)
	// Send the categories array to the front-end
	c.HTML(200, "category.html", gin.H{
		"categories": categories,
		"username":   user.Username,
	})
}
func sendFeedback(c *gin.Context) {
	session, _ := store.Get(c.Request, "session")
	user := session.Values["user"].(*models.User)
	var sendCom models.SendCommentData

	sendCom.UserId = user.Id
	e := c.BindJSON(&sendCom)
	fmt.Println(sendCom)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": e.Error(),
		})
		return
	}

	fmt.Println("comment 1")
	e = sendCom.Send()

	if e != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": "Не удалось зарегистрировать пользователя",
		})
		return
	}

	c.Redirect(http.StatusFound, "/")
}

func handlerRecipe(c *gin.Context) {
	view.Recipe(c, store)
}

func handleSearch(c *gin.Context) {
	view.Search(c, store)

}

func handlerIndex(c *gin.Context) {
	view.MainPage(c, store)
}
func welcomePage(c *gin.Context) {
	session, _ := store.Get(c.Request, "session")
	user := session.Values["user"].(*models.User)
	c.HTML(200, "index.html", gin.H{"username": user.Username})
}
func signupView(c *gin.Context) {
	c.HTML(200, "signup.html", gin.H{})
}

func loginView(c *gin.Context) {
	c.HTML(200, "login.html", gin.H{})
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
