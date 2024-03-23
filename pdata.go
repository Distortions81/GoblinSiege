package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	"github.com/hako/durafmt"
)

var players playersListData
var pLock sync.Mutex

type playersListData struct {
	dirty bool
	idmap map[int64]*playerData
}

type playerData struct {
	Name   string
	Points int
}

func playersAutosave() {
	for ServerRunning {

		pLock.Lock()
		if players.dirty {
			players.dirty = false
			writePlayers()
		}
		pLock.Unlock()
		time.Sleep(time.Second * 30)
	}
}

func writePlayers() {

	pLock.Lock()
	defer pLock.Unlock()

	qlog("Saving players...")
	startTime := time.Now()

	tempPath := playersFile + ".tmp"
	finalPath := playersFile

	outbuf := new(bytes.Buffer)
	enc := json.NewEncoder(outbuf)

	if err := enc.Encode(players.idmap); err != nil {
		log.Fatal("writePlayers: enc.Encode failure")
		return
	}

	qlog("serialize players took: %v", durafmt.Parse(time.Since(startTime)).LimitFirstN(2))

	_, err := os.Create(tempPath)

	if err != nil {
		log.Fatal("writePlayers: os.Create failure")
		return
	}

	err = os.WriteFile(tempPath, outbuf.Bytes(), 0644)

	if err != nil {
		log.Fatal("writePlayers: WriteFile failure")
		return
	}

	err = os.Rename(tempPath, finalPath)

	if err != nil {
		log.Fatal("Couldn't rename players file.")
		return
	}

	qlog("Write player file took: %v", durafmt.Parse(time.Since(startTime)).LimitFirstN(2))
}

// Load player scores
func readPlayers() {

	pLock.Lock()
	defer pLock.Unlock()

	qlog("Reading players.")

	_, err := os.Stat(playersFile)
	notfound := os.IsNotExist(err)

	if !notfound { /* Otherwise just read in the config */
		file, err := os.ReadFile(playersFile)

		if file != nil && err == nil {

			err := json.Unmarshal([]byte(file), &players.idmap)
			if len(players.idmap) == 0 {
				//Empty database, create a map
				players.idmap = make(map[int64]*playerData)
			}
			if err != nil {
				log.Fatal("readPlayers: Unmarshal failure")
			}
		} else {
			qlog("No players file.")
		}
	}
}
