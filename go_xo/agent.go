package xo

import (
	"math/rand"
	"slices"
)

// TODO: use pointers for the agents

type Agent interface {
	TakeAction(observation [9]int) int
	GetMark() int
}

type RandomAgent struct {
	Mark int
}

func (a *RandomAgent) TakeAction(observation [9]int) int {
	empty_cells := GetIndexOfEmptyCells(observation)
	return empty_cells[rand.Intn(len(empty_cells))]
}

func (a *RandomAgent) GetMark() int {
	return a.Mark
}

type FillFirstEmptyAgent struct {
	Mark int
}

func (a *FillFirstEmptyAgent) TakeAction(observation [9]int) int {
	empty_cells := GetIndexOfEmptyCells(observation)

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

type MinMaxAgent struct {
	AgentMark    int
	OpponentMark int
}

func (a *MinMaxAgent) TakeAction(observation [9]int) int {
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

func (a *MinMaxAgent) TakeRandomAction(observation [9]int) int {
	empty_cells := GetIndexOfEmptyCells(observation)
	return empty_cells[rand.Intn(len(empty_cells))]
}

func (a *MinMaxAgent) TakeMinMaxAction(observation [9]int) int {
	best_score := -10
	var best_moves []int

	for _, cell := range GetIndexOfEmptyCells(observation) {
		possible_game := observation
		possible_game[cell] = a.AgentMark

		next_player := a.GetNextPlayer(a.AgentMark)
		score := a.MinMax(possible_game, 0, next_player)

		if score == best_score {
			best_moves = append(best_moves, cell)
		} else if score > best_score {
			best_score = score
			best_moves = []int{cell}
		}
	}

	randIdx := rand.Intn(len(best_moves))
	return best_moves[randIdx]
}

func (a *MinMaxAgent) GameOver(observation [9]int) bool {
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

func (a *MinMaxAgent) ScoreBoard(observation [9]int, depth int) int {
	score := 0

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
	case 0:
		return score - depth
	case -10:
		// The player lost
		return score + depth
	default:
		// The game is a draw
		return 0
	}
}

func (a *MinMaxAgent) GetPossibleGames(observation [9]int, next_player int) [][9]int {
	var possible_games [][9]int

	for _, cell := range GetIndexOfEmptyCells(observation) {
		possible_game := observation
		possible_game[cell] = next_player
		possible_games = append(possible_games, possible_game)
	}

	return possible_games
}

func (a *MinMaxAgent) GetNextPlayer(current_player int) int {
	if current_player == Cross {
		return Naught
	} else {
		return Cross
	}
}

func (a *MinMaxAgent) MinMax(observation [9]int, depth int, current_player int) int {
	// If game over, return the score
	if a.GameOver(observation) {
		return a.ScoreBoard(observation, depth)
	}

	depth++
	var scores []int
	possible_games := a.GetPossibleGames(observation, current_player)

	for _, pospossible_game := range possible_games {
		next_player := a.GetNextPlayer(current_player)
		possible_score := a.MinMax(pospossible_game, depth, next_player)
		scores = append(scores, possible_score)
	}

	if current_player == a.AgentMark {
		return slices.Max(scores)
	} else {
		return slices.Min(scores)
	}
}

func (a *MinMaxAgent) GetMinMaxScoreMap() map[int]int {
	score_map := make(map[int]int)
	score_map[a.AgentMark] = 10
	score_map[a.OpponentMark] = -10
	score_map[Empty] = 0

	return score_map
}

func (a *MinMaxAgent) GetMark() int {
	return a.AgentMark
}

func GetIndexOfEmptyCells(observation [9]int) []int {
	var empty_cells []int
	for i, cell := range observation {
		if cell == Empty {
			empty_cells = append(empty_cells, i)
		}
	}
	return empty_cells
}
