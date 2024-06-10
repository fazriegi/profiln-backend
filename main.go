package main

import (
	"fmt"
	"os"

	"profiln-be/config"
	"profiln-be/delivery/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	db := config.NewDatabase()
	defer db.Close()

	log, file := config.NewLogger()
	defer file.Close()

	app := gin.Default()
	routes.NewRoute(app, db, log)
	port := os.Getenv("PORT")

	log.Fatal(app.Run(fmt.Sprintf(":%s", port)))
}
