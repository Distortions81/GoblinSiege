package main

import "time"

func main() {
	readDB()

	dbLock.Lock()
	WriteDB() //Unlocks after serialize

	connectTwitch()
	go dbAutoSave()
	startEbiten()
}

func dbAutoSave() {
	for {

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
