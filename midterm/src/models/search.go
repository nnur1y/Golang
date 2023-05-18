package models

import (
	"database/sql"
)

type SearchItem struct {
	SearchItem string `json:"searchItem"`
}

func (s SearchItem) Search(result *sql.Rows) ([]Recipe, error) {

	var err error
	recipes := []Recipe{}

	for result.Next() {
		var r Recipe
		// rc.Id_r, rc.name, rc.description ,rc.categories , rt.image_r, ,rt.rating
		err = result.Scan(&r.Id_r, &r.Name, &r.Description, &r.Categories, &r.ImgURL, &r.Rate)

		recipes = append(recipes, r)

		if err != nil {
			panic(err)
		}

	}
	return recipes, err
}
