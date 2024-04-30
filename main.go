package main

import (
	"fmt"
	"log"
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

	app := gin.Default()
	port := os.Getenv("PORT")
	route.NewRoute(app, db)

	log.Fatal(app.Run(fmt.Sprintf(":%s", port)))
}
