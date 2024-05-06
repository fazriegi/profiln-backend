package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

func NewLogger() (*logrus.Logger, *os.File) {
	outputPath := os.Getenv("LOG_OUTPUT_PATH")
	filePath := fmt.Sprintf("%s/profiln.log", outputPath)

	// check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// if the file doesn't exist, create it
		if _, err := os.Create(filePath); err != nil {
			log.Fatal("failed to create file:", err)
		}
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal("failed to open file:", err)
	}

	log := logrus.New()
	logLevel, _ := strconv.Atoi(os.Getenv("LOG_LEVEL"))

	log.SetLevel(logrus.Level(logLevel))
	log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
	log.SetOutput(file)

	return log, file
}
