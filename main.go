package main

import (
	"fmt"
	"os"

	"profiln-be/config"
	"profiln-be/delivery/http/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	db := config.NewDatabase()
	defer db.Close()

	log, file := config.NewLogger()
	defer file.Close()

	app := gin.New()
	app.RedirectTrailingSlash = false
	app.Use(gin.Logger())

	routes.NewRoute(app, db, log)
	port := os.Getenv("PORT")

	certFile := os.Getenv("TLS_CERT_FILE")
	keyFile := os.Getenv("TLS_KEY_FILE")

	if certFile == "" || keyFile == "" {
		log.Fatal("TLS_CERT_FILE and TLS_KEY_FILE must be set")
	}

	log.Fatal(app.RunTLS(fmt.Sprintf(":%s", port), certFile, keyFile))
}
