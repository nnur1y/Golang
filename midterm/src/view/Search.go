package view

import (
	"Golang/config"
	"Golang/models"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

// handleSearch,
func Search(c *gin.Context, store *sessions.CookieStore) {
	query := c.Query("searchItem")
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

	c.HTML(200, "index.html", gin.H{
		"search":   false,
		"content":  RecipesList,
		"username": user.Username,
	})
}
