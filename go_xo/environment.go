package xo

import (
	"fmt"
	"math/rand"
)

const (
	Empty = 0
	Naught = 1
	Cross = 2
)

type Board [3][3]int

type NaughtsAndCrossesEnvironment struct {
	Board
	timeStep int
}

func (b *NaughtsAndCrossesEnvironment) Reset() {
	b.Board = Board{
		{Empty, Empty, Empty},
		{Empty, Empty, Empty},
		{Empty, Empty, Empty},
	}
	b.timeStep = 0

	// Randomly decide if the agent goes first
	if rand.Intn(2) == 1 {
		row, col := b.randomSpace()
		b.Board[row][col] = Naught
	}
}

func (b *NaughtsAndCrossesEnvironment) randomSpace() (int, int) {
	row := rand.Intn(len(b.Board))
	col := rand.Intn(len(b.Board))
	return row, col
}

func (b *NaughtsAndCrossesEnvironment) agentTakeTurn() {
	// Select the first non-empty space
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if b.Board[i][j] == Empty {
				b.Board[i][j] = Naught
				return
			}
		}
	}
}

func (b *NaughtsAndCrossesEnvironment) Step(action int) (observation [9]int, reward int, terminated bool, truncated bool) {
    if !b.Terminated() {
        row := action / 3
        col := action % 3
        if b.Board[row][col] == Empty {
            b.Board[row][col] = Cross
        }
    }

    reward = b.Reward()
    terminated = b.Terminated()

    if !terminated {
        b.agentTakeTurn()
        b.timeStep += 1
    }

	observation = b.Observation()
	truncated = false

	return
}

func (b *NaughtsAndCrossesEnvironment) Render() {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			switch b.Board[i][j] {
			case Empty:
				fmt.Print(" ")
			case Naught:
				fmt.Print("O")
			case Cross:
				fmt.Print("X")
			}

			if j < 2 {
				fmt.Print("|")
			}
		}
		fmt.Println()
	}
}

func (b *NaughtsAndCrossesEnvironment) Observation() [9]int {
	obs := [9]int{}

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			obs[i*3+j] = b.Board[i][j]
		}
	}
	return obs
}

func (b *NaughtsAndCrossesEnvironment) gameIsDraw() bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if b.Board[i][j] == Empty {
                return false
            }
		}
	}

	return true
}

func (b *NaughtsAndCrossesEnvironment) gameWinner() int {
	for i := 0; i < 3; i++ {
		// Check for a win in the rows
        if b.Board[i][0] == b.Board[i][1] && b.Board[i][1] == b.Board[i][2] && b.Board[i][0] != Empty{
            return b.Board[i][0]
        }

        // Check for a win in the columns
        if b.Board[0][i] == b.Board[1][i] && b.Board[1][i] == b.Board[2][i] && b.Board[0][i] != Empty{
            return b.Board[0][i]
        }
    }

    // Check for a win in the diagonals
    if b.Board[0][0] == b.Board[1][1] && b.Board[1][1] == b.Board[2][2] && b.Board[0][0] != Empty{
        return b.Board[0][0]
    }

    if b.Board[0][2] == b.Board[1][1] && b.Board[1][1] == b.Board[2][0] && b.Board[0][2] != Empty{
        return b.Board[0][2]
    }

	return Empty
}

func (b *NaughtsAndCrossesEnvironment) Terminated() bool {
	return b.gameIsDraw() || b.gameWinner() != Empty
}

func (b *NaughtsAndCrossesEnvironment) Reward() int {
	switch b.gameWinner() {
		case Naught:
			return -10
		case Cross:
			return 10
		}

	if b.gameIsDraw() {
		return 1
	}

	return 0
}