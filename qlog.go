package main

import (
	"fmt"
	"log"
)

func qlog(format string, args ...interface{}) {
	buf := fmt.Sprintf(format, args...)
	log.Print(buf)
}
