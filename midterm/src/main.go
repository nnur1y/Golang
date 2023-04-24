package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"github.com/nnur1y/Golang/tree/main/midterm/src/config"
	"github.com/nnur1y/Golang/tree/main/midterm/src/models"

	_ "github.com/go-sql-driver/mysql"
)

var router *gin.Engine
var store = sessions.NewCookieStore([]byte("super-secret"))

type CommentData struct {
	Username string
	ComText  string
}
type SendCommentData struct {
	CommentText string
	UserId      string
	RecipeId    string
	Rate        string
}
type RecipePageData struct {
	RecipeData  []models.Recipe
	CommentList []CommentData
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
	// authRouter := router.Group("/user", auth)

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
	// authRouter.GET("/profile", profileHandler)

	_ = router.Run(":8080")
}

var RecipesList []models.Recipe

func handleSendComment(c *gin.Context) {
	var sendCom SendCommentData
	e := c.BindJSON(&sendCom)

	if e != nil {
		c.JSON(200, gin.H{
			"Error": e.Error(),
		})

	}
	fmt.Println("comment 1")
	e = sendCom.SendComment()

	if e != nil {
		c.JSON(200, gin.H{
			"Error": "Не удалось зарегистрировать пользователя",
		})
		// c.Redirect(http.StatusFound, "/authorization")

	} else {
		fmt.Println("Comment sent")
		c.JSON(200, gin.H{
			"response": "Comment sent",
		})
		return
	}

}

func (s SendCommentData) SendComment() error {

	{ // Insert a new user

		commentText := s.CommentText
		userId := s.UserId
		recipeId := s.RecipeId
		rate := s.Rate

		db, _ := config.LoadDB()
		result, err := db.Exec(`INSERT INTO feedback( userid, recipeid, comment,rate) VALUES (?, ?, ?,?)`, userId, recipeId, commentText, rate)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(result)

	}
	return nil

}

func handlerRecipe(c *gin.Context) {
	recipeid := c.Param("id")
	// fmt.Println("userid " + recipeid)

	var searchItem models.SearchItem
	db, _ := config.LoadDB()
	result, err := db.Query("select rc.Id_r, rc.name, rc.description ,rc.categories ,rt.rating from recipe rc join ratings rt on rc.ID_r = rt.recipeId WHERE  rc.id_r= ?", recipeid)
	if err != nil {
		fmt.Print(err)
	}
	RecipesList, e := searchItem.Search(result)

	if e != nil {
		fmt.Print(e)
	}
	commentsResult, err := db.Query("SELECT u.Username, c.comment FROM comments c join users u on u.id =c.userid WHERE recipeid = ?", recipeid)
	if err != nil {
		fmt.Print(err)
	}
	comment := CommentData{}
	commentsList := []CommentData{}

	for commentsResult.Next() {
		var username string
		var comText string

		err = commentsResult.Scan(&username, &comText)

		comment.Username = username
		comment.ComText = comText

		commentsList = append(commentsList, comment)

		if err != nil {
			panic(err)
		}

	}
	fmt.Println(commentsList)

	c.HTML(200, "singleRecipe.html", gin.H{
		"recipeData":   RecipesList,
		"commentsData": commentsList,
	})

}
func handleProductBreakfast(c *gin.Context) {
	var searchItem models.SearchItem
	db, _ := config.LoadDB()
	result, err := db.Query("select rc.Id_r, rc.name, rc.description ,rc.categories ,rt.rating from recipe rc join ratings rt on rc.ID_r = rt.recipeId WHERE  rc.categories= 'Breakfast' order by rt.rating desc ")
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
	result, err := db.Query("select rc.Id_r, rc.name, rc.description ,rc.categories ,rt.rating from recipe rc join ratings rt on rc.ID_r = rt.recipeId WHERE  rc.categories='Salads'  order by rt.rating desc")
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
	result, err := db.Query("select rc.Id_r, rc.name, rc.description ,rc.categories ,rt.rating from recipe rc join ratings rt on rc.ID_r = rt.recipeId where rc.name LIKE ? ", text)
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
	result, err := db.Query("select rc.Id_r, rc.name, rc.description ,rc.categories ,rt.rating from recipe rc join ratings rt on rc.ID_r = rt.recipeId WHERE rc.name LIKE ? ", text)
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
