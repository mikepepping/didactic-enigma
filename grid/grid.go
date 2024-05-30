package grid

import (
	"errors"
)

type Grid struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Cells  []byte `json:"cells"`
}

func New(width int, height int) *Grid {
	return &Grid{
		Width:  width,
		Height: height,
		Cells:  make([]byte, width*height),
	}
}

func (g Grid) Neighbours(i int) ([]byte, error) {
	if i >= len(g.Cells) {
		return []byte{}, errors.New("index out of bounds")
	}
	// it kind of sucks that I need to make a new slice for each cell
	// think about using a channel or appending to a user provided slice instead

	var (
		n = []byte{}
		x = i % g.Width
		y = i / g.Width
	)

	// this stores the coordinates of the would be neighbours
	// we can use this to check if they are out of bounds before
	// calculating the index of the cell
	neigbourIndexs := [][]int{
		{x - 1, y - 1}, {x, y - 1}, {x + 1, y - 1},
		{x - 1, y}, {x + 1, y},
		{x - 1, y + 1}, {x, y + 1}, {x + 1, y + 1},
	}

	for _, ni := range neigbourIndexs {
		if ni[1] < 0 || ni[1] >= g.Height || ni[0] < 0 || ni[0] >= g.Width {
			continue
		}

		n = append(n, g.Cells[ni[0]+(ni[1]*g.Width)])
	}

	return n, nil
}

func (g Grid) Next() (*Grid, error) {
	next := New(g.Width, g.Height)
	for i := range g.Cells {
		n, err := g.Neighbours(i)
		if err != nil {
			return nil, err
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
			next.Cells[i] = 0
			continue
		}

		// 2. Any live cell with two or three live neighbors lives on to the next generation.
		if g.Cells[i] == 1 {
			next.Cells[i] = 1
			continue
		}

		// 4. Any dead cell with exactly three live neighbors becomes a live cell, as if by reproduction.
		if alive == 3 {
			next.Cells[i] = 1
			continue
		}

	}

	return next, nil
}

func Generate(start *Grid) <-chan *Grid {
	grids := make(chan *Grid)
	go func(start *Grid, grids chan *Grid) {
		curr := start
		for {
			next, err := curr.Next()
			if err != nil {
				panic(err)
			}
			grids <- next
			curr = next
		}
	}(start, grids)

	return grids
}
