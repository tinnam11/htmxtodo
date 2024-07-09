package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("./*.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/todos", func(c *gin.Context) {
		c.HTML(http.StatusOK, "todos.html", todos)
	})
	// c.Redirect(http.StatusFound, "/")

	log.Println("Starting server on :8080")
	r.Run()
}

type Todo struct {
	ID    int
	Title string
	Done  bool
}

var todos = []Todo{
	{1, "Learn Go", false},
	{2, "Build a Todo App", false},
}
