package main

import (
	"log"
	"math/rand"
	"time"
)

const (
	TOWER_RED = iota
	TOWER_GREEN
	TOWER_BLUE
	TOWER_PURPLE
	TOWER_MAX
)

func addTower() {

	if *noTowers {
		return
	}

	//Tower constuction progresss
	for _, item := range board.towerMap {
		if item.building < 2 {
			item.building++
		}
	}

	notUpgradable := false

	//Get votes, or choose a point if no votes
	if votes.VoteCount > 0 &&
		votes.Result.X > 0 &&
		votes.Result.Y > 0 &&
		votes.Result.X <= boardSizeX &&
		votes.Result.Y <= boardSizeY {

		tpos := votes.Result
		//Handle placing tower on top of goblin
		moved := false
		if board.goblinMap[tpos] != nil {
			for x := tpos.X; x <= boardSizeX+enemyBoardX; x++ {
				newPos := xyi{X: x, Y: tpos.Y}
				if board.goblinMap[newPos] != nil {
					//Move to new pos
					board.goblinMap[newPos] = board.goblinMap[tpos]
					delete(board.goblinMap, tpos)
					log.Println("Tower placed on top of goblin, moving goblin.")
					moved = true
					break
				}
			}
			//Should never happen
			if !moved {
				log.Println("Unable to find new location for goblin! tower not placed!")
				return
			}
		}
		if board.towerMap[tpos] == nil {
			board.towerMap[tpos] = &objectData{
				pos:          tpos,
				sheetP:       &obj_tower1,
				health:       obj_tower1.health,
				aniOffset:    uint64(rand.Intn(obj_tower1.frames)),
				building:     0,
				worldObjType: OTYPE_TOWER}
		} else {
			//Tower upgrade
			if board.towerMap[tpos].upgrade < TOWER_PURPLE {
				board.towerMap[tpos].upgrade++
			} else {
				notUpgradable = true
			}
		}
		if !notUpgradable {
			return
		}
	}

	var foundSmart bool
	if *smartMove && board.moveNum != 0 {
		for x := 0; x < boardSizeX+enemyBoardX; x++ {
			if foundSmart {
				break
			}
			for y := 0; y < boardSizeY; y++ {
				enemy := board.goblinMap[xyi{X: x, Y: y}]
				if enemy == nil {
					continue
				}

				tpos := xyi{X: 0, Y: y}
				if x-5 > 0 && x <= boardSizeX {
					var found bool
					for xx := 0; xx < boardSizeX+enemyBoardX; xx++ {
						checkT := board.towerMap[xyi{X: xx, Y: y}]
						if checkT != nil && !checkT.dead {
							found = true
							break
						}
					}
					if !found {
						tpos = xyi{X: x - 5, Y: y}
					} else {
						continue
					}

				}

				tower := board.towerMap[tpos]
				checkForEnemy := board.goblinMap[tpos]
				if tower == nil && checkForEnemy == nil {
					board.towerMap[tpos] = &objectData{
						pos:          tpos,
						sheetP:       &obj_tower1,
						health:       obj_tower1.health,
						aniOffset:    uint64(rand.Intn(obj_tower1.frames)),
						building:     0,
						worldObjType: OTYPE_TOWER}
					foundSmart = true
					break
				}
			}
		}
	}
	//If smartMove isnt' enabled, or fails to find a point choose random
	if !*smartMove || !foundSmart {
		log.Println("Not enough votes, picking random.")
		for t := 0; t < 1000; t++ {
			tpos := xyi{X: rand.Intn(boardSizeX-1) + 1, Y: rand.Intn(boardSizeY-1) + 1}
			if board.goblinMap[tpos] == nil && board.towerMap[tpos] == nil {
				board.towerMap[tpos] = &objectData{
					pos:          tpos,
					sheetP:       &obj_tower1,
					health:       obj_tower1.health,
					aniOffset:    uint64(rand.Intn(obj_tower1.frames)),
					building:     0,
					worldObjType: OTYPE_TOWER}
				break
			}
		}
	}
}

