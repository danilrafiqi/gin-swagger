package ginSwagger

import (
	"html/template"
	"net/http"
	"path"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type swaggerUIBundle struct {
	URL         string
	DeepLinking bool
	SwaggerUrl  string
}

func checkAuth(swaggerUrl string) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if session.Get("login") == true {
			c.Next()
		} else {
			session.Clear()
			c.Redirect(http.StatusMovedPermanently, swaggerUrl+"/login")
			c.Next()
		}
	}
}

func InitAuth(r *gin.RouterGroup, user map[string]string, swaggerUrl string) {
	var loginPath = path.Join("views", "login.html")
	var login, _ = template.ParseFiles(loginPath)
	r.Static("/assets", "./assets")
	r.StaticFile("/swagger.json", "./docs/swagger.json")
	r.GET("/login", func(c *gin.Context) {
		login.Execute(c.Writer, swaggerUIBundle{
			URL: swaggerUrl,
		})
	})

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		session := sessions.Default(c)

		for key, val := range user {
			if key == username && val == password {
				session.Set("login", true)
				session.Save()
				c.Redirect(http.StatusMovedPermanently, swaggerUrl)
			} else {
				session.Clear()
				c.Redirect(http.StatusMovedPermanently, swaggerUrl+"/login")
			}
		}

	})
	var indexPath = path.Join("views", "index.html")
	var index, _ = template.ParseFiles(indexPath)
	r.GET("/", checkAuth(swaggerUrl), func(c *gin.Context) {
		index.Execute(c.Writer, swaggerUIBundle{
			URL: swaggerUrl + "/swagger.json",
		})
	})
}
