from python_xo.environment import NaughtsAndCrossesEnvironment, NaughtsAndCrossesEnvironmentGym, Owner
from python_xo.agent import MinMaxAgent, RandomAgent, FillFirstEmptyAgent
import random
from stable_baselines3 import PPO, DQN
from stable_baselines3.common.evaluation import evaluate_policy


def run_simulation() -> None:
    env = NaughtsAndCrossesEnvironment(
        board=[
            [Owner.EMPTY, Owner.EMPTY, Owner.EMPTY],
            [Owner.EMPTY, Owner.EMPTY, Owner.EMPTY],
            [Owner.EMPTY, Owner.EMPTY, Owner.EMPTY],
        ],
        time_step=0
    )
    agent = MinMaxAgent(agent_mark=Owner.NAUGHT)
    # rand_agent = RandomAgent()

    # 50/50 change of the agent taking the first move
    if random.choice([True, False]):
        env.step(agent.take_action(env.observation()))

    while not env.terminated():
        # env.step(rand_agent.take_action(env.observation()))

        env.render()

        _input = int(input()) - 1
        row = _input // 3
        col = _input % 3

        env.board[row][col] = Owner.CROSS

        if not env.terminated():
            env.step(agent.take_action(env.observation()))

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