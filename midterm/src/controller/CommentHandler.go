package controller

import (
	"Golang/models"
	"Golang/view"
	"fmt"

	"github.com/gin-gonic/gin"
)

func SendComment(c *gin.Context) {
	var sendCom models.SendCommentData
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

		view.Recipe(c)
		return
	}

}
