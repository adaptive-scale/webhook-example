package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {

	sharedSecret := os.Getenv("SHARED_SECRET")
	outputType := os.Getenv("OUTPUT_TYPE")
	fileLocation := os.Getenv("FILE_LOCATION")
	maxSize := os.Getenv("MAX_SIZE")
	maxBackup := os.Getenv("MAX_BACKUP")
	maxAge := os.Getenv("MAX_AGE")
	port := os.Getenv("PORT")
	formatter := os.Getenv("FORMATTER")

	log := logrus.New()

	if outputType == "" {
		outputType = "stdout"
	}

	if outputType == "file" {

		if fileLocation == "" {
			fileLocation = "/tmp/adaptive.log"
		}

		if maxSize == "" {
			maxSize = "10"
		}

		if maxBackup == "" {
			maxBackup = "3"
		}

		if maxAge == "" {
			maxAge = "28"
		}

		maxSizeI, err := strconv.Atoi(maxSize)
		if err != nil {
			log.Fatal(err)
		}

		maxBackupI, err := strconv.Atoi(maxBackup)
		if err != nil {
			log.Fatal(err)
		}

		maxAgeI, err := strconv.Atoi(maxAge)
		if err != nil {
			log.Fatal(err)
		}

		// Configure log rotation
		log.SetOutput(&lumberjack.Logger{
			Filename:   fileLocation,
			MaxSize:    maxSizeI, // Max megabytes before rotating
			MaxBackups: maxBackupI,  // Max number of old log files to retain
			MaxAge:     maxAgeI, // Max days to retain old logs
			Compress:   true, // Compress rotated files
		})

		if formatter == "json" {
			log.SetFormatter(&logrus.JSONFormatter{})
		}
	}

	http.HandleFunc("/healthz", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "ok")
	})

	http.HandleFunc("/api/hook", func(writer http.ResponseWriter, request *http.Request) {

		if request.Method != http.MethodPost {
			http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		incomingSharedSecret := request.Header.Get("Authorization")

		if incomingSharedSecret != sharedSecret {
			http.Error(writer, "Unauthorized", http.StatusUnauthorized)
			return
		}

		d, _ := io.ReadAll(request.Body)
		log.Println(string(d))

		fmt.Fprintf(writer, "ok")
	})

	if port == "" {
		port = "8080"
	}

	fmt.Println("Starting webhook server with mode: " + outputType)
	fmt.Println("Starting webhook server on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
