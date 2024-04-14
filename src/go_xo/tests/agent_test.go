package xo_test

import (
	"reflect"
	"slices"
	"testing"
	xo "xo/src/go_xo"
)

func TestMinMaxAgentGameOver(t *testing.T) {
	agent := xo.MinMaxAgent{
		AgentMark:    xo.Naught,
		OpponentMark: xo.Cross,
	}

	game_over_observations := [][9]byte{
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
		{1, 2, 2, 0, 1, 0, 0, 0, 1},
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

	game_not_over_observations := [][9]byte{
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

	observations := [][9]byte{
		{1, 2, 2, 0, 0, 2, 1, 1, 2},
		{2, 2, 1, 1, 2, 2, 1, 1, 2},
		{1, 2, 1, 2, 2, 2, 1, 1, 2},
		{0, 2, 2, 1, 0, 2, 1, 1, 2},
		{0, 2, 2, 0, 1, 2, 1, 1, 2},
		{2, 1, 1, 0, 2, 0, 0, 0, 2},
		{2, 1, 0, 1, 1, 0, 2, 2, 2},
	}
	var correct_score int8 = -10

	for _, observation := range observations {
		score := agent.ScoreBoard(observation, 0)

		if score != correct_score {
			t.Errorf("Observation %v had score %v! Was meant to be %v!", observation, score, correct_score)
		}
	}
}

func TestMinMaxAgentScoreBoardWinning(t *testing.T) {
	agent := xo.MinMaxAgent{
		AgentMark:    xo.Naught,
		OpponentMark: xo.Cross,
	}

	observations := [][9]byte{
		{1, 2, 2, 0, 1, 0, 0, 0, 1},
	}
	var correct_score int8 = 10

	for _, observation := range observations {
		score := agent.ScoreBoard(observation, 0)

		if score != correct_score {
			t.Errorf("Observation %v had score %v! Was meant to be %v!", observation, score, correct_score)
		}
	}
}

func TestMinMaxAgentPossibleGames(t *testing.T) {
	agent := xo.MinMaxAgent{
		AgentMark:    xo.Naught,
		OpponentMark: xo.Cross,
	}

	observations := [][9]byte{
		{1, 1, 0, 0, 0, 2, 0, 0, 2},
		{1, 2, 1, 0, 0, 2, 1, 1, 2},
		{1, 2, 1, 2, 2, 2, 1, 1, 2},
	}

	expected_possible_games := [][][9]byte{
		{
			{1, 1, 1, 0, 0, 2, 0, 0, 2},
			{1, 1, 0, 1, 0, 2, 0, 0, 2},
			{1, 1, 0, 0, 1, 2, 0, 0, 2},
			{1, 1, 0, 0, 0, 2, 1, 0, 2},
			{1, 1, 0, 0, 0, 2, 0, 1, 2},
		},
		{
			{1, 2, 1, 1, 0, 2, 1, 1, 2},
			{1, 2, 1, 0, 1, 2, 1, 1, 2},
		},
		nil,
	}

	for i := range observations {
		result := agent.GetPossibleGames(observations[i], agent.AgentMark)

		if !reflect.DeepEqual(result, expected_possible_games[i]) {
			t.Errorf("Expected %v, but got %v", expected_possible_games[i], result)
		}
	}
}

func TestMinMaxAgentNextPlayer(t *testing.T) {
	agent := xo.MinMaxAgent{
		AgentMark:    xo.Naught,
		OpponentMark: xo.Cross,
	}

	if agent.GetNextPlayer(xo.Cross) != xo.Naught {
		t.Errorf("Player %v was supposed to be %v!", xo.Cross, xo.Naught)
	}

	if agent.GetNextPlayer(xo.Naught) != xo.Cross {
		t.Errorf("Player %v was supposed to be %v!", xo.Naught, xo.Cross)
	}
}

func TestMinMaxAgentAlgorithm(t *testing.T) {
	agent := xo.MinMaxAgent{
		AgentMark:    xo.Naught,
		OpponentMark: xo.Cross,
	}

	observations := [][9]byte{
		{1, 2, 0, 0, 0, 2, 1, 1, 2},
		{0, 2, 1, 0, 0, 2, 1, 1, 2},
		{0, 2, 0, 1, 0, 2, 1, 1, 2},
		{0, 2, 0, 0, 1, 2, 1, 1, 2},
		{0, 2, 0, 0, 0, 2, 1, 1, 2},
		{1, 2, 2, 0, 0, 0, 0, 0, 1},
		{2, 1, 1, 0, 0, 0, 0, 0, 2},
		{1, 2, 2, 0, 0, 2, 1, 1, 2},
		{1, 2, 2, 0, 0, 2, 1, 1, 2},
		{2, 0, 1, 0, 0, 0, 1, 0, 2},
		{2, 0, 1, 0, 2, 0, 0, 0, 2},
		{2, 1, 0, 1, 1, 0, 2, 2, 2},
		{2, 1, 0, 1, 1, 0, 0, 2, 2},
	}

	correct_scores := []int8{-8, -6, -8, -8, -6, 9, -9, -10, -10, -9, -10, -10, -9}
	depths := []byte{1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	current_players := []byte{
		xo.Cross,
		xo.Cross,
		xo.Cross,
		xo.Cross,
		xo.Naught,
		xo.Naught,
		xo.Cross,
		xo.Naught,
		xo.Cross,
		xo.Cross,
		xo.Naught,
		xo.Cross,
		xo.Cross,
	}

	for i := range observations {
		score := agent.MinMax(observations[i], depths[i], current_players[i])

		if score != correct_scores[i] {
			t.Errorf(
				"Observation %v had score %v! Was meant to be %v! Current player: %v, Depth: %v!",
				observations[i],
				score,
				correct_scores[i],
				current_players[i],
				depths[i],
			)
		}
	}
}

func TestMinMaxAgentGetMinMaxBestMoves(t *testing.T) {
	agent := xo.MinMaxAgent{
		AgentMark:    xo.Naught,
		OpponentMark: xo.Cross,
	}

	observations := [][9]byte{
		{1, 2, 1, 0, 2, 2, 1, 1, 2},
		{0, 2, 1, 1, 2, 2, 1, 1, 2},
		{0, 2, 0, 0, 1, 2, 1, 1, 2},
		{1, 0, 2, 2, 1, 0, 0, 1, 0},
		{1, 0, 2, 0, 0, 0, 0, 0, 1},
		{1, 2, 0, 0, 2, 0, 0, 1, 1},
	}
	correct_moves := [][]byte{{3}, {0}, {2}, {1, 8}, {4}, {6}}

	for i := range observations {
		actions := agent.GetMinMaxBestMoves(observations[i])

		for _, action := range actions {
			if !slices.Contains(correct_moves[i], action) {
				t.Errorf(
					"Observation %v had action %v! Correct actions: %v!",
					observations[i],
					action,
					correct_moves[i],
				)
			}
		}
	}
}

func TestGetIndexOfEmptyCells(t *testing.T) {
	observations := [][9]byte{
		{1, 0, 2, 0, 1, 0, 2, 0, 1},
		{0, 1, 0, 2, 0, 1, 0, 2, 0},
		{2, 0, 1, 0, 2, 0, 1, 0, 2},
		{0, 2, 0, 1, 0, 2, 0, 1, 0},
		{2, 0, 1, 0, 0, 0, 1, 2, 1},
	}
	expected := [][]byte{
		{1, 3, 5, 7},
		{0, 2, 4, 6, 8},
		{1, 3, 5, 7},
		{0, 2, 4, 6, 8},
		{1, 3, 4, 5},
	}

	for i := range observations {
		result := xo.GetIndexOfEmptyCells(observations[i])

		if !reflect.DeepEqual(result, expected[i]) {
			t.Errorf("Expected %v, but got %v", expected, result)
		}
	}
}
