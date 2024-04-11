package xo

import (
	"math/rand"
)

type Agent interface {
	TakeAction(observation [9]int) int
	GetMark() int
}

type RandomAgent struct {
	Mark int
}

func (a *RandomAgent) TakeAction(observation [9]int) int {
	empty_cells := getIndexOfEmptyCells(observation)
	return empty_cells[rand.Intn(len(empty_cells))]
}

func (a *RandomAgent) GetMark() int {
	return a.Mark
}

type FillFirstEmptyAgent struct {
	Mark int
}

func (a *FillFirstEmptyAgent) TakeAction(observation [9]int) int {
	empty_cells := getIndexOfEmptyCells(observation)

	for _, cell := range empty_cells {
		if observation[cell] == Empty {
			return cell
		}
	}

	panic("No empty cells found!")
}

func (a *FillFirstEmptyAgent) GetMark() int {
	return a.Mark
}

func getIndexOfEmptyCells(observation [9]int) []int {
	var empty_cells []int
	for i, cell := range observation {
		if cell == Empty {
			empty_cells = append(empty_cells, i)
		}
	}
	return empty_cells
}
