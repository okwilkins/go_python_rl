from abc import ABC, abstractmethod
import random
from typing import Literal

from python_xo.environment import Owner


class Agent(ABC):
    action_space = list[int]

    @abstractmethod
    def take_action(self, observation: tuple[int, ...]) -> int:
        ...


class RandomAgent(Agent):
    def __init__(self) -> None:
        self.action_space = list(range(9))
    
    def _get_index_of_empty_cells(self, observation: tuple[int, ...]) -> list[int]:
        return [i for i, cell in enumerate(observation) if cell == Owner.EMPTY.value]

    def take_action(self, observation: tuple[int, ...]) -> int:
        empty_indexes = self._get_index_of_empty_cells(observation=observation)
        return random.choice(empty_indexes)


class FillFirstEmptyAgent(Agent):
    def __init__(self) -> None:
        self.action_space = list(range(9))

    def take_action(self, observation: tuple[int, ...]) -> int:
        for i, cell in enumerate(observation):
            if cell == Owner.EMPTY:
                return i

        return random.choice(self.action_space)


class MinMaxAgent(Agent):
    def __init__(self, agent_mark: Literal[Owner.CROSS, Owner.NAUGHT]) -> None:
        self._AGENT_MARK = agent_mark
        self._FOE_MARK = Owner.NAUGHT if agent_mark == Owner.CROSS else Owner.CROSS
        self._SCORE_MAP = {
            self._AGENT_MARK.value: 10,
            self._FOE_MARK.value: -10,
            Owner.EMPTY.value: 0,
        }
        self.action_space = list(range(9))
    
    def take_action(self, observation: tuple[int, ...]) -> int:
        turn_counter = len(self._get_index_of_empty_cells(observation=observation)) - 9

        if turn_counter == 0:
            # Take random action to avoid heavy computation
            return self.take_random_action(observation=observation)
        else:
            return self.take_min_max_action(observation=observation)
    
    def take_random_action(self, observation: tuple[int, ...]) -> int:
        return random.choice(self._get_index_of_empty_cells(observation=observation))

    def take_min_max_action(self, observation: tuple[int, ...]) -> int:
        possible_games: list[tuple[tuple[int, ...], int, int]] = []
        starting_depth = 0

        for i in self._get_index_of_empty_cells(observation=observation):
            possible_game = list(observation)
            possible_game[i] = self._AGENT_MARK.value

            possible_games.append((
                tuple(possible_game),
                i,
                self._min_max(possible_game, starting_depth, self._AGENT_MARK)
            ))

        return max(possible_games, key=lambda x: x[2])[1]

    def _game_over(self, observation: tuple[int, ...]) -> bool:
        for i in range(3):
            # Check for a win in the rows
            if observation[i] == observation[i + 1] == observation[i + 2] != Owner.EMPTY.value:
                return True
            
            # Check for a win in the columns
            if observation[i] == observation[i + 3] == observation[i + 6] != Owner.EMPTY.value:
                return True

        # Check for a win in the diagonals
        if observation[0] == observation[4] == observation[8] != Owner.EMPTY.value:
            return True
        elif observation[2] == observation[4] == observation[6] != Owner.EMPTY.value:
            return True

        if Owner.EMPTY.value not in observation:
            return True

        return False

    def _score_board(self, observation: tuple[int, ...], depth: int) -> int:
        score = 0

        for i in range(3):
            # Check for a win in the rows
            if observation[i * 3] == observation[i * 3 + 1] == observation[i * 3 + 2] != Owner.EMPTY.value:
                score = self._SCORE_MAP[observation[i]]
                break
            
            # Check for a win in the columns
            if observation[i] == observation[i + 3] == observation[i + 6] != Owner.EMPTY.value:
                score = self._SCORE_MAP[observation[i]]
                break
        
        if score == 0:
            # Check for a win in the diagonals
            if observation[0] == observation[4] == observation[8] != Owner.EMPTY.value:
                score = self._SCORE_MAP[observation[0]]
            elif observation[2] == observation[4] == observation[6] != Owner.EMPTY.value:
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
    
    def _get_index_of_empty_cells(self, observation: tuple[int, ...]) -> list[int]:
        return [i for i, cell in enumerate(observation) if cell == Owner.EMPTY.value]

    def _min_max(self, observation: tuple[int, ...], depth: int, current_player = Literal[Owner.CROSS, Owner.NAUGHT]) -> int:
        # If game over, return the score
        if self._game_over(observation):
            return self._score_board(observation, depth)

        depth += 1
        scores = []
        possible_games = []

        for i in self._get_index_of_empty_cells(observation=observation):
            # Create a new game state
            possible_game = list(observation)
            possible_game[i] = current_player.value
            possible_games.append(tuple(possible_game))

        for possible_game in possible_games:
            if current_player == self._AGENT_MARK:
                next_player = self._FOE_MARK
            else:
                next_player = self._AGENT_MARK

            possible_score = self._min_max(possible_game, depth, next_player)
            scores.append(possible_score)

        if current_player == self._AGENT_MARK:
            return max(scores)
        else:
            return min(scores)
