package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"time"

	grid "github.com/mikepepping/didactic-enigma/grid"
)

var (
	width  int
	height int
	state  grid.Grid
	isDemo bool
)

func main() {
	flag.IntVar(&width, "width", 10, "Width of the grid")
	flag.IntVar(&height, "height", 10, "Height of the grid")
	flag.BoolVar(&isDemo, "demo", false, "Run the demo")
	flag.Parse()

	// Consider making this a linked list queue so we can pop off really old states
	// for now we will consider the last state as the current state
	state = grid.New(width, height)

	if isDemo {
		initDemo(&state)
	}

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
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			status := g.Cells[x+(y*g.Width)]
			if status == 1 {
				fmt.Print("*")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

func initDemo(g *grid.Grid) error {
	if g.Width < 3 || g.Height < 3 {
		return errors.New("grid is too small for demo")
	}

	// adds a walker
	//  0,1,0
	//  0,0,1
	//  1,1,1
	alive := [][]int{{1, 0}, {2, 1}, {0, 2}, {1, 2}, {2, 2}}
	for _, pos := range alive {
		g.Cells[pos[0]+(pos[1]*g.Width)] = 1
	}
	return nil
}
