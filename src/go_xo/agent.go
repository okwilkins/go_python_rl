package xo

import (
	"fmt"
	"math/rand"
	"slices"
)

type Agent interface {
	TakeAction(observation *[9]byte) (byte, error)
	GetMark() byte
}

type RandomAgent struct {
	Mark byte
}

func (a *RandomAgent) TakeAction(observation *[9]byte) (byte, error) {
	empty_cells := GetIndexOfEmptyCells(observation)
	if len(empty_cells) == 0 {
		return 0, fmt.Errorf("no empty cell found in observation: %v", observation)
	}
	return empty_cells[rand.Intn(len(empty_cells))], nil
}

func (a *RandomAgent) GetMark() byte {
	return a.Mark
}

type FillFirstEmptyAgent struct {
	Mark byte
}

func (a *FillFirstEmptyAgent) TakeAction(observation *[9]byte) (byte, error) {
	empty_cells := GetIndexOfEmptyCells(observation)

	for _, cell := range empty_cells {
		if observation[cell] == Empty {
			return cell, nil
		}
	}

	return 0, fmt.Errorf("no empty cell found in observation: %v", observation)
}

func (a *FillFirstEmptyAgent) GetMark() byte {
	return a.Mark
}

type MinMaxAgent struct {
	AgentMark    byte
	OpponentMark byte
}

func (a *MinMaxAgent) TakeAction(observation *[9]byte) (byte, error) {
	first_turn := true

	for _, cell := range observation {
		if cell != Empty {
			first_turn = false
			break
		}
	}

	if first_turn {
		// Take random action to avoid heavy computation
		return a.TakeRandomAction(observation)
	} else {
		return a.TakeMinMaxAction(observation)
	}
}

func (a *MinMaxAgent) TakeRandomAction(observation *[9]byte) (byte, error) {
	empty_cells := GetIndexOfEmptyCells(observation)
	if len(empty_cells) == 0 {
		return 0, fmt.Errorf("no empty cells were found in observation: %v", observation)
	}
	return empty_cells[rand.Intn(len(empty_cells))], nil
}

func (a *MinMaxAgent) TakeMinMaxAction(observation *[9]byte) (byte, error) {
	best_moves := a.GetMinMaxBestMoves(observation)

	if len(best_moves) == 0 {
		return 0, fmt.Errorf(
			"min max algorithm could not find moves as there are no available moves in observation: %v",
			observation,
		)
	}

	randIdx := rand.Intn(len(best_moves))
	return best_moves[randIdx], nil
}

func (a *MinMaxAgent) GetMinMaxBestMoves(observation *[9]byte) []byte {
	var best_score int8 = -10
	var best_moves []byte

	for _, cell := range GetIndexOfEmptyCells(observation) {
		possible_game := *observation
		possible_game[cell] = a.AgentMark

		next_player := a.GetNextPlayer(a.AgentMark)
		score := a.MinMax(&possible_game, 0, next_player)

		if score == best_score {
			best_moves = append(best_moves, cell)
		} else if score > best_score {
			best_score = score
			best_moves = []byte{cell}
		}
	}

	return best_moves
}

func (a *MinMaxAgent) GameOver(observation *[9]byte) bool {
	for i := 0; i < 3; i++ {
		// Check for win in the rows
		if observation[i*3] == observation[i*3+1] &&
			observation[i*3+1] == observation[i*3+2] &&
			observation[i*3+2] != Empty {
			return true
		}

		// Check for win in the columns
		if observation[i] == observation[i+3] &&
			observation[i+3] == observation[i+6] &&
			observation[i+6] != Empty {
			return true
		}
	}

	// Check for win in diagonals
	if observation[0] == observation[4] &&
		observation[4] == observation[8] &&
		observation[8] != Empty {
		return true
	}

	if observation[2] == observation[4] &&
		observation[4] == observation[6] &&
		observation[6] != Empty {
		return true
	}

	// If all the win conditions are not satisfied, check if a draw as occurred
	game_is_draw := true

	for _, cell := range observation {
		if cell == Empty {
			game_is_draw = false
		}
	}

	return game_is_draw
}

func (a *MinMaxAgent) ScoreBoard(observation *[9]byte, depth byte) int8 {
	var score int8 = 0

	for i := 0; i < 3; i++ {
		// Check for win in the rows
		if observation[i*3] == observation[i*3+1] &&
			observation[i*3+1] == observation[i*3+2] &&
			observation[i*3+2] != Empty {
			score = a.GetMinMaxScoreMap()[observation[i*3]]
			break
		}

		// Check for win in the columns
		if observation[i] == observation[i+3] &&
			observation[i+3] == observation[i+6] &&
			observation[i+6] != Empty {
			score = a.GetMinMaxScoreMap()[observation[i]]
			break
		}
	}

	if score == 0 {
		// Check for a win in the diagonals
		if observation[0] == observation[4] &&
			observation[4] == observation[8] &&
			observation[8] != Empty {
			score = a.GetMinMaxScoreMap()[observation[0]]
		} else if observation[2] == observation[4] &&
			observation[4] == observation[6] &&
			observation[6] != Empty {
			score = a.GetMinMaxScoreMap()[observation[2]]
		}
	}

	switch score {
	case 10:
		// The player won
		return score - int8(depth)
	case -10:
		// The player lost
		return score + int8(depth)
	default:
		// The game is a draw
		return 0
	}
}

func (a *MinMaxAgent) GetPossibleGames(observation *[9]byte, next_player byte) [][9]byte {
	var possible_games [][9]byte

	for _, cell := range GetIndexOfEmptyCells(observation) {
		possible_game := *observation
		possible_game[cell] = next_player
		possible_games = append(possible_games, possible_game)
	}

	return possible_games
}

func (a *MinMaxAgent) GetNextPlayer(current_player byte) byte {
	if current_player == Cross {
		return Naught
	} else {
		return Cross
	}
}

func (a *MinMaxAgent) MinMax(observation *[9]byte, depth byte, current_player byte) int8 {
	// If game over, return the score
	if a.GameOver(observation) {
		return a.ScoreBoard(observation, depth)
	}

	depth++
	var scores []int8
	possible_games := a.GetPossibleGames(observation, current_player)

	for _, possible_game := range possible_games {
		next_player := a.GetNextPlayer(current_player)
		possible_score := a.MinMax(&possible_game, depth, next_player)
		scores = append(scores, possible_score)
	}

	if current_player == a.AgentMark {
		return slices.Max(scores)
	} else {
		return slices.Min(scores)
	}
}

func (a *MinMaxAgent) GetMinMaxScoreMap() map[byte]int8 {
	score_map := make(map[byte]int8)
	score_map[a.AgentMark] = 10
	score_map[a.OpponentMark] = -10
	score_map[Empty] = 0

	return score_map
}

func (a *MinMaxAgent) GetMark() byte {
	return a.AgentMark
}

func GetIndexOfEmptyCells(observation *[9]byte) []byte {
	var empty_cells []byte
	for i, cell := range observation {
		if cell == Empty {
			empty_cells = append(empty_cells, byte(i))
		}
	}
	return empty_cells
}
