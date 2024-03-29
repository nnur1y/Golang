package view

import (
	"Golang/config"
	"Golang/models"
	"fmt"

	"github.com/gin-gonic/gin"
)

// handleProductBreakfast functions.
func GetCategory(c *gin.Context, categoryName string) []models.Recipe {
	var searchItem models.SearchItem
	db, _ := config.LoadDB()
	result, err := db.Query("select rc.Id_r, rc.name, rc.description ,rc.categories,rc.image_r ,rt.rating from recipe rc join ratings rt on rc.ID_r = rt.recipeId WHERE  rc.categories= ? order by rt.rating desc ", categoryName)
	if err != nil {
		fmt.Print(err)
	}
	RecipesList, e := searchItem.Search(result)
	if e != nil {
		fmt.Print(e)
	}

	return RecipesList
}
