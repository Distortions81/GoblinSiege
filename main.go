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

	connectTwitch()

	err := twitchWriter.Say("xboxtv81", "testing.\n")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Sent 'testing.' to channel xboxtv81.")

	go func() {
		time.Sleep(time.Second * 2)
		UserMsgDict.Lock.Lock()
		startVote()
		UserMsgDict.Lock.Unlock()
	}()

	go dbAutoSave()

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
