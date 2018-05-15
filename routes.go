package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday"
)

func (s *server) routes() {
	s.router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	s.router.GET("/mark", func(c *gin.Context) {
		c.String(http.StatusOK, string(blackfriday.MarkdownBasic([]byte("**hi!**"))))
	})

	s.router.GET("/signup", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signup.tmpl.html", nil)
	})
	s.router.POST("/signup", signupFunc)

	s.router.POST("/signin", signinFunc)

	s.router.POST("/signout", signoutFunc)

	s.router.GET("/repeat", repeatHandler)
	s.router.GET("/db", dbFunc)
}
