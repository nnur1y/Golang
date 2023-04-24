package models

import (
	"database/sql"
)

type SearchItem struct {
	SearchItem string `json:"searchItem"`
}

func (s SearchItem) Search(result *sql.Rows) ([]Recipe, error) {

	var err error
	recipe := Recipe{}
	recipes := []Recipe{}

	for result.Next() {
		var id int
		var name string
		var description string
		var Categories string
		var rating string

		err = result.Scan(&id, &name, &description, &Categories, &rating)

		recipe.Id_r = id
		recipe.Name = name
		recipe.Description = description
		recipe.Categories = Categories
		recipe.Rate = rating

		recipes = append(recipes, recipe)

		if err != nil {
			panic(err)
		}

	}
	return recipes, err
}
