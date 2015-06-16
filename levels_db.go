package main

import (
	"bufio"
	"io"
	"os"
	"strings"
)

const game_data_file = "sokoban_levels.txt"

type ldb struct {
	maxLevel int
	data     [][]string
}

// convert strings to 2-d matrix with bytes
func (db *ldb) getLevel(l int) [][]Cell {
	board := make([][]Cell, len(db.data[l-1]))
	for i, str := range db.data[l-1] {
		board[i] = make([]Cell, len(str))
		for j := 0; j < len(str); j++ {
			var c Cell
			switch str[j] {
			case WALL:
				c = Cell{base: WALL, obj: WALL}
			case FLOOR:
				c = Cell{base: FLOOR, obj: FLOOR}
			case GUY:
				c = Cell{base: FLOOR, obj: GUY}
			case SLOT:
				c = Cell{base: SLOT, obj: SLOT}
			case BOX:
				c = Cell{base: FLOOR, obj: BOX}
			case BOXIN:
				c = Cell{base: SLOT, obj: BOX}
			}
			board[i][j] = c
		}
	}
	return board
}

// Read game database from .txt
func (db *ldb) loadAll() {
	f, err := os.Open(game_data_file)
	if err != nil {
		panic(err)
	}

	rd := bufio.NewReader(f)

	var l = 0
	var matrix = make([]string, 0)
	for {
		line, err := rd.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				db.data = append(db.data, matrix)
				db.maxLevel = len(db.data)
				return
			}
			panic(err)
		}
		line = strings.TrimRight(line, "\t\n\f\r")
		if len(line) == 0 {
			db.data = append(db.data, matrix)
			l = l + 1
			matrix = make([]string, 0)
		} else {
			matrix = append(matrix, line)
		}
	}
}
