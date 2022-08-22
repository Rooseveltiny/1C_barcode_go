package main

import (
	"log"
	"os"
)

func writeDownLogMessage(msg string, err error) {

	var errorMsg string
	if err != nil {
		errorMsg = err.Error()
	}

	file, errFile := os.OpenFile("logging.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if errFile != nil {
		log.Fatal(err)
	}
	defer file.Close()
	logger := log.New(file, "App logger!", log.Ltime)
	logger.Fatalf("%s   %s", msg, errorMsg)
}
