# import random
from stable_baselines3 import PPO

from python_xo.agent import MinMaxAgent
from python_xo.environment import (
    NaughtsAndCrossesEnvironment,
    NaughtsAndCrossesEnvironmentGym,
    Owner,
)

# from stable_baselines3.common.evaluation import evaluate_policy

# import argparse


def run_simulation() -> None:
    agent = MinMaxAgent(agent_mark=Owner.NAUGHT)
    env = NaughtsAndCrossesEnvironment(
        board=[
            [Owner.EMPTY, Owner.EMPTY, Owner.EMPTY],
            [Owner.EMPTY, Owner.EMPTY, Owner.EMPTY],
            [Owner.EMPTY, Owner.EMPTY, Owner.EMPTY],
        ],
        time_step=0,
        user_mark=Owner.CROSS,
        agent=agent,
    )

    while not env.terminated():
        env.render()

        _input = int(input())

        if not env.terminated():
            env.step(_input)

    env.render()


def main() -> None:
    agent = MinMaxAgent(agent_mark=Owner.NAUGHT)
    env = NaughtsAndCrossesEnvironment(
        board=[
            [Owner.EMPTY, Owner.EMPTY, Owner.EMPTY],
            [Owner.EMPTY, Owner.EMPTY, Owner.EMPTY],
            [Owner.EMPTY, Owner.EMPTY, Owner.EMPTY],
        ],
        time_step=0,
        user_mark=Owner.CROSS,
        agent=agent,
    )
    gym_env = NaughtsAndCrossesEnvironmentGym(env)
    gym_env.reset()

    model = PPO("MlpPolicy", gym_env, verbose=2)
    model.learn(10_000, progress_bar=True)


if __name__ == "__main__":
    main()
