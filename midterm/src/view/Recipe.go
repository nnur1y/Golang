package view

import (
	"Golang/config"
	"Golang/models"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

//  handlerRecipe,
// handlerIndex,

func Recipe(c *gin.Context) {
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
	commentsResult, err := db.Query("SELECT u.Username, c.comment FROM feedback c join users u on u.id =c.userid WHERE c.recipeid = ?", recipeid)
	if err != nil {
		fmt.Print(err)
	}
	comment := models.CommentData{}
	commentsList := []models.CommentData{}

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

func MainPage(c *gin.Context, store *sessions.CookieStore) {
	var user models.User
	session, err := store.Get(c.Request, "session")
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println("session:", session)
	userVal, ok := session.Values["user"]

	if !ok {
		fmt.Println("no user", userVal)
		user.Username = ""
	} else {
		user = *userVal.(*models.User)
	}

	db, _ := config.LoadDB()
	var searchItem models.SearchItem

	result, err := db.Query("select rc.Id_r, rc.name, rc.description ,rc.categories ,rt.rating from recipe rc join ratings rt on rc.ID_r = rt.recipeId ")
	if err != nil {
		fmt.Print(err)
	}
	RecipesList, e := searchItem.Search(result)

	if e != nil {
		fmt.Println(e)
	}

	c.HTML(200, "layout.html", gin.H{
		"search":   false,
		"content":  RecipesList,
		"username": user.Username,
	})

}
