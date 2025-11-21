package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

type State struct {
	cost float64
	sum  float64
}

func prune(states []State) []State {
	sort.Slice(states, func(i, j int) bool {
		if states[i].cost == states[j].cost {
			return states[i].sum > states[j].sum
		}
		return states[i].cost < states[j].cost
	})
	res := make([]State, 0, len(states))
	bestSum := -1.0
	for _, st := range states {
		if st.sum <= bestSum+1e-12 {
			continue
		}
		bestSum = st.sum
		res = append(res, st)
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	groups := make([][]int, 101)
	for i := 0; i < n; i++ {
		var p, w int
		fmt.Fscan(in, &p, &w)
		groups[p] = append(groups[p], w)
	}

	baseSum := 0.0
	for _, w := range groups[100] {
		baseSum += float64(w)
	}
	groups[100] = nil

	states := []State{{0, baseSum}}

	for p := 1; p <= 99; p++ {
		ws := groups[p]
		if len(ws) == 0 {
			continue
		}
		sort.Slice(ws, func(i, j int) bool { return ws[i] > ws[j] })
		limit := (100 + (100 - p) - 1) / (100 - p)
		if len(ws) > limit {
			ws = ws[:limit]
		}
		costUnit := math.Log(100.0 / float64(p))
		options := make([]State, len(ws))
		prefix := 0.0
		for i, w := range ws {
			prefix += float64(w)
			options[i] = State{costUnit * float64(i+1), prefix}
		}
		newStates := make([]State, len(states)*(len(options)+1))
		copy(newStates, states)
		idx := len(states)
		for _, st := range states {
			for _, opt := range options {
				newStates[idx] = State{st.cost + opt.cost, st.sum + opt.sum}
				idx++
			}
		}
		states = prune(newStates[:idx])
	}

	ans := 0.0
	for _, st := range states {
		if st.sum == 0 {
			continue
		}
		val := math.Exp(-st.cost) * st.sum
		if val > ans {
			ans = val
		}
	}

	fmt.Printf("%.10f\n", ans)
}
