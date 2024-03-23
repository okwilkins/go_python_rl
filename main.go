package main

import (
	"math/rand"
	xo "xo/go_xo"
)

func run_simulation(env *xo.NaughtsAndCrossesEnvironment) {
    env.Reset()

    for !env.Terminated() {
        env.Step(rand.Intn(9))
    }
}

func run_simulation_non_ref() {
    env := xo.NaughtsAndCrossesEnvironment{}
    env.Reset()

    for !env.Terminated() {
        env.Step(rand.Intn(9))
    }
}


func main() {
	env := xo.NaughtsAndCrossesEnvironment{}
	env.Reset()

    for i := 0; i < 1_000_0000; i++ {
        go run_simulation_non_ref()
    }
}
