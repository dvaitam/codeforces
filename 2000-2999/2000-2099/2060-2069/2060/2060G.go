package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

// Global scanner for fast input
var sc = bufio.NewScanner(os.Stdin)

func init() {
	sc.Split(bufio.ScanWords)
	// Set a sufficiently large buffer for tokens
	buf := make([]byte, 0, 1<<18)
	sc.Buffer(buf, 1<<22)
}

func nextInt() int {
	sc.Scan()
	i, _ := strconv.Atoi(sc.Text())
	return i
}

func main() {
	// Initialize scanner by reading the first token (t)
	if !sc.Scan() {
		return
	}
	t, _ := strconv.Atoi(sc.Text())

	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	for i := 0; i < t; i++ {
		solve(writer)
	}
}

type Pair struct {
	u, v int
}

type State struct {
	a, b int
}

func solve(w *bufio.Writer) {
	n := nextInt()
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = nextInt()
	}
	b := make([]int, n)
	for i := 0; i < n; i++ {
		b[i] = nextInt()
	}

	pairs := make([]Pair, n)
	inputParity := 0
	for i := 0; i < n; i++ {
		u, v := a[i], b[i]
		if u > v {
			u, v = v, u
			inputParity ^= 1
		}
		pairs[i] = Pair{u, v}
	}

	// Sort pairs based on their maximum value (v)
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].v < pairs[j].v
	})

	// dp[p] stores a list of Pareto optimal states (last_a, last_b) for accumulated parity p
	// p is 0 or 1
	dp := make([][]State, 2)
	dp[0] = []State{{0, 0}}
	dp[1] = []State{}

	for i := 0; i < n; i++ {
		u, v := pairs[i].u, pairs[i].v
		
		nextDp0 := make([]State, 0)
		nextDp1 := make([]State, 0)

		// Transition from dp[0]
		for _, s := range dp[0] {
			// Orientation 0: (u, v) - Adds 0 to parity
			if u > s.a && v > s.b {
				nextDp0 = append(nextDp0, State{u, v})
			}
			// Orientation 1: (v, u) - Adds 1 to parity
			if v > s.a && u > s.b {
				nextDp1 = append(nextDp1, State{v, u})
			}
		}

		// Transition from dp[1]
		for _, s := range dp[1] {
			// Orientation 0: (u, v) - Adds 0 to parity (remains 1)
			if u > s.a && v > s.b {
				nextDp1 = append(nextDp1, State{u, v})
			}
			// Orientation 1: (v, u) - Adds 1 to parity (becomes 0)
			if v > s.a && u > s.b {
				nextDp0 = append(nextDp0, State{v, u})
			}
		}

		dp[0] = prune(nextDp0)
		dp[1] = prune(nextDp1)
	}

	if len(dp[inputParity]) > 0 {
		fmt.Fprintln(w, "YES")
	} else {
		fmt.Fprintln(w, "NO")
	}
}

// prune removes dominated states. A state (a1, b1) dominates (a2, b2) if a1 <= a2 and b1 <= b2.
// We only keep states that are not dominated by any other state.
func prune(states []State) []State {
	if len(states) <= 1 {
		return states
	}
	// Sort by 'a' ascending. If 'a' is same, sort by 'b' ascending.
	sort.Slice(states, func(i, j int) bool {
		if states[i].a != states[j].a {
			return states[i].a < states[j].a
		}
		return states[i].b < states[j].b
	})

	res := states[:0]
	minB := 2000000000 

	// Iterate through sorted states. Since 'a' is increasing, we only keep a state
	// if its 'b' is strictly smaller than all previous 'b's.
	for _, s := range states {
		if s.b < minB {
			res = append(res, s)
			minB = s.b
		}
	}
	return res
}