package ginSwagger

import (
	"html/template"
	"path"

	"github.com/gin-gonic/gin"
)

type swaggerUIBundle struct {
	URL         string
	DeepLinking bool
	SwaggerUrl  string
}

func InitAuth(r *gin.RouterGroup, user map[string]string) {
	var indexPath = path.Join("views", "index.html")
	var index, _ = template.ParseFiles(indexPath)
	r.GET("/", func(c *gin.Context) {
		index.Execute(c.Writer, &swaggerUIBundle{
			URL:         "http://localhost:8200/docs.json",
			DeepLinking: true,
		})
	})
}
