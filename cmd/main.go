package main

import (
	"log"
	"sachahjkl/htmx_go/pkg/config"
	"sachahjkl/htmx_go/pkg/db"
	"sachahjkl/htmx_go/pkg/handler"
	"sachahjkl/htmx_go/pkg/router"

	"github.com/joho/godotenv"
)

func main() {

	// Try to load a potentially existing ".env" file
	// ignore errors
	godotenv.Load()

	c, err := config.LoadConfig()

	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	log.Printf("Version : %v\n", c.Version)

	db := db.Init(c.DBUrl)

	app := router.New(c)

	router.RegisterDefaultRoutes(app, db)
	handler.RegisterTodoRoutes(app, db, c)
	handler.RegisterAuthRoutes(app, db, c)

	log.Fatal(app.Listen(":" + c.Port))

}
