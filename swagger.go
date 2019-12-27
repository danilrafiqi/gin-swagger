package ginSwagger

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag"
)

type Config struct {
	Route          *gin.RouterGroup
	User           map[string]string
	Url            string
	Urls           []Urls
	Authentication bool
}

type Urls struct {
	Url  string `json:"url"`
	Name string `json:"name"`
}

type swaggerUIBundle struct {
	URL            string
	URLS           interface{}
	Authentication bool
	DeepLinking    bool
	SwaggerUrl     string
}

func authenticate(conf Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var session = sessions.Default(c)
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
		var loginPath = path.Join("swagger/views", "login.html")
		var login, _ = template.ParseFiles(loginPath)
		login.Execute(c.Writer, swaggerUIBundle{
			URL: conf.Url,
		})
	}
}

func loginProccess(conf Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var username = c.PostForm("username")
		var password = c.PostForm("password")
		var session = sessions.Default(c)

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
		var session = sessions.Default(c)
		session.Clear()
		session.Save()
		c.Redirect(http.StatusMovedPermanently, conf.Url+"/login")
	}
}

func index(conf Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var indexPath = path.Join("swagger/views", "index.html")
		var index, _ = template.ParseFiles(indexPath)
		var jsonData, _ = json.Marshal(conf.Urls)
		if len(conf.Urls) > 0 {
			index.Execute(c.Writer, swaggerUIBundle{
				URLS:           string(jsonData),
				Authentication: conf.Authentication,
			})
		} else {
			doc, _ := swag.ReadDoc()
			var url string
			if doc == "" {
				url = conf.Url + "/docs/swagger.json"
			} else {
				url = conf.Url + "/docs.json"
			}

			index.Execute(c.Writer, swaggerUIBundle{
				URL:            url,
				Authentication: conf.Authentication,
			})
		}
	}
}

func docToJson(conf Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		doc, _ := swag.ReadDoc()
		c.Writer.Write([]byte(doc))
	}
}

func Init(config Config) {
	var store = cookie.NewStore([]byte("secret"))
	config.Route.Use(sessions.Sessions("mysession", store))
	config.Route.Static("/assets", "./swagger/assets")
	config.Route.Static("/docs", "./swagger/docs/")
	config.Route.GET("/docs.json", docToJson(config))
	config.Route.GET("/login", login(config))
	config.Route.POST("/login", loginProccess(config))
	config.Route.GET("/logout", logout(config))
	if config.Authentication == true {
		config.Route.GET("/", authenticate(config), index(config))
	} else {
		config.Route.GET("/", index(config))
	}
}
