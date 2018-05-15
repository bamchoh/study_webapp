package main

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/bamchoh/study_webapp/dao"
	"github.com/bamchoh/study_webapp/models"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	_ "github.com/lib/pq"
)

var (
	repeat int
	db     *sql.DB
)

func repeatHandler(c *gin.Context) {
	var buffer bytes.Buffer
	for i := 0; i < repeat; i++ {
		buffer.WriteString("Hello from Go!!\n")
	}
	c.String(http.StatusOK, buffer.String())
}

func dbFunc(c *gin.Context) {
	if _, err := db.Exec("INSERT INTO ticks VALUES (now())"); err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error incrementing tick: %q", err))
		return
	}

	rows, err := db.Query("SELECT tick FROM ticks")
	if err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error reading ticks: %q", err))
		return
	}

	defer rows.Close()
	var tickList []string
	for rows.Next() {
		var tick time.Time
		if err := rows.Scan(&tick); err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error scanning ticks: %q", err))
			return
		}
		tickList = append(tickList, tick.String())
	}
	c.HTML(http.StatusOK, "db.tmpl.html", gin.H{
		"title": "hogehoge",
		"list":  tickList,
	})
}

func singinFunc(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	builder := &dao.UserDao{DB: db}

	user, err := builder.Get(email, password)
	if err != nil {
		log.Println("Get", err)
		err = errors.New("Signin was failed")
		c.HTML(http.StatusInternalServerError,
			"index.tmpl.html",
			fmt.Sprintf("Error Signin: %q", err))
		return
	}

	c.Redirect(http.StatusMovedPermanently, "/")
}

func signoutFunc(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/")
}

func signupFunc(c *gin.Context) {
	id := c.PostForm("id")
	name := c.PostForm("name")
	email := c.PostForm("email")
	password := c.PostForm("password")
	confirmpassword := c.PostForm("confirmpassword")

	builder := &dao.UserDao{DB: db}

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

type server struct {
	router gin.IRouter
}

func main() {
	var err error
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	tStr := os.Getenv("REPEAT")
	repeat, err = strconv.Atoi(tStr)
	if err != nil {
		log.Printf("Error converting $REPEAT to an int: %q - Using default\n", err)
		repeat = 5
	}

	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}

	engine := gin.New()
	store := sessions.NewCookieStore([]byte("secret"))
	engine.Use(sessions.Sessions("SessionName", store))
	engine.Use(gin.Logger())
	engine.LoadHTMLGlob("templates/*.tmpl.html")
	engine.Static("/assets", "./assets")
	engine.Static("/static", "static")

	router := server{engine}

	engine.Run(":" + port)
}
