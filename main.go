package main

import (
	"log"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"github.com/templ-exemple/model"
	"github.com/templ-exemple/views"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB Config
var DB *gorm.DB

func Connect() {
	db, err := gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	DB = db
}

func migration() {
	DB.AutoMigrate(&model.Article{})
}

func Render(c *gin.Context, status int, template templ.Component) error {
	c.Status(status)
	return template.Render(c.Request.Context(), c.Writer)
}

func main() {
	r := gin.Default()

	Connect()
	migration()

	r.GET("/", home)
	r.GET("/new", newArticle)
	r.POST("/article/new", createArticle)

	r.Run("localhost:8080")
}

func home(c *gin.Context) {
	var articles []model.Article
	DB.Find(&articles)
	Render(c, 200, views.Home(articles))
}

func newArticle(c *gin.Context) {
	Render(c, 200, views.NewArticle())
}

func createArticle(c *gin.Context) {
	title := c.PostForm("title")
	content := c.PostForm("content")
	image, err := c.FormFile("image")
	if err != nil {
		log.Fatal("Failed to get image:", err)
	}
	imageName := image.Filename
	c.SaveUploadedFile(image, "uploads/images/"+imageName)

	DB.Create(&model.Article{Title: title, Content: content, Image: imageName})
	c.Redirect(302, "/")
}
