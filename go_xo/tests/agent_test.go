package xo_test

import (
	"reflect"
	"testing"
	xo "xo/go_xo"
)

func TestMinMaxAgentGameOver(t *testing.T) {
	agent := xo.MinMaxAgent{
		AgentMark:    xo.Naught,
		OpponentMark: xo.Cross,
	}

	game_over_observations := [][9]int{
		{1, 1, 1, 1, 1, 1, 1, 1, 1},
		{2, 2, 2, 2, 2, 2, 2, 2, 2},
		{1, 0, 0, 1, 0, 0, 1, 0, 0},
		{1, 0, 0, 0, 1, 0, 0, 0, 1},
		{0, 0, 1, 0, 1, 0, 1, 0, 0},
		{0, 0, 1, 0, 0, 1, 0, 0, 1},
		{1, 1, 1, 0, 0, 2, 0, 0, 2},
		{1, 2, 1, 2, 2, 2, 1, 1, 2},
		{1, 2, 2, 0, 1, 0, 0, 0, 1},
		{2, 1, 0, 1, 1, 0, 2, 2, 2},
		{1, 2, 1, 2, 1, 1, 1, 2, 2},
	}

	for _, observation := range game_over_observations {
		if !agent.GameOver(observation) {
			t.Errorf("Observation %v was supposed to be a game over!", observation)
		}
	}
}

func TestMinMaxAgentNotGameOver(t *testing.T) {
	agent := xo.MinMaxAgent{
		AgentMark:    xo.Naught,
		OpponentMark: xo.Cross,
	}

	game_not_over_observations := [][9]int{
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 1},
		{0, 0, 1, 0, 2, 0, 0, 0, 1},
	}

	for _, observation := range game_not_over_observations {
		if agent.GameOver(observation) {
			t.Errorf("Observation %v was supposed to not be game over!", observation)
		}
	}
}

func TestMinMaxAgentScoreBoardLosing(t *testing.T) {
	agent := xo.MinMaxAgent{
		AgentMark:    xo.Naught,
		OpponentMark: xo.Cross,
	}

	observations := [][9]int{
		{1, 2, 2, 0, 0, 2, 1, 1, 2},
		{2, 2, 1, 1, 2, 2, 1, 1, 2},
		{1, 2, 1, 2, 2, 2, 1, 1, 2},
		{0, 2, 2, 1, 0, 2, 1, 1, 2},
		{0, 2, 2, 0, 1, 2, 1, 1, 2},
		{2, 1, 1, 0, 2, 0, 0, 0, 2},
		{2, 1, 0, 1, 1, 0, 2, 2, 2},
	}
	correct_score := -10

	for _, observation := range observations {
		score := agent.ScoreBoard(observation, 0)

		if score != correct_score {
			t.Errorf("Observation %v had score %v! Was meat to be %v!", observation, score, correct_score)
		}
	}
}

func TestGetIndexOfEmptyCells(t *testing.T) {
	observations := [][9]int{
		{1, 0, 2, 0, 1, 0, 2, 0, 1},
		{0, 1, 0, 2, 0, 1, 0, 2, 0},
		{2, 0, 1, 0, 2, 0, 1, 0, 2},
		{0, 2, 0, 1, 0, 2, 0, 1, 0},
	}
	expected := [][]int{
		{1, 3, 5, 7},
		{0, 2, 4, 6, 8},
		{1, 3, 5, 7},
		{0, 2, 4, 6, 8},
	}

	for i := range observations {
		result := xo.GetIndexOfEmptyCells(observations[i])

		if !reflect.DeepEqual(result, expected[i]) {
			t.Errorf("Expected %v, but got %v", expected, result)
		}
	}
}
