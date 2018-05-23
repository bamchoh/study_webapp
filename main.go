package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bamchoh/study_webapp/dao"
	"github.com/bamchoh/study_webapp/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	_ "github.com/lib/pq"
)

type SessionInfo struct {
	ID             string
	Name           string
	IsSessionAlive bool
}

type server struct {
	db     *sql.DB
	router gin.IRouter
}

func (s *server) appIndex(c *gin.Context) {
	session := sessions.Default(c)
	alive := session.Get("alive")
	user_id := session.Get("user_id")
	c.HTML(http.StatusOK, "index.tmpl.html", gin.H{"alive": alive, "user_id": user_id})
}

func (s *server) signupGet(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.tmpl.html", nil)
}

func (s *server) signupPost(c *gin.Context) {
	id := c.PostForm("id")
	name := c.PostForm("name")
	email := c.PostForm("email")
	password := c.PostForm("password")
	confirmpassword := c.PostForm("confirmpassword")

	builder := &dao.UserDao{DB: s.db}

	user := models.User{
		ID:    id,
		Name:  name,
		Email: email,
	}

	if err := builder.Create(user, password, confirmpassword); err != nil {
		log.Println("Create", err)
		err = errors.New("User creation was failed")
		c.HTML(http.StatusInternalServerError,
			"signup.tmpl.html",
			fmt.Sprintf("Error Signup: %q", err))
		return
	}

	c.Redirect(http.StatusMovedPermanently, "/")
}

func (s *server) signinPost(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	builder := &dao.UserDao{DB: s.db}

	user, err := builder.Get(email, password)
	if err != nil {
		log.Println("Get", err)
		err = errors.New("Signin was failed")
		c.HTML(http.StatusInternalServerError,
			"index.tmpl.html",
			fmt.Sprintf("Error Signin: %q", err))
		return
	}

	session := sessions.Default(c)
	session.Set("alive", true)
	session.Set("user_id", user.ID)
	session.Set("user_name", user.Name)
	session.Save()

	c.Redirect(http.StatusMovedPermanently, "/")
}

func (s *server) signoutPost(c *gin.Context) {
	fmt.Println("signoutPost")
	session := sessions.Default(c)
	session.Set("alive", nil)
	session.Set("user_id", nil)
	session.Set("user_name", nil)
	session.Save()
	alive := session.Get("alive")
	fmt.Println("  alive", alive)
	c.Redirect(http.StatusMovedPermanently, "/")
}

func (s *server) settingsGet(c *gin.Context) {
	session := sessions.Default(c)
	alive := session.Get("alive")
	if alive != nil {
		id := session.Get("user_id")
		if id != nil {
			c.String(http.StatusOK, fmt.Sprintf("test: %s", id))
			return
		}
	}
	c.Redirect(http.StatusOK, "/")
}

func main() {
	var err error
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	proto := os.Getenv("PROTO")
	if proto == "" {
		proto = "https"
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}

	engine := gin.New()
	store := cookie.NewStore([]byte("secret"))
	engine.Use(sessions.Sessions("SessionName", store))
	engine.Use(gin.Logger())
	engine.Use(func(c *gin.Context) {
		if proto != "https" {
			return
		}

		header := c.Request.Header
		isssl := false
		if params, ok := header["X-Forwarded-Proto"]; ok {
			if len(params) != 0 && params[0] == "https" {
				isssl = true
			}
		}
		if !isssl {
			req := c.Request
			c.Redirect(http.StatusMovedPermanently, "https://"+req.Host+req.URL.Path)
			// c.String(http.StatusInternalServerError, "http does not support. please access over https")
			c.Abort()
		}
	})
	engine.LoadHTMLGlob("templates/*.tmpl.html")
	engine.Static("/assets", "./assets")
	engine.Static("/static", "static")

	svr := server{db, engine}

	svr.routes()

	engine.Run(":" + port)
}
