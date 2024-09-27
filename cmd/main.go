package main

import (
	"log"

	"github.com/Vatsal-S-Patel/Bloggy/internal/api"
	"github.com/Vatsal-S-Patel/Bloggy/internal/app"
)

func main() {
	// make instance of app
	app, err := app.New()
	if err != nil {
		log.Println("ERROR failed to initialize app:", err.Error())
		return
	}

	// Initialize apis, register routes and start HTTP server
	api.ListenAndServe(app)
}
