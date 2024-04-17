package main

import (
	"fmt"
	xo "xo/src/go_xo"
)

func run_simulation(env *xo.NaughtsAndCrossesEnvironment) {
	env.Reset()
	opponent_agent := xo.MinMaxAgent{
		AgentMark:    *env.Agent.GetMark(),
		OpponentMark: env.UserMark,
	}

	for !env.Terminated() {
		observation := env.Observation()
		env.Step(opponent_agent.TakeAction(&observation))
	}
}

func main() {
	agent := &xo.MinMaxAgent{AgentMark: xo.Naught, OpponentMark: xo.Cross}
	env := xo.NaughtsAndCrossesEnvironment{
		Board: xo.Board{
			{xo.Empty, xo.Empty, xo.Empty},
			{xo.Empty, xo.Empty, xo.Empty},
			{xo.Empty, xo.Empty, xo.Empty},
		},
		TimeStep: 0,
		UserMark: xo.Cross,
		Agent:    agent,
	}
	for i := 0; i < 9; i++ {
		fmt.Printf("Simulation %d\n", i+1)
		run_simulation(&env)
		env.Render()
	}
}
