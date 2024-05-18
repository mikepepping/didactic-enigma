package main

import (
	"flag"
	"fmt"
	"time"

	grid "github.com/mikepepping/didactic-enigma/grid"
)

var (
	width  int
	height int
	state  grid.Grid
)

func main() {
	flag.IntVar(&width, "width", 10, "Width of the grid")
	flag.IntVar(&height, "height", 10, "Height of the grid")
	flag.Parse()

	// Consider making this a linked list queue so we can pop off really old states
	// for now we will consider the last state as the current state
	state = grid.New(width, height)

	for {
		state = next(state)
		print(state)
		time.Sleep(300 * time.Millisecond)
	}
}

func next(curr grid.Grid) grid.Grid {
	g := grid.New(curr.Width, curr.Height)
	for i := range curr.Cells {
		n, err := curr.Neighbours(i)
		if err != nil {
			panic(err)
		}

		alive := 0
		for _, cell := range n {
			if cell != 0 {
				alive++
			}
		}

		// 1. Any live cell with fewer than two live neighbors dies, as if by underpopulation.
		// 3. Any live cell with more than three live neighbors dies, as if by overpopulation.
		if alive < 2 || alive > 3 {
			g.Cells[i] = 0
			continue
		}

		// 2. Any live cell with two or three live neighbors lives on to the next generation.
		if curr.Cells[i] == 1 {
			g.Cells[i] = 1
			continue
		}

		// 4. Any dead cell with exactly three live neighbors becomes a live cell, as if by reproduction.
		if alive == 3 {
			g.Cells[i] = 1
			continue
		}

	}

	return g
}

func print(g grid.Grid) {
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			fmt.Print(g.Cells[x+(y*g.Width)])
		}
		fmt.Println()
	}
}
