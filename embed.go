package main

import (
	"embed"
	"log"
)

var (
	//go:embed data
	f embed.FS
)

func getFont(name string) []byte {
	data, err := f.ReadFile("data/fonts/" + name)
	if err != nil {
		log.Fatal(err)
	}
	return data

}
