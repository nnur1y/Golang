package main

import (
	"Golang/config"
	"Golang/controller"
	"Golang/models"
	"Golang/view"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"

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
	router.GET("/products/breakfast", handleProductBreakfast)
	router.GET("/products/salads", handleProductSalad)
	router.GET("/signup", signupView)
	router.GET("/addrecipe", addrecipe)
	router.POST("/addrecipe", addrecipe)
	router.GET("/category", category)
	router.GET("/category/:title", categoryRecipes)
	router.GET("/login", loginView)
	router.GET("/buy", buyPage)
	router.POST("/payment", Payment)
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

func buyPage(c *gin.Context) {

	c.HTML(200, "buy.html", gin.H{})
}

func Payment(c *gin.Context) {
	var payment models.Charge

	c.BindJSON(&payment)

	amount, err := strconv.ParseFloat(payment.Amount, 64)
	if err != nil {
		fmt.Print("Invalid payment amount")
		return
	}
	fmt.Println("payment.Amount:", payment)
	apiKey := "sk_test_51N8LwQLBQrS0LN3Qn57GnQluvgAJdaFAXPQQvJfhzJ0ZpQfmDoIf1sgM6E63zC7H0lkZPXP8Ufedv03HkSCLX8CW00tfo0PsTp"
	stripe.Key = apiKey
	_, err = charge.New(&stripe.ChargeParams{
		Amount:      stripe.Int64(int64(amount * 100)),
		Currency:    stripe.String(string(stripe.CurrencyUSD)),
		Description: stripe.String(payment.RecipeId),
		Source:      &stripe.SourceParams{Token: stripe.String("tok_visa")},
	})
	if err != nil {
		c.String(http.StatusBadRequest, "Payment Unsuccessful")
		return
	}
	db, err := config.LoadDB()
	stmt, err := db.Prepare("INSERT INTO payments (recipe_id, user_id) VALUES (?, ?)")
	if err != nil {
		fmt.Print(err)
	}
	defer stmt.Close()
	session, _ := store.Get(c.Request, "session")
	user := session.Values["user"].(*models.User)

	// Execute the SQL statement with the payment data
	_, err = stmt.Exec(payment.RecipeId, user.Id)
	if err != nil {
		fmt.Print(err)
	}
	c.Redirect(http.StatusFound, "/")
}

// func SavePayment(charge *models.Charge) error {

// }

func addrecipe(c *gin.Context) {
	if c.Request.Method == "GET" {
		session, _ := store.Get(c.Request, "session")
		user := session.Values["user"].(*models.User)

		c.HTML(http.StatusOK, "addrecipe.html", gin.H{
			"username": user.Username,
		})
	} else if c.Request.Method == "POST" {

		// Parse the multipart form
		err := c.Request.ParseMultipartForm(32 << 20) // 32MB is the maximum form size
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to parse form",
			})
			return
		}

		// Retrieve the uploaded file from the "image" field
		file, handler, err := c.Request.FormFile("imagefile")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		defer file.Close()
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to read file data",
			})
			return
		}

		folderPath := "assets/img"
		filePath := filepath.Join(folderPath, handler.Filename)
		absPath, err := filepath.Abs(filePath)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to get absolute file path",
			})
			return
		}
		filePath = strings.ReplaceAll(filePath, "\\", "/")
		err = ioutil.WriteFile(absPath, fileBytes, 0644)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to save file to disk",
			})
			return

		}
		var r models.Recipe
		fmt.Println("Image saved:", filePath)
		r.Name = c.Request.FormValue("name")
		r.Description = c.Request.FormValue("description")
		r.ImgURL = "/" + filePath
		r.Categories = c.Request.FormValue("category")
		db, err := config.LoadDB()
		_, err = db.Exec("INSERT INTO recipe (name, description, image_r, Categories) VALUES (?, ?, ?, ?)", &r.Name, &r.Description, &r.ImgURL, &r.Categories)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		c.Redirect(http.StatusFound, "/")
	}
}
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
