package main

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

var userSettings settingsData
var settingsLock sync.Mutex

type settingsData struct {
	UserName  string
	AuthToken string
	CmdPrefix string
}

func readSettings() {

	settingsLock.Lock()
	defer settingsLock.Unlock()

	_, err := os.Stat(authFile)
	notfound := os.IsNotExist(err)

	if !notfound {
		file, err := os.ReadFile(authFile)

		if file != nil && err == nil {
			qlog("Reading settings.")

			err := json.Unmarshal([]byte(file), &userSettings)
			if err != nil {
				log.Fatal("readAuth: Unmarshal failure")
				return
			}

			if userSettings.AuthToken == "" || userSettings.UserName == "" {
				log.Printf("readSettings: Missing UserName, BotName or AuthToken in settings.")
				return
			}

			if userSettings.CmdPrefix == "" {
				userSettings.CmdPrefix = "!"
			}

			writeSettings()
			return
		}
	}

	qlog("No settings found, attempting to create.")
	writeSettings()
	log.Fatal("Please add your UserName, BotName and AuthToken to the settings file.")
}

func writeSettings() {

	settingsLock.Lock()
	defer settingsLock.Unlock()

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
