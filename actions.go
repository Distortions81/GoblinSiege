package main

import (
	"log"
	"math/rand"
	"time"
)

func addTower() {
	tower1 := getOtype("Stone Tower")

	if votes.VoteCount > 0 &&
		votes.Result.X > 0 &&
		votes.Result.Y > 0 &&
		votes.Result.X <= boardSizeX &&
		votes.Result.Y <= boardSizeY {

		tpos := votes.Result
		if board.enemyMap[tpos] == nil && board.playMap[tpos] == nil {
			board.playMap[tpos] = &objectData{Pos: tpos, oTypeP: tower1, Health: tower1.maxHealth}
		} else {
			log.Println("COLLISION!")
		}
	} else {

		log.Println("Not enough votes, picking random.")
		//Invalid or not enough votes, pick a pos at random
		tpos := xyi{X: rand.Intn(boardSizeX-1) + 1, Y: rand.Intn(boardSizeY-1) + 1}
		if board.enemyMap[tpos] == nil && board.playMap[tpos] == nil {
			board.playMap[tpos] = &objectData{Pos: tpos, oTypeP: tower1, Health: tower1.maxHealth}
		}
	}

}

func towerShootArrow() {
	curTime := time.Now()

	for _, item := range board.playMap {
		if item.dead {
			continue
		}
		for _, enemy := range board.enemyMap {
			if enemy.dead {
				continue
			}
			//If enemy within range
			if Distance(item.Pos, enemy.Pos) < 6 {

				if rand.Intn(2) != 0 {
					arrow := arrowData{tower: item.Pos, target: enemy.Pos, missed: true, shot: curTime}
					board.arrowsShot = append(board.arrowsShot, arrow)
					break
				}
				arrow := arrowData{tower: item.Pos, target: enemy.Pos, missed: false, shot: curTime}
				board.arrowsShot = append(board.arrowsShot, arrow)

				dmgAmt := 5 + rand.Intn(20)
				enemy.Health -= dmgAmt

				if enemy.Health <= 0 {
					board.enemyMap[enemy.Pos].dead = true
					//For tweening
					board.enemyMap[enemy.Pos].OldPos = board.enemyMap[enemy.Pos].Pos
				}
				break
			}
		}
	}
}
