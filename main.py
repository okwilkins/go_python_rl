from python_xo.environment import NaughtsAndCrossesEnvironment, NaughtsAndCrossesEnvironmentGym, Owner
import random
from stable_baselines3 import PPO, DQN
from stable_baselines3.common.evaluation import evaluate_policy


def run_simulation(env: NaughtsAndCrossesEnvironment) -> None:
    env.reset()

    while not env.terminated():
        env.step(random.randint(0, 8))


def main() -> None:
    env = NaughtsAndCrossesEnvironment((
        (Owner.EMPTY, Owner.EMPTY, Owner.EMPTY),
        (Owner.EMPTY, Owner.EMPTY, Owner.EMPTY),
        (Owner.EMPTY, Owner.EMPTY, Owner.EMPTY),
    ))
    env.reset()

    for _ in range(1_000_000):
        run_simulation(env)

if __name__ == "__main__":
    # main()
    env = NaughtsAndCrossesEnvironment((
        (Owner.EMPTY, Owner.EMPTY, Owner.EMPTY),
        (Owner.EMPTY, Owner.EMPTY, Owner.EMPTY),
        (Owner.EMPTY, Owner.EMPTY, Owner.EMPTY),
    ))
    gym_env = NaughtsAndCrossesEnvironmentGym(env)
    model = PPO("MlpPolicy", gym_env, verbose=1)
    
    model = PPO.load("ppo_naughts_and_crosses", env=gym_env)
    # mean_reward, std_reward = evaluate_policy(
    #     model,
    #     gym_env,
    #     n_eval_episodes=1000,
    #     deterministic=True,
    # )
    # print(f'Mean reward: {mean_reward} +/- {std_reward}')


    obs, _ = gym_env.reset()
    terminated = False

    while not terminated:
        gym_env.render()
        action, _ = model.predict(obs, deterministic=True)
        obs, reward, terminated, _, _ = gym_env.step(action)
        print(f'Reward: {reward}')

    gym_env.render()

    # model.learn(total_timesteps=50_000)
    # model.save("ppo_naughts_and_crosses")