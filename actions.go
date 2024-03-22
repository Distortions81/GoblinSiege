package main

import (
	"log"
	"math/rand"
	"time"
)

func addTower() {

	if *noTowers {
		return
	}

	//Increment build step
	for _, item := range board.playMap {
		if item.building < 2 {
			item.building++
		}
	}

	//Get votes, or choose random point if no votes
	if votes.VoteCount > 0 &&
		votes.Result.X > 0 &&
		votes.Result.Y > 0 &&
		votes.Result.X <= boardSizeX &&
		votes.Result.Y <= boardSizeY {

		tpos := votes.Result
		if board.enemyMap[tpos] == nil && board.playMap[tpos] == nil {
			board.playMap[tpos] = &objectData{
				pos:          tpos,
				sheetP:       &obj_tower1,
				health:       obj_tower1.health,
				aniOffset:    uint64(rand.Intn(obj_tower1.frames)),
				building:     0,
				worldObjType: OTYPE_TOWER}
		} else {
			log.Println("COLLISION!")
		}
	} else {
		var foundSmart bool

		if *smartMove && board.moveNum != 0 {
			for x := 0; x < boardSizeX+enemyBoardX; x++ {
				if foundSmart {
					break
				}
				for y := 0; y < boardSizeY; y++ {
					enemy := board.enemyMap[xyi{X: x, Y: y}]
					if enemy == nil {
						continue
					}

					tpos := xyi{X: 0, Y: y}
					if x-5 > 0 && x <= boardSizeX {
						var found bool
						for xx := 0; xx < boardSizeX+enemyBoardX; xx++ {
							checkT := board.playMap[xyi{X: xx, Y: y}]
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

					tower := board.playMap[tpos]
					checkForEnemy := board.enemyMap[tpos]
					if tower == nil && checkForEnemy == nil {
						board.playMap[tpos] = &objectData{
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
		if !*smartMove || !foundSmart {
			log.Println("Not enough votes, picking random.")
			tpos := xyi{X: rand.Intn(boardSizeX-1) + 1, Y: rand.Intn(boardSizeY-1) + 1}
			if board.enemyMap[tpos] == nil && board.playMap[tpos] == nil {
				board.playMap[tpos] = &objectData{
					pos:          tpos,
					sheetP:       &obj_tower1,
					health:       obj_tower1.health,
					aniOffset:    uint64(rand.Intn(obj_tower1.frames)),
					building:     0,
					worldObjType: OTYPE_TOWER}
			}
		}
	}

}

func towerShootArrow() {
	curTime := time.Now()

	for _, item := range board.playMap {
		//If tower is dead or not fully built, skip.
		if item.dead || item.building < 2 || item.worldObjType != OTYPE_TOWER {
			continue
		}

		//Shoot from tower top
		towerPos := item.pos
		towerPos.Y -= 1

		for _, enemy := range board.enemyMap {
			if enemy.dead {
				continue
			}

			//If enemy within range
			if Distance(item.pos, enemy.pos) < 6 {
				go func() {
					time.Sleep(time.Millisecond * time.Duration(rand.Intn(200)))
					playVariated(SND_ARROW_SHOOT)
				}()

				if rand.Intn(2) != 0 {
					arrow := arrowData{tower: towerPos, target: enemy.pos, missed: true, shot: curTime}
					board.arrowsShot = append(board.arrowsShot, arrow)
					break
				}
				arrow := arrowData{tower: towerPos, target: enemy.pos, missed: false, shot: curTime}
				board.arrowsShot = append(board.arrowsShot, arrow)

				dmgAmt := 5 + rand.Intn(15)
				enemy.health -= dmgAmt

				if enemy.health <= 0 {
					board.enemyMap[enemy.pos].dead = true
					board.enemyMap[enemy.pos].diedAt = time.Now()
					go func() {
						time.Sleep(deathDelay + (time.Millisecond * time.Duration(rand.Intn(200))))
						playVariated(SND_GOBLIN_DIE)
					}()

					//For tweening
					board.enemyMap[enemy.pos].prevPos = board.enemyMap[enemy.pos].pos
				}
				break
			}

		}
	}
}

func spawnGoblins() {
	//Spawn goblins
	if board.moveNum%2 == 0 {
		rpos := xyi{X: boardSizeX + enemyBoardX + 1, Y: 1 + rand.Intn(boardSizeY-1)}
		if board.enemyMap[rpos] == nil {
			board.enemyMap[rpos] = &objectData{
				pos:       rpos,
				sheetP:    &obj_goblinBarb,
				health:    obj_goblinBarb.health,
				prevPos:   xyi{X: rpos.X, Y: rpos.Y},
				aniOffset: uint64(rand.Intn(obj_goblinBarb.frames))}
		}
	}
}

func goblinAttack() {

	var newitems []*objectData
	//Detect defeat, defeat, do damage to towers, remove dead towers
	for _, item := range board.enemyMap {
		if item.dead {
			continue
		}

		//Detect game over
		oldItem := item

		//Setup next enemy position
		nextPos := item.pos
		oldItem.prevPos = oldItem.pos
		nextPos.X -= 1

		//Check towers and enemy positions before moving
		tower := board.playMap[nextPos]
		self := board.enemyMap[nextPos]
		if self != nil && !self.dead {
			continue
		}

		//If a tower is in our way, do damage
		item.attacking = false
		if tower != nil && !tower.dead && tower.building >= 2 {
			tower.health -= 10 + rand.Intn(10)
			item.attacking = true
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
			continue
		}

		//Delete enemy, add to list
		delete(board.enemyMap, item.pos)
		oldItem.pos = nextPos
		newitems = append(newitems, oldItem)
	}

	//Add enemy back to new position
	for i, item := range newitems {
		board.enemyMap[item.pos] = newitems[i]
		if item.pos.X < -2 {
			board.gameover = GAME_DEFEAT
			playSound(SND_GAMEOVER)
			votes.RoundTime = time.Now()
			votes.GameRunning = false
		}
	}

}
