package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// App is the application object. It contains the router and settings
type App struct {
	router *httprouter.Router
	port   int
}

// NewApp builds a new app object
func NewApp() (*App, error) {
	app := &App{
		router: httprouter.New(),
		port:   8080,
	}

	return app, nil
}

// Configure sets up the applications routes
func (a *App) Configure() {
	a.router.GET("/posts", GetPostsHandler)
	a.router.GET("/posts/:id", GetPostHandler)
	a.router.POST("/posts", CreatePostHandler)
	a.router.PUT("/posts/:id", UpdatePostHandler)
	a.router.PATCH("/posts/:id", PatchPostHandler)
	a.router.DELETE("/posts/:id", DeletePostHandler)
}

// Start fires up the web server
func (a *App) Start() {
	log.Printf("starting http server on port: %d", a.port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(a.port), a.router))
}
