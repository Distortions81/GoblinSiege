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

type playersListData struct {
	lock  sync.Mutex
	dirty bool
	idmap map[int64]*playerData
}

type playerData struct {
	Name   string
	Points int
}

func playersAutosave() {
	for ServerRunning {

		players.lock.Lock()
		if players.dirty {
			players.dirty = false
			writePlayers() //This unlocks after serialize
		} else {
			//No write to do, unlock
			players.lock.Unlock()
		}

		time.Sleep(time.Second * 30)
	}
}

// This unlocks playersLock after serialize
func writePlayers() {

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

	players.lock.Unlock()
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

	qlog("Reading players.")

	_, err := os.Stat(playersFile)
	notfound := os.IsNotExist(err)

	if !notfound { /* Otherwise just read in the config */
		file, err := os.ReadFile(playersFile)

		if file != nil && err == nil {

			//Only lock after file is read
			players.lock.Lock()
			defer players.lock.Unlock()

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
