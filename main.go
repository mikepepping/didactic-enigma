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

func main() {
	var (
		width  int
		height int
		isDemo bool
		curr   *grid.Grid
		grids  <-chan *grid.Grid
	)

	flag.IntVar(&width, "width", 10, "Width of the grid")
	flag.IntVar(&height, "height", 10, "Height of the grid")
	flag.BoolVar(&isDemo, "demo", false, "Run the demo")
	flag.Parse()

	curr = grid.New(width, height)

	if isDemo {
		initDemo(curr)
	}

	grids = grid.Generate(curr)

	for {
		if curr = <-grids; curr == nil {
			fmt.Println("an error occured generating grids")
			return
		}

		print(curr)
		time.Sleep(100 * time.Millisecond)
	}
}

func print(g *grid.Grid) {
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
