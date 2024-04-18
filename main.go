package main

import (
	"fmt"
	"os"
	"sync"
	"time"
	xo "xo/src/go_xo"
)

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}

func run_simulation(env xo.NaughtsAndCrossesEnvironment) {
	env.Reset()
	opponent_agent := xo.MinMaxAgent{
		AgentMark:    env.Agent.GetMark(),
		OpponentMark: env.UserMark,
	}

	for !env.Terminated() {
		observation := env.Observation()
		action, err := opponent_agent.TakeAction(&observation)

		if err == nil {
			env.Step(action)
		} else {
			fmt.Printf("%v", err)
		}
	}
}

func main() {
	defer timer("main")()
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
	wg := sync.WaitGroup{}

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			run_simulation(env)
		}(i)
	}

	wg.Wait()
	fmt.Println("Done")
	os.Exit(0)
}
