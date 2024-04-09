from python_xo.environment import (
    NaughtsAndCrossesEnvironment,
    Owner,
)
from python_xo.agent import MinMaxAgent

# import random
# from stable_baselines3 import PPO, DQN
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
    for _ in range(1):
        run_simulation()


if __name__ == "__main__":
    main()
    # env = NaughtsAndCrossesEnvironment(
    #     board=[
    #         [Owner.EMPTY, Owner.EMPTY, Owner.EMPTY],
    #         [Owner.EMPTY, Owner.EMPTY, Owner.EMPTY],
    #         [Owner.EMPTY, Owner.EMPTY, Owner.EMPTY],
    #     ],
    #     agent=ag
    # )
    # gym_env = NaughtsAndCrossesEnvironmentGym(env)
    # model = PPO("MlpPolicy", gym_env, verbose=1)

    # model = PPO.load("ppo_naughts_and_crosses", env=gym_env)
    # mean_reward, std_reward = evaluate_policy(
    #     model,
    #     gym_env,
    #     n_eval_episodes=1000,
    #     deterministic=True,
    # )
    # print(f'Mean reward: {mean_reward} +/- {std_reward}')

    # obs = gym_env.naughts_and_crosses_environment.observation()
    # terminated = False

    # while not terminated:
    #     gym_env.render()
    #     t = input("Move: ")
    #     row = int(t) // 3
    #     col = int(t) % 3
    #     gym_env.naughts_and_crosses_environment.board[row][col] = Owner.NAUGHT

    #     action, _ = model.predict(obs, deterministic=True)
    #     obs, reward, terminated, _, _ = gym_env.step(action)
    #     print(f'Reward: {reward}')

    # gym_env.render()

    # model.learn(total_timesteps=1_000_000)
    # model.save("ppo_naughts_and_crosses")
