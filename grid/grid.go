package grid

import (
	"errors"
)

type Grid struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Cells  []byte `json:"cells"`
}

func New(width int, height int) Grid {
	return Grid{
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
		if ni[1] > 0 || ni[1] >= g.Height || ni[0] < 0 || ni[0] >= g.Width {
			continue
		}

		n = append(n, g.Cells[ni[0]+(ni[1]*g.Width)])
	}

	return n, nil
}
