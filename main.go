package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
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
	env.Reset()

	reader := bufio.NewReader(os.Stdin)

	for !env.Terminated() {
		env.Render()
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		if !env.Terminated() {
			num, err := strconv.Atoi(text)
			if err != nil || num < 0 || num > 8 {
				fmt.Println("Invalid input. Please enter a valid integer.")
				continue
			}
			env.Step(num)
		}
	}

	env.Render()
}
