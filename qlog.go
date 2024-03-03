package main

import (
	"fmt"
	"log"
)

func qlog(format string, args ...interface{}) {
	buf := fmt.Sprintf(format, args...)
	log.Print(buf)
}

func sayLog(format string, args ...interface{}) {
	buf := fmt.Sprintf(format, args...)
	log.Print(buf)
	client.Say(userSettings.Channel, buf)
}
