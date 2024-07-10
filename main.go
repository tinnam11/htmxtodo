package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

func main() {
	//creates context that listens to signals for program to be terminated
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	r := gin.Default()
	r.LoadHTMLGlob("./*.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/todos", getTodo)
	r.DELETE("/todos", deleteTodo)
	r.POST("/todos", addTodo)
	// c.Redirect(http.StatusFound, "/")

	srv := http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: r,
	}

	go func() {
		//request gets cancelled
		<-ctx.Done()
		fmt.Println("shuttign down...")
		//emits a timeout signal through the context
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Println(err)
			}
		}
	}()

	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
	}
}

// function appends other items into a new list that does not need to be removed
func removeItem(id string, todos []Todo) []Todo {
	newList := []Todo{}
	for _, v := range todos {
		if v.ID != id {
			newList = append(newList, v)
		}
	}
	return newList
}

func deleteTodo(ctx *gin.Context) {
	var todo Todo
	if err := ctx.BindJSON(&todo); err != nil {
		return
	}
	todos := readJsonFile()
	//rewrites items that should not be deleted
	newTodo := removeItem(todo.ID, todos)
	file, _ := json.MarshalIndent(newTodo, "", " ")
	_ = os.WriteFile("todos.json", file, 0644)
	getTodo(ctx)
}

func getTodo(ctx *gin.Context) {
	todos := readJsonFile()
	ctx.HTML(http.StatusOK, "todos.html", todos)
}

func addTodo(ctx *gin.Context) {
	var newTodo Todo
	if err := ctx.BindJSON(&newTodo); err != nil {
		return
	}
	writeJsonFile(newTodo)
	getTodo(ctx)
}

func readJsonFile() []Todo {
	//opends a file and creates a new one if it does not exist
	file, err := os.OpenFile("todos.json", os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	byteValue, _ := io.ReadAll(file)
	var todos []Todo
	json.Unmarshal(byteValue, &todos)
	return todos
}

func writeJsonFile(todo Todo) {
	todo.ID = uuid.New().String()
	todos := readJsonFile()
	todos = append(todos, todo)
	file, _ := json.MarshalIndent(todos, "", " ")
	_ = os.WriteFile("todos.json", file, 0644)

}

type Todo struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}
