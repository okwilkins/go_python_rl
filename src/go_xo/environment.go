package xo

import (
	"fmt"
	"math/rand"
)

const (
	Empty  = 0
	Naught = 1
	Cross  = 2
)

type Board [3][3]byte

type NaughtsAndCrossesEnvironment struct {
	Board
	TimeStep   byte
	UserMark   byte
	Agent      Agent
	LastPlayer byte
}

func (env *NaughtsAndCrossesEnvironment) Reset() {
	env.Board = Board{
		{Empty, Empty, Empty},
		{Empty, Empty, Empty},
		{Empty, Empty, Empty},
	}
	env.TimeStep = 0

	// Randomly decide if the agent goes first
	if rand.Intn(2) == 1 {
		observation := env.Observation()
		env.Agent.TakeAction(&observation)
		env.LastPlayer = *env.Agent.GetMark()
	}
}

func (env *NaughtsAndCrossesEnvironment) Render() {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			switch env.Board[i][j] {
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

func (env *NaughtsAndCrossesEnvironment) Observation() [9]byte {
	obs := [9]byte{
		env.Board[0][0], env.Board[0][1], env.Board[0][2],
		env.Board[1][0], env.Board[1][1], env.Board[1][2],
		env.Board[2][0], env.Board[2][1], env.Board[2][2],
	}
	return obs
}

func (env *NaughtsAndCrossesEnvironment) GameIsDraw() bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if env.Board[i][j] == Empty {
				return false
			}
		}
	}

	return true
}

func (env *NaughtsAndCrossesEnvironment) GameWinner() byte {
	for i := 0; i < 3; i++ {
		// Check for a win in the rows
		if env.Board[i][0] == env.Board[i][1] && env.Board[i][1] == env.Board[i][2] && env.Board[i][0] != Empty {
			return env.Board[i][0]
		}

		// Check for a win in the columns
		if env.Board[0][i] == env.Board[1][i] && env.Board[1][i] == env.Board[2][i] && env.Board[0][i] != Empty {
			return env.Board[0][i]
		}
	}

	// Check for a win in the diagonals
	if env.Board[0][0] == env.Board[1][1] && env.Board[1][1] == env.Board[2][2] && env.Board[0][0] != Empty {
		return env.Board[0][0]
	}

	if env.Board[0][2] == env.Board[1][1] && env.Board[1][1] == env.Board[2][0] && env.Board[0][2] != Empty {
		return env.Board[0][2]
	}

	return Empty
}

func (env *NaughtsAndCrossesEnvironment) Terminated() bool {
	return env.GameIsDraw() || env.GameWinner() != Empty
}

func (env *NaughtsAndCrossesEnvironment) Truncated() bool {
	return env.TimeStep > 9
}

func (env *NaughtsAndCrossesEnvironment) Reward() int {
	switch env.GameWinner() {
	case Naught:
		return -100
	case Cross:
		return 100
	}

	if env.GameIsDraw() {
		return 50
	}

	return 0
}

func (env *NaughtsAndCrossesEnvironment) Step(action byte) (
	observation [9]byte,
	reward int,
	terminated bool,
	truncated bool,
) {
	if !env.Terminated() {
		env.PlaceMarker(action, &env.UserMark)

		if !env.Terminated() {
			observation := env.Observation()
			agent_action := env.Agent.TakeAction(&observation)
			env.PlaceMarker(agent_action, env.Agent.GetMark())
		}
	}

	env.TimeStep += 1
	reward = env.Reward()
	terminated = env.Terminated()
	observation = env.Observation()
	truncated = env.Truncated()

	return
}

func (env *NaughtsAndCrossesEnvironment) PlaceMarker(
	action byte,
	player_mark *byte,
) {
	row := action / 3
	col := action % 3

	if env.Board[row][col] == Empty {
		env.Board[row][col] = *player_mark
	}
}
