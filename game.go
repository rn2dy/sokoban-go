package main

import (
	_ "fmt"
)

type Direction int

const (
	UP Direction = iota
	DOWN
	LEFT
	RIGHT
)

var baseTokens = []byte("#@XO* ")

var (
	WALL  = baseTokens[0]
	GUY   = baseTokens[1]
	SLOT  = baseTokens[2]
	BOX   = baseTokens[3]
	BOXIN = baseTokens[4]
	FLOOR = baseTokens[5]
)

type Cell struct {
	base, obj byte
}

type Game struct {
	level int
	board [][]Cell
	steps [][][]Cell
	db    *ldb
	x     int
	y     int
	debug bool
}

func NewGame() *Game {
	g := &Game{
		level: 1,
		db:    &ldb{},
	}
	g.db.loadAll()
	g.reset()
	return g
}

// reset current level
func (g *Game) reset() {
	g.board = g.db.getLevel(g.level)
	g.whereami()
}

func (g *Game) record() {
	step := make([][]Cell, len(g.board))
	for i, cs := range g.board {
		step[i] = make([]Cell, len(cs))
		for j, c := range cs {
			step[i][j] = c
		}
	}
	g.steps = append(g.steps, step)
	// record the last 100 steps
	if len(g.steps) > 100 {
		g.steps = g.steps[len(g.steps)-100:]
	}
}

func (g *Game) undo() {
	if len(g.steps) > 0 {
		g.board = g.steps[len(g.steps)-1]
		g.whereami()
		g.steps = g.steps[:len(g.steps)-1]
	}
}

func (g *Game) whereami() {
	for y, cels := range g.board {
		for x, cel := range cels {
			if cel.obj == GUY {
				g.x, g.y = x, y
			}
		}
	}
}

func (g *Game) checkMove(dx, dy int) {
	var (
		x, y int
		k    = 1
	)
	for {
		x, y = g.x+dx*k, g.y+dy*k
		cell := g.board[y][x]
		switch cell.obj {
		case FLOOR, SLOT:
			g.record()
			// move obj along the way forward
			for k > 0 {
				xp, yp := x-dx, y-dy
				g.board[y][x].obj = g.board[yp][xp].obj
				x, y = xp, yp
				k--
			}
			g.board[y][x].obj = g.board[y][x].base
			g.x, g.y = x+dx, y+dy
			return
		case WALL:
			return
		case BOX:
			if g.board[y+dy][x+dx].obj == BOX {
				return
			}
		}
		k++
	}
}

func (g *Game) move(dir Direction) {
	switch dir {
	case UP:
		g.checkMove(0, -1)
	case DOWN:
		g.checkMove(0, 1)
	case LEFT:
		g.checkMove(-1, 0)
	case RIGHT:
		g.checkMove(1, 0)
	}
	done := g.checkState()
	if done {
		g.nextLevel()
		g.steps = g.steps[0:0]
	}
}

func (g *Game) checkState() bool {
	for _, cells := range g.board {
		for _, cell := range cells {
			if cell.base == SLOT && cell.obj != BOX {
				return false
			}
		}
	}
	return true
}

func (g *Game) nextLevel() {
	if g.level < g.db.maxLevel {
		g.level = g.level + 1
		g.reset()
	}
}

func (g *Game) prevLevel() {
	if g.level > 1 {
		g.level = g.level - 1
		g.reset()
	}
}

func (g *Game) toggleDebug() {
	g.debug = !g.debug
}

// func main() {
// 	var p = func(g *Game) {
// 		for i, row := range g.board {
// 			fmt.Printf("%-2d %q\n", i, row)
// 		}
// 		fmt.Println()
// 	}
//
// 	g := NewGame()
// 	p(g)
// 	g.move(RIGHT)
// 	p(g)
// 	g.undo()
// 	p(g)
// }
