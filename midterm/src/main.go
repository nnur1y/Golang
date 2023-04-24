package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"github.com/nnur1y/Golang/tree/main/midterm/src/config"
	"github.com/nnur1y/Golang/tree/main/midterm/src/models"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"
)

var store = sessions.NewCookieStore([]byte("super-secret"))

func init() {
	store.Options.HttpOnly = true //since we are not accessing ANY COOKIES W/ JavaScript, set to true
	store.Options.Secure = true   // requires secure HTTPS connection
	gob.Register(&models.User{})
}

type CommentData struct {
	Username string
	ComText  string
}
type SendCommentData struct {
	CommentText string
	UserId      string
	RecipeId    string
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
	fmt.Println("aith middleware running")
	session, _ := store.Get(c.Request, "session")
	fmt.Println("session:", session)
	_, ok := session.Values["user"]
	if !ok {
		c.HTML(http.StatusForbidden, "authorization.html", nil)
		c.Abort()
		return
	}
	fmt.Println("middleware done")
	c.Next()
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
	// commentText := c.FormValue("commentText")
	// fmt.Println(commentText)
	// if len(commentText) == 0 {
	// 	result = "write some comment"
	// 	tpl.ExecuteTemplate(w, "product.html", result)
	// 	return
	// } else {

	// 	userId := r.FormValue("userId")
	// 	productId := r.FormValue("productId")
	// 	var insertStmt *sql.Stmt
	// 	insertStmt, err2 := database.Prepare("INSERT INTO comments (productid,userid, comment) VALUES (?, ?,?);")
	// 	// fmt.Println(userId)
	// 	if err2 != nil {
	// 		fmt.Println("error preparing statement:", err2)
	// 		tpl.ExecuteTemplate(w, "index.html", "there was a problem registering account")
	// 		return
	// 	}
	// 	defer insertStmt.Close()
	// 	var result sql.Result

	// 	result, err2 = insertStmt.Exec(productId, userId, commentText)
	// 	lastIns, _ := result.LastInsertId()
	// 	fmt.Println("lastIns comment:", lastIns)
	// 	if err2 != nil {
	// 		fmt.Println("error inserting new user")
	// 		tpl.ExecuteTemplate(w, "registration.html", "there was a problem registering account")
	// 		return
	// 	}
	// 	p, err := getProduct(productId)
	// 	fmt.Println(productId)
	// 	if err != nil {
	// 		log.Println(err)
	// 		http.Error(w, "Failed to retrieve product", http.StatusInternalServerError)
	// 		return
	// 	}

	// 	comments, err := getComments(productId)
	// 	if err != nil {
	// 		log.Println(err)
	// 		http.Error(w, "Failed to retrieve comments", http.StatusInternalServerError)
	// 		return
	// 	}

	// 	data := ProductPage{
	// 		Product:  p,
	// 		Comments: comments,
	// 	}

	// 	http.Redirect(w, r, "/product:"+productId, http.StatusSeeOther)
	// 	tpl.ExecuteTemplate(w, "product.html", data)
	// }

}

func (s SendCommentData) SendComment() error {

	{ // Insert a new user

		commentText := s.CommentText
		userId := s.UserId
		recipeId := s.RecipeId

		db, _ := config.LoadDB()
		result, err := db.Exec(`INSERT INTO comments( userid, recipeid, comment) VALUES (?, ?, ?)`, userId, recipeId, commentText)
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
	result, err := db.Query("SELECT id_r,name,description,categories FROM recipe WHERE  id_r= ?", recipeid)
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
	var u models.User
	u.Username = c.PostForm("Username")
	u.Email = c.PostForm("Email")
	u.Password = c.PostForm("Password")

	exists := u.UsernameExists()
	if exists {
		c.HTML(http.StatusBadRequest, "registration.html", gin.H{
			"message": "Username already taken please try another",
		})
		return
	}

	err := u.Create()
	if err != nil {
		fmt.Println("create new error: ", err)
		err = errors.New("there was an issue creating account, please try again")
		c.HTML(http.StatusBadRequest, "registration.html", gin.H{
			"message":  err,
			"username": u.Username,
		})
		return
	}

	/*
		db, _ := config.LoadDB()
		var searchItem models.SearchItem
		result, err := db.Query("SELECT id_r,name,description,categories FROM recipe WHERE name LIKE ? ", text)
		if err != nil {
			fmt.Print(err)
		}
		RecipesList, e := searchItem.Search(result)

		if e != nil {
			fmt.Println(e)
		}
	*/
	c.HTML(http.StatusOK, "layout.html", gin.H{
		"username": u.Username,
		//"content":  RecipesList,
	})
}

func handlerUserAuthorization(c *gin.Context) {
	var user models.User
	user.Username = c.PostForm("Username")
	password := c.PostForm("Password")
	e := user.Select()
	if e != nil {
		fmt.Println("error selecting hashed password in db by Username, err:", e)
		c.HTML(http.StatusUnauthorized, "authorization.html", gin.H{"message": "Incorrect username and password!"})
		return
	}
	e = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if e == nil {
		session, _ := store.Get(c.Request, "session")
		session.Values["user"] = user
		session.Save(c.Request, c.Writer)
		c.HTML(http.StatusOK, "layout.html", gin.H{"username": user.Username})
		return
	}
	fmt.Println("err:", e)
	c.HTML(http.StatusUnauthorized, "authorization.html", gin.H{"message": "Incorrect username or password!"})
}

func profileHandler(c *gin.Context) {
	session, _ := store.Get(c.Request, "session")
	var user = &models.User{}
	val := session.Values["user"]
	var ok bool
	if user, ok = val.(*models.User); !ok {
		fmt.Println("was not of type *User")
		c.HTML(http.StatusForbidden, "authorization.html", nil)
		return
	}
	c.HTML(http.StatusOK, "layout.html", gin.H{"username": user.Username})

}

func logoutHandler(c *gin.Context) {
	session, _ := store.Get(c.Request, "session")
	delete(session.Values, "user")
	session.Save(c.Request, c.Writer)
	c.HTML(http.StatusOK, "authorization.html", gin.H{"message": "Logged out"})
}
