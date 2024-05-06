package main

import (
	"fmt"
	"os"

	"profiln-be/config"
	"profiln-be/delivery/http/route"

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
	route.NewRoute(app, db, log)
	port := os.Getenv("PORT")

	log.Fatal(app.Run(fmt.Sprintf(":%s", port)))
}
