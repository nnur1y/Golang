package view

import (
	"Golang/config"
	"Golang/models"
	"fmt"

	"github.com/gin-gonic/gin"
)

// handleSearch,
func Search(c *gin.Context) {
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
