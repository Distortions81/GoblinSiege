package main

import (
	"encoding/json"
	"log"
	"os"
)

var userSettings settingsData

type settingsData struct {
	Username  string
	AuthToken string
}

func readSettings() {

	_, err := os.Stat(authFile)
	notfound := os.IsNotExist(err)

	if !notfound {
		file, err := os.ReadFile(authFile)

		if file != nil && err == nil {
			log.Println("Reading settings.")

			err := json.Unmarshal([]byte(file), &userSettings)
			if err != nil {
				log.Fatal("readAuth: Unmarshal failure")
				return
			}

			if userSettings.AuthToken == "" || userSettings.Username == "" {
				log.Fatal("readSettings: Missing username or token in settings.")
				return
			}

			writeSettings()
			return
		}
	}

	log.Println("No settings found, attempting to create.")
	writeSettings()
	log.Fatal("Please add your username and token to the settings file.")
}

func writeSettings() {

	var err error
	outbuf, err := json.MarshalIndent(userSettings, "", "    ")

	if err != nil {
		log.Fatal("WriteGCfg: json marshal error.")
		return
	}

	_, err = os.Create(authFile)

	if err != nil {
		log.Fatal("WriteGCfg: os.Create failure")
		return
	}

	err = os.WriteFile(authFile, outbuf, 0644)

	if err != nil {
		log.Fatal("WriteGCfg: WriteFile failure")
		return
	}

}