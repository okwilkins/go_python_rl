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
        (1, 2, 2, 0, 1, 0, 0, 0, 1),
        (2, 1, 0, 1, 1, 0, 2, 2, 2),
        (1, 2, 1, 2, 1, 1, 1, 2, 2),
    ],
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
    ],
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
        (2, 1, 1, 0, 2, 0, 0, 0, 2),
        (2, 1, 0, 1, 1, 0, 2, 2, 2),
    ],
)
def test_min_max_agent_score_board_losing(observation) -> None:
    agent = MinMaxAgent(agent_mark=Owner.NAUGHT)
    assert agent._score_board(observation=observation, depth=0) == -10


@pytest.mark.parametrize(
    ("observation", "possible_games"),
    [
        (
            (1, 1, 0, 0, 0, 2, 0, 0, 2),
            [
                (1, 1, 1, 0, 0, 2, 0, 0, 2),
                (1, 1, 0, 1, 0, 2, 0, 0, 2),
                (1, 1, 0, 0, 1, 2, 0, 0, 2),
                (1, 1, 0, 0, 0, 2, 1, 0, 2),
                (1, 1, 0, 0, 0, 2, 0, 1, 2),
            ],
        ),
        (
            (1, 2, 1, 0, 0, 2, 1, 1, 2),
            [
                (1, 2, 1, 1, 0, 2, 1, 1, 2),
                (1, 2, 1, 0, 1, 2, 1, 1, 2),
            ],
        ),
        ((1, 2, 1, 2, 2, 2, 1, 1, 2), []),
    ],
)
def test_min_max_get_possible_games(observation, possible_games) -> None:
    agent = MinMaxAgent(agent_mark=Owner.NAUGHT)
    calculated_games = agent._get_possible_games(
        observation=observation, next_player=Owner.NAUGHT
    )
    assert set(calculated_games) == set(possible_games)


def test_min_max_get_next_player() -> None:
    agent = MinMaxAgent(agent_mark=Owner.CROSS)
    assert agent._get_next_player(current_player=Owner.CROSS) == Owner.NAUGHT
    assert agent._get_next_player(current_player=Owner.NAUGHT) == Owner.CROSS


@pytest.mark.parametrize(
    ("observation", "correct_score", "depth", "current_player"),
    [
        ((1, 2, 0, 0, 0, 2, 1, 1, 2), -8, 1, Owner.CROSS),
        ((0, 2, 1, 0, 0, 2, 1, 1, 2), -6, 1, Owner.CROSS),
        ((0, 2, 0, 1, 0, 2, 1, 1, 2), -8, 1, Owner.CROSS),
        ((0, 2, 0, 0, 1, 2, 1, 1, 2), -8, 1, Owner.CROSS),
        ((0, 2, 0, 0, 0, 2, 1, 1, 2), -6, 0, Owner.NAUGHT),
        ((1, 2, 2, 0, 0, 0, 0, 0, 1), 9, 0, Owner.NAUGHT),
        ((2, 1, 1, 0, 0, 0, 0, 0, 2), -9, 0, Owner.CROSS),
        ((1, 2, 2, 0, 0, 2, 1, 1, 2), -10, 0, Owner.NAUGHT),
        ((1, 2, 2, 0, 0, 2, 1, 1, 2), -10, 0, Owner.CROSS),
        ((2, 0, 1, 0, 0, 0, 1, 0, 2), -9, 0, Owner.CROSS),
        ((2, 0, 1, 0, 2, 0, 0, 0, 2), -10, 0, Owner.NAUGHT),
        ((2, 1, 0, 1, 1, 0, 2, 2, 2), -10, 0, Owner.CROSS),
        ((2, 1, 0, 1, 1, 0, 0, 2, 2), -9, 0, Owner.CROSS),
    ],
)
def test_min_max_agent_algorithm(
    observation, correct_score, depth, current_player
) -> None:
    agent = MinMaxAgent(agent_mark=Owner.NAUGHT)
    assert (
        agent._min_max(
            observation=observation, depth=depth, current_player=current_player
        )
        == correct_score
    )


@pytest.mark.parametrize(
    ("observation", "correct_moves"),
    [
        ((1, 2, 1, 0, 2, 2, 1, 1, 2), [3]),
        ((0, 2, 1, 1, 2, 2, 1, 1, 2), [0]),
        ((0, 2, 0, 0, 1, 2, 1, 1, 2), [2]),
        ((2, 0, 1, 1, 2, 0, 0, 2, 0), [1, 8]),
        ((1, 0, 2, 0, 0, 0, 0, 0, 1), [4]),
        ((1, 2, 0, 0, 2, 0, 0, 1, 1), [6]),
    ],
)
def test_min_max_agent_correct_move(observation, correct_moves) -> None:
    agent = MinMaxAgent(agent_mark=Owner.CROSS)
    assert agent.take_min_max_action(observation=observation) in correct_moves
