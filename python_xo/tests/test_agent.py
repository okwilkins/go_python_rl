import pytest
from python_xo.agent import MinMaxAgent
from python_xo.owner import Owner


@pytest.mark.parametrize(
    ("observation"),
    [
        (1, 1, 1, 1, 1, 1, 1, 1, 1),
        (2, 2, 2, 2, 2, 2, 2, 2, 2),
        (1, 0, 0, 1, 0, 0, 1, 0, 0),
        (1, 0, 0, 0, 1, 0, 0, 0, 1),
        (0, 0, 1, 0, 1, 0, 1, 0, 0),
        (0, 0, 1, 0, 0, 1, 0, 0, 1),
        (1, 1, 1, 0, 0, 2, 0, 0, 2),
        (1, 2, 1, 2, 2, 2, 1, 1, 2),
    ]
)
def test_min_max_agent_game_over(observation) -> None:
    agent = MinMaxAgent(agent_mark=Owner.NAUGHT)
    assert agent._game_over(observation=observation)


@pytest.mark.parametrize(
    ("observation"),
    [
        (0, 0, 0, 0, 0, 0, 0, 0, 0),
        (0, 0, 0, 0, 0, 0, 0, 0, 1),
        (0, 0, 1, 0, 2, 0, 0, 0, 1),
    ]
)
def test_min_max_agent_game_not_over(observation) -> None:
    agent = MinMaxAgent(agent_mark=Owner.NAUGHT)
    assert not agent._game_over(observation=observation)


@pytest.mark.parametrize(
    ("observation"),
    [
        (1, 2, 2, 0, 0, 2, 1, 1, 2),
        (2, 2, 1, 1, 2, 2, 1, 1, 2),
        (1, 2, 1, 2, 2, 2, 1, 1, 2),
        (0, 2, 2, 1, 0, 2, 1, 1, 2),
        (0, 2, 2, 0, 1, 2, 1, 1, 2),
    ]
)
def test_min_max_agent_score_board_losing(observation) -> None:
    agent = MinMaxAgent(agent_mark=Owner.NAUGHT)
    assert agent._score_board(observation=observation, depth=0) == -10


@pytest.mark.parametrize(
    ("observation", "correct_score", "depth", "current_player"),
    [
        ((1, 2, 0, 0, 0, 2, 1, 1, 2), -8, 1, Owner.CROSS),
        ((0, 2, 1, 0, 0, 2, 1, 1, 2), -6, 1, Owner.CROSS),
        ((0, 2, 0, 1, 0, 2, 1, 1, 2), -8, 1, Owner.CROSS),
        ((0, 2, 0, 0, 1, 2, 1, 1, 2), -8, 1, Owner.CROSS),
        ((0, 2, 0, 0, 0, 2, 1, 1, 2), -6, 0, Owner.NAUGHT),
    ]
)
def test_min_max_agent_algorithm(observation, correct_score, depth, current_player) -> None:
    agent = MinMaxAgent(agent_mark=Owner.NAUGHT)
    assert agent._min_max(
        observation=observation,
        depth=depth,
        current_player=current_player
    ) == correct_score


@pytest.mark.parametrize(
    ("observation", "correct_moves"),
    [
        ((1, 2, 1, 0, 2, 2, 1, 1, 2), [3]),
        ((0, 2, 1, 1, 2, 2, 1, 1, 2), [0]),
        ((0, 2, 0, 0, 1, 2, 1, 1, 2), [2]),
        ((2, 0, 1, 1, 2, 0, 0, 2, 0), [2, 8]),
    ]
)
def test_min_max_agent_correct_move(observation, correct_moves) -> None:
    agent = MinMaxAgent(agent_mark=Owner.CROSS)
    assert agent.take_min_max_action(observation=observation) in correct_moves

