package controller

import (
	"Golang/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

func Auth(c *gin.Context, store *sessions.CookieStore) {
	fmt.Println("auth middleware running")
	session, err := store.Get(c.Request, "session")
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println("session:", session)
	val, ok := session.Values["user"]
	if !ok {
		c.HTML(http.StatusForbidden, "authorization.html", nil)
		c.Abort()
		return
	}
	user, ok := val.(*models.User)
	if !ok {
		fmt.Println("was not of type *User")
		c.HTML(http.StatusForbidden, "authorization.html", nil)
		return
	}
	c.Set("user", user)
	fmt.Println("middleware done")
	c.Next()
}

func Logout(c *gin.Context, store *sessions.CookieStore) {
	session, err := store.Get(c.Request, "session")
	if err != nil {
		fmt.Print(err)
	}
	delete(session.Values, "user")
	session.Save(c.Request, c.Writer)
	c.HTML(http.StatusOK, "authorization.html", gin.H{"message": "Logged out"})
}

func Profile(c *gin.Context, store *sessions.CookieStore) {
	session, _ := store.Get(c.Request, "session")
	var user = &models.User{}
	val := session.Values["user"]
	var ok bool
	if user, ok = val.(*models.User); !ok {
		fmt.Println("was not of type *User")
		c.HTML(http.StatusForbidden, "authorization.html", nil)
		return
	}
	fmt.Println(user)
	// c.HTML(http.StatusOK, "profile.html", gin.H{"user": user})

}
