package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"sync"
)

var dbLock sync.Mutex
var dbDirty bool

var Players map[int64]*playerData

type playerData struct {
	Points int
}

// This unlocks dbLock after serialize
func WriteDB() {

	tempPath := dbFile + ".tmp"
	finalPath := dbFile

	outbuf := new(bytes.Buffer)
	enc := json.NewEncoder(outbuf)

	if err := enc.Encode(Players); err != nil {
		log.Fatal("WriteGCfg: enc.Encode failure")
		return
	}

	dbLock.Unlock()

	_, err := os.Create(tempPath)

	if err != nil {
		log.Fatal("WriteGCfg: os.Create failure")
		return
	}

	err = os.WriteFile(tempPath, outbuf.Bytes(), 0644)

	if err != nil {
		log.Fatal("WriteGCfg: WriteFile failure")
		return
	}

	err = os.Rename(tempPath, finalPath)

	if err != nil {
		log.Fatal("Couldn't rename Gcfg file.")
		return
	}

	log.Println("Wrote db.")
}

/* Read in cached list of Discord players with specific roles */
func readDB() {
	dbLock.Lock()
	defer dbLock.Unlock()

	_, err := os.Stat(dbFile)
	notfound := os.IsNotExist(err)

	if !notfound { /* Otherwise just read in the config */
		file, err := os.ReadFile(dbFile)

		if file != nil && err == nil {

			log.Println("Reading db.")

			err := json.Unmarshal([]byte(file), &Players)
			if len(Players) == 0 {
				//Empty database, create a map
				Players = make(map[int64]*playerData)
			}
			if err != nil {
				log.Fatal("Readcfg.RoleList: Unmarshal failure")
			}
		} else {
			log.Println("No database file.")
		}
	}
}
