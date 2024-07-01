package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func info(message string) {
	if os.Getenv("LOG_LEVEL") == "DEBUG" || os.Getenv("LOG_LEVEL") == "INFO" {
		log.Printf("Hestia %s: %s", formatDate(time.Now()), message)
	}
}

func error(message string) {
	log.Printf("Hestia %s", message)
}

func padTo2Digits(num int) string {
	return fmt.Sprintf("%02d", num)
}

func formatDate(date time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d.%03d",
		date.Year(), date.Month(), date.Day(),
		date.Hour(), date.Minute(), date.Second(), date.Nanosecond()/1e6)
}
