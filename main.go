package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var ServerRunning bool = true
var ServerIsStopped bool
var roundTime time.Duration = time.Second * 10

func main() {
	go startEbiten()

	readDB()

	dbLock.Lock()
	WriteDB() //Unlocks after serialize

	UserMsgDict.Users = make(map[int64]*userMsgData)
	gameMap = make(map[xyi]*objectData)

	go dbAutoSave()
	go connectTwitch()

	time.Sleep(time.Second * 5)
	startVote()

	//After starting loops, wait here for process signals
	signalHandle := make(chan os.Signal, 1)

	signal.Notify(signalHandle, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-signalHandle

	ServerRunning = false

	log.Println("Saving DB...")
	dbLock.Lock()
	WriteDB()
}

func dbAutoSave() {
	for ServerRunning {

		dbLock.Lock()
		if dbDirty {
			dbDirty = false
			WriteDB() //This unlocks after serialize
		} else {
			//No write to do, unlock
			dbLock.Unlock()
		}

		time.Sleep(time.Second * 30)
	}
}
