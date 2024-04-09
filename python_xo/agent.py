from abc import ABC, abstractmethod
import random
from typing import Literal

from python_xo.owner import Owner


OBS_TYPE = tuple[int, int, int, int, int, int, int, int, int]


class Agent(ABC):
    action_space: list[int]
    agent_mark: Literal[Owner.CROSS, Owner.NAUGHT]

    @abstractmethod
    def take_action(self, observation: OBS_TYPE) -> int: ...


class RandomAgent(Agent):
    def __init__(self, agent_mark: Literal[Owner.CROSS, Owner.NAUGHT]) -> None:
        self.agent_mark = agent_mark
        self.action_space = list(range(9))

    def _get_index_of_empty_cells(self, observation: OBS_TYPE) -> list[int]:
        return [i for i, cell in enumerate(observation) if cell == Owner.EMPTY.value]

    def take_action(self, observation: OBS_TYPE) -> int:
        empty_indexes = self._get_index_of_empty_cells(observation=observation)
        return random.choice(empty_indexes)


class FillFirstEmptyAgent(Agent):
    def __init__(self, agent_mark: Literal[Owner.CROSS, Owner.NAUGHT]) -> None:
        self.agent_mark = agent_mark
        self.action_space = list(range(9))

    def take_action(self, observation: OBS_TYPE) -> int:
        for i, cell in enumerate(observation):
            if cell == Owner.EMPTY:
                return i

        return random.choice(self.action_space)


class MinMaxAgent(Agent):
    def __init__(self, agent_mark: Literal[Owner.CROSS, Owner.NAUGHT]) -> None:
        self.agent_mark = agent_mark
        self._FOE_MARK = Owner.NAUGHT if agent_mark == Owner.CROSS else Owner.CROSS
        self._SCORE_MAP = {
            self.agent_mark.value: 10,
            self._FOE_MARK.value: -10,
            Owner.EMPTY.value: 0,
        }
        self.action_space = list(range(9))

    def take_action(self, observation: OBS_TYPE) -> int:
        turn_counter = len(self._get_index_of_empty_cells(observation=observation)) - 9

        if turn_counter == 0:
            # Take random action to avoid heavy computation
            return self.take_random_action(observation=observation)
        else:
            return self.take_min_max_action(observation=observation)

    def take_random_action(self, observation: OBS_TYPE) -> int:
        return random.choice(self._get_index_of_empty_cells(observation=observation))

    def take_min_max_action(self, observation: OBS_TYPE) -> int:
        possible_games_and_scores: list[tuple[OBS_TYPE, int, int]] = []

        for cell_idx in self._get_index_of_empty_cells(observation=observation):
            possible_game: list[int] = list(observation)
            possible_game[cell_idx] = self.agent_mark.value
            possible_game: OBS_TYPE = tuple(possible_game)  # type: ignore

            next_player = self._get_next_player(current_player=self.agent_mark)
            score = self._min_max(possible_game, 0, next_player)  # type: ignore

            possible_games_and_scores.append((next_player, cell_idx, score))  # type: ignore

        max_score = max(possible_games_and_scores, key=lambda x: x[2])[2]
        # Randomly choose from the best moves
        best_moves = [
            move for move in possible_games_and_scores if move[2] == max_score
        ]
        return random.choice(best_moves)[1]

    def _game_over(self, observation: OBS_TYPE) -> bool:
        for i in range(3):
            # Check for a win in the rows
            if (
                observation[i * 3]
                == observation[i * 3 + 1]
                == observation[i * 3 + 2]
                != Owner.EMPTY.value
            ):
                return True

            # Check for a win in the columns
            if (
                observation[i]
                == observation[i + 3]
                == observation[i + 6]
                != Owner.EMPTY.value
            ):
                return True

        # Check for a win in the diagonals
        if observation[0] == observation[4] == observation[8] != Owner.EMPTY.value:
            return True
        elif observation[2] == observation[4] == observation[6] != Owner.EMPTY.value:
            return True

        if Owner.EMPTY.value not in observation:
            return True

        return False

    def _score_board(self, observation: OBS_TYPE, depth: int) -> int:
        score = 0

        for i in range(3):
            # Check for a win in the rows
            if (
                observation[i * 3]
                == observation[i * 3 + 1]
                == observation[i * 3 + 2]
                != Owner.EMPTY.value
            ):
                score = self._SCORE_MAP[observation[i * 3]]
                break

            # Check for a win in the columns
            if (
                observation[i]
                == observation[i + 3]
                == observation[i + 6]
                != Owner.EMPTY.value
            ):
                score = self._SCORE_MAP[observation[i]]
                break

        if score == 0:
            # Check for a win in the diagonals
            if observation[0] == observation[4] == observation[8] != Owner.EMPTY.value:
                score = self._SCORE_MAP[observation[0]]
            elif (
                observation[2] == observation[4] == observation[6] != Owner.EMPTY.value
            ):
                score = self._SCORE_MAP[observation[2]]

        # The player won
        if score == 10:
            return score - depth
        # The player lost
        elif score == -10:
            return score + depth
        # The game is a draw
        else:
            return 0

    def _get_index_of_empty_cells(self, observation: OBS_TYPE) -> list[int]:
        return [i for i, cell in enumerate(observation) if cell == Owner.EMPTY.value]

    def _get_possible_games(
        self,
        observation: OBS_TYPE,
        next_player: Literal[Owner.CROSS, Owner.NAUGHT],
    ) -> list[OBS_TYPE]:
        possible_games = []

        for i in self._get_index_of_empty_cells(observation=observation):
            # Create a new game state
            possible_game = list(observation)
            possible_game[i] = next_player.value
            possible_game: OBS_TYPE = tuple(possible_game)  # type: ignore
            possible_games.append(possible_game)

        return possible_games  # type: ignore

    def _get_next_player(
        self,
        current_player: Literal[Owner.CROSS, Owner.NAUGHT],
    ) -> Literal[Owner.CROSS, Owner.NAUGHT]:
        if current_player == self.agent_mark:
            next_player = self._FOE_MARK
        else:
            next_player = self.agent_mark

        return next_player  # type: ignore

    def _min_max(
        self,
        observation: OBS_TYPE,
        depth: int,
        current_player=Literal[Owner.CROSS, Owner.NAUGHT],
    ) -> int:
        # If game over, return the score
        if self._game_over(observation):
            return self._score_board(observation, depth)

        depth += 1
        scores = []
        possible_games = self._get_possible_games(observation, current_player)

        for possible_game in possible_games:
            next_player = self._get_next_player(current_player)
            possible_score = self._min_max(possible_game, depth, next_player)
            scores.append(possible_score)

        if current_player == self.agent_mark:
            return max(scores)
        else:
            return min(scores)
