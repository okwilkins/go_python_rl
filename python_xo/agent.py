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
        self._agent_mark = agent_mark
        self._foe_mark = Owner.NAUGHT if agent_mark == Owner.CROSS else Owner.CROSS
        self._score_map = {
            self._agent_mark.value: 10,
            self._foe_mark.value: -10,
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
            possible_game[i] = self._agent_mark.value

            possible_games.append((
                tuple(possible_game),
                i,
                self._min_max(possible_game, starting_depth, self._agent_mark)
            ))

        return max(possible_games, key=lambda x: x[2])[1]

    def _game_over(self, observation: tuple[int, ...]) -> bool:
        return all([cell != Owner.EMPTY.value for cell in observation])

    def _score_board(self, observation: tuple[int, ...], depth: int) -> int:
        score = 0

        for i in range(3):
            # Check for a win in the rows
            if observation[i] == observation[i + 3] == observation[i + 6] != Owner.EMPTY:
                score = self._score_map[observation[i]]
                break
            
            # Check for a win in the columns
            if observation[i] == observation[i + 3] == observation[i + 6]:
                score = self._score_map[observation[i]]
                break
        
        if score != 0:
            # Check for a win in the diagonals
            if observation[0] == observation[4] == observation[8] != Owner.EMPTY:
                score = self._score_map[observation[0]]
            elif observation[2] == observation[4] == observation[6] != Owner.EMPTY:
                score = self._score_map[observation[2]]
        
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

    def _min_max(self, observation: tuple[int, ...], depth: int, last_move = Owner) -> int:
        # If game over, return the score
        if self._game_over(observation):
            return self._score_board(observation, depth)

        depth += 1
        scores = []

        for i in self._get_index_of_empty_cells(observation=observation):
            possible_game = list(observation)
            
            match last_move:
                case self._agent_mark:
                    last_move = self._foe_mark
                    possible_game[i] = self._foe_mark.value
                case self._foe_mark:
                    last_move = self._agent_mark
                    possible_game[i] = self._agent_mark.value

            possible_score = self._min_max(possible_game, depth, last_move)
            scores.append(possible_score)

        return max(scores)
