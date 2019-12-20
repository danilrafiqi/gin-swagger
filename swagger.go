package ginSwagger

import (
	"fmt"
	"html/template"
	"net/http"
	"path"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type Config struct {
	Route          *gin.RouterGroup
	User           map[string]string
	Url            string
	Docs           string
	Authentication bool
}

type swaggerUIBundle struct {
	URL         string
	DeepLinking bool
	SwaggerUrl  string
}

func authenticate(conf Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		fmt.Println("login status", session.Get("login"))
		if session.Get("login") == true {
			c.Next()
		} else {
			session.Clear()
			session.Save()
			c.Redirect(http.StatusMovedPermanently, conf.Url+"/login")
			c.Next()
		}
	}
}

func login(conf Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginPath = path.Join("views", "login.html")
		var login, _ = template.ParseFiles(loginPath)
		login.Execute(c.Writer, swaggerUIBundle{
			URL: conf.Url,
		})
	}
}

func loginProccess(conf Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		session := sessions.Default(c)

		for key, val := range conf.User {
			if key == username && val == password {
				session.Set("login", true)
				session.Save()
				c.Redirect(http.StatusMovedPermanently, conf.Url)
			} else {
				session.Clear()
				session.Save()
				c.Redirect(http.StatusMovedPermanently, conf.Url+"/login")
			}
		}
	}
}

func logout(conf Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()
		c.Redirect(http.StatusMovedPermanently, conf.Url+"/login")
	}
}

func index(conf Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var indexPath = path.Join("views", "index.html")
		var index, _ = template.ParseFiles(indexPath)
		index.Execute(c.Writer, swaggerUIBundle{
			URL: conf.Url + conf.Docs,
		})
	}
}

func Init(config Config) {
	config.Route.Static("/assets", "./assets")
	config.Route.StaticFile(config.Docs, "."+config.Docs)
	store := cookie.NewStore([]byte("secret"))
	config.Route.Use(sessions.Sessions("mysession", store))
	config.Route.GET("/login", login(config))
	config.Route.POST("/login", loginProccess(config))
	config.Route.GET("/logout", logout(config))
	if config.Authentication == true {
		config.Route.GET("/", authenticate(config), index(config))
	} else {
		config.Route.GET("/", index(config))
	}
}
