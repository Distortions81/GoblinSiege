package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
)

var authInfo authData

type authData struct {
	Username  string
	AuthToken string
}

func readAuth() {

	_, err := os.Stat(authFile)
	notfound := os.IsNotExist(err)

	if !notfound {
		file, err := os.ReadFile(authFile)

		if file != nil && err == nil {
			log.Println("Reading auth info.")

			err := json.Unmarshal([]byte(file), &authInfo)
			if err != nil {
				log.Fatal("readAuth: Unmarshal failure")
				return
			}

			if authInfo.AuthToken == "" || authInfo.Username == "" {
				log.Fatal("Missing username or token in auth file.")
				return
			}
			return
		}
	}

	log.Println("No auth file found, attempting to create.")
	writeAuth()
	log.Fatal("Please add your username and token to the auth file.")
}

func writeAuth() {

	outbuf := new(bytes.Buffer)
	enc := json.NewEncoder(outbuf)

	if err := enc.Encode(authInfo); err != nil {
		log.Fatal("WriteGCfg: enc.Encode failure")
		return
	}

	_, err := os.Create(authFile)

	if err != nil {
		log.Fatal("WriteGCfg: os.Create failure")
		return
	}

	err = os.WriteFile(authFile, outbuf.Bytes(), 0644)

	if err != nil {
		log.Fatal("WriteGCfg: WriteFile failure")
		return
	}

}
