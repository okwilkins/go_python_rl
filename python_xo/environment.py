import random
from typing import Literal

import gymnasium as gym

from python_xo.agent import Agent
from python_xo.owner import Owner


class NaughtsAndCrossesEnvironment:
    def __init__(
        self,
        board: list[list[Owner]],
        user_mark: Literal[Owner.CROSS, Owner.NAUGHT],
        agent: Agent,
        time_step: int = 0,
    ) -> None:
        self.board = board
        self.time_step = time_step
        self.user_mark = user_mark
        self.agent = agent
        self.last_player = agent.agent_mark

    def reset(self) -> None:
        self.board = [
            [Owner.EMPTY, Owner.EMPTY, Owner.EMPTY],
            [Owner.EMPTY, Owner.EMPTY, Owner.EMPTY],
            [Owner.EMPTY, Owner.EMPTY, Owner.EMPTY],
        ]
        self.time_step = 0

        # 50/50 change of the agent taking the first move
        if random.choice([True, False]):
            self.agent_take_turn(self.agent.take_action(self.observation()))

    def agent_take_turn(self, action: int):
        row = action // 3
        col = action % 3
        if self.board[row][col] == Owner.EMPTY:
            self.board[row][col] = self.agent.agent_mark

    def render(self) -> None:
        for row in self.board:
            for i, cell in enumerate(row):
                if cell == Owner.EMPTY:
                    print(" ", end="")
                elif cell == Owner.NAUGHT:
                    print("O", end="")
                elif cell == Owner.CROSS:
                    print("X", end="")

                if i < 2:
                    print("|", end="")

            print()

    def observation(self) -> tuple[int, int, int, int, int, int, int, int, int]:
        return (
            self.board[0][0].value,
            self.board[0][1].value,
            self.board[0][2].value,
            self.board[1][0].value,
            self.board[1][1].value,
            self.board[1][2].value,
            self.board[2][0].value,
            self.board[2][1].value,
            self.board[2][2].value,
        )

    def game_is_draw(self) -> bool:
        return all([cell != Owner.EMPTY for row in self.board for cell in row])

    def game_winner(self) -> Owner:
        for i in range(3):
            # Check for a win in the rows
            if (
                self.board[i][0] == self.board[i][1] == self.board[i][2]
                and self.board[i][0] != Owner.EMPTY
            ):
                return self.board[i][0]

            # Check for a win in the columns
            if (
                self.board[0][i] == self.board[1][i] == self.board[2][i]
                and self.board[0][i] != Owner.EMPTY
            ):
                return self.board[0][i]

        # Check for a win in the diagonals
        if (
            self.board[0][0] == self.board[1][1] == self.board[2][2]
            and self.board[0][0] != Owner.EMPTY
        ):
            return self.board[0][0]

        if (
            self.board[0][2] == self.board[1][1] == self.board[2][0]
            and self.board[0][2] != Owner.EMPTY
        ):
            return self.board[0][2]

        return Owner.EMPTY

    def terminated(self) -> bool:
        return self.game_is_draw() or self.game_winner() != Owner.EMPTY

    def truncated(self) -> bool:
        return self.time_step > 9

    def reward(self) -> int:
        match (self.game_winner(), self.game_is_draw()):
            case (Owner.NAUGHT, _):
                return -100
            case (Owner.CROSS, _):
                return 100
            case (_, True):
                return 50
            case _:
                return 0

    def step(
        self, action: int
    ) -> tuple[tuple[int, int, int, int, int, int, int, int, int], int, bool, bool]:
        if not self.terminated():
            row = action // 3
            col = action % 3
            if self.board[row][col] == Owner.EMPTY:
                self.board[row][col] = self.user_mark

            if not self.terminated():
                self.agent_take_turn(self.agent.take_action(self.observation()))

        self.time_step += 1

        return self.observation(), self.reward(), self.terminated(), self.truncated()


class NaughtsAndCrossesEnvironmentGym(gym.Env):
    def __init__(self, naughts_and_crosses_environment: NaughtsAndCrossesEnvironment) -> None:
        self.naughts_and_crosses_environment = naughts_and_crosses_environment
        self.action_space = gym.spaces.Discrete(9)
        self.observation_space = gym.spaces.MultiDiscrete([3] * 9)

    def reset(self, seed=None, options=None):
        super().reset(seed=seed)

        self.naughts_and_crosses_environment.reset()
        return self.naughts_and_crosses_environment.observation(), None

    def step(self, action: int):
        observation, reward, terminated, truncated = self.naughts_and_crosses_environment.step(
            action
        )
        return observation, reward, terminated, truncated, {}

    def render(self):
        self.naughts_and_crosses_environment.render()

    def close(self):
        pass