func towerShootArrow() {
	curTime := time.Now()

	//Cycle list of towers
	for _, tower := range board.towerMap {
		//If tower is dead or not fully built, skip.
		if tower.dead || tower.building < 2 || tower.worldObjType != OTYPE_TOWER {
			continue
		}

		for count := 0; count < tower.upgrade+1; count++ {

			//Look for targets
			for _, enemy := range board.goblinMap {
				if enemy.dead {
					continue
				}

				if Distance(tower.pos, enemy.pos) >= 6 {
					continue
				}

				go func() {
					time.Sleep(time.Millisecond * time.Duration(rand.Intn(200)))
					playVariated(SND_ARROW_SHOOT)
				}()

				//50-50 hit or miss
				if rand.Intn(2) != 0 {
					//Missed
					arrow := arrowData{
						tower:  xyi{X: ((tower.pos.X + offX - 1) * mag), Y: ((tower.pos.Y) * mag) + 12},
						target: xyi{X: ((enemy.pos.X + offX - 1) * mag), Y: ((enemy.pos.Y) * mag) + 12},
						missed: true, shot: curTime}
					board.arrowsShot = append(board.arrowsShot, arrow)
					break
				}

				arrow := arrowData{
					tower:  xyi{X: ((tower.pos.X + offX - 1) * mag), Y: ((tower.pos.Y) * mag) + 12},
					target: xyi{X: ((enemy.pos.X + offX - 1) * mag), Y: ((enemy.pos.Y) * mag) + 12},
					missed: false, shot: curTime}
				board.arrowsShot = append(board.arrowsShot, arrow)

				//RNG damage
				dmgAmt := 5 + rand.Intn(15)
				enemy.health -= dmgAmt

				if enemy.health <= 0 {
					enemy.dead = true
					enemy.diedAt = time.Now()

					go func() {
						time.Sleep(deathDelay + (time.Millisecond * time.Duration(rand.Intn(200))))
						playVariated(SND_GOBLIN_DIE)
					}()

					//For tweening
					enemy.prevPos = enemy.pos
				}
				break
			}
		}
	}
}

func spawnGoblins() {
	//Every other move
	if board.moveNum%2 == 0 {
		rpos := xyi{X: boardSizeX + enemyBoardX + 1, Y: 1 + rand.Intn(boardSizeY-1)}
		if board.goblinMap[rpos] == nil {
			board.goblinMap[rpos] = &objectData{
				pos:       rpos,
				sheetP:    &obj_goblinBarb,
				health:    obj_goblinBarb.health,
				prevPos:   xyi{X: rpos.X, Y: rpos.Y},
				aniOffset: uint64(rand.Intn(obj_goblinBarb.frames))}
		}
	}
}

func goblinAttack() {

	var moveList []*objectData
	for _, goblin := range board.goblinMap {
		if goblin.dead {
			continue
		}

		//Setup next enemy position
		nextPos := goblin.pos
		nextPos.X -= 1

		//Check towers and enemy positions before moving
		tower := board.towerMap[nextPos]
		enemy := board.goblinMap[nextPos]
		if enemy != nil && !enemy.dead {
			//Because we aren't moving forward, set prevpos to pos so we don't teleport back and re-tween
			goblin.prevPos = goblin.pos
			continue
		}

		//If a tower is in our way, do damage
		goblin.attacking = false
		if tower != nil && !tower.dead && tower.building >= 2 {
			tower.health -= 10 + rand.Intn(10)
			goblin.attacking = true
			tower.lastAttacked = time.Now()

			go func() {
				time.Sleep(time.Millisecond*100 + time.Duration(rand.Intn(100)))
				playVariated(SND_AXE)
			}()
			if tower.health <= 0 {
				go playSound(SND_TOWER_DIE)
				tower.dead = true
			}
			if tower.worldObjType == OTYPE_VWALL {
				//Wall damaged, play a sound to alert players
				if time.Since(board.wallDmgTime) > time.Second*30 {
					board.wallDmgTime = time.Now()
					playSound(SND_TENSION)
				}
			}
			//Because we aren't moving forward, set prevpos to pos so we don't teleport back and re-tween
			goblin.prevPos = goblin.pos
			continue
		}

		//Delete enemy, add to list
		goblin.prevPos = goblin.pos
		goblin.pos = nextPos
		moveList = append(moveList, goblin)
		delete(board.goblinMap, goblin.prevPos)
	}

	//Add enemy back to new position
	for i, item := range moveList {
		board.goblinMap[item.pos] = moveList[i]
		if item.pos.X < -2 {
			board.gameover = GAME_DEFEAT
			playSound(SND_GAMEOVER)
			votes.RoundTime = time.Now()
			votes.GameRunning = false
		}
	}

}
