package main

type spriteSheet struct {
	name   string
	file   string
	size   xyi
	frames int
	anims  []animationData
}

type animationData struct {
	name   string
	row    int
	mirror bool
}

var goblinBarb = spriteSheet{
	name:   "goblin barbarian",
	file:   "Goblin_Barbarian",
	size:   xyi{X: 32, Y: 32},
	frames: 4,

	anims: []animationData{
		{
			name:   "idle",
			row:    2,
			mirror: true,
		},
		{
			name:   "run",
			row:    5,
			mirror: true,
		},
		{
			name:   "fade",
			row:    8,
			mirror: true,
		},
		{
			name:   "die",
			row:    10,
			mirror: true,
		},
		{
			name:   "attack",
			row:    13,
			mirror: true,
		},
		{
			name:   "swim",
			row:    16,
			mirror: true,
		},
	},
}
