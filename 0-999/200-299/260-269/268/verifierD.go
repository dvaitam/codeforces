package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const MOD = 1000000009

func solveK(n, h, k int) int64 {
	type State struct {
		x [4]int
	}
	states := []State{}
	stateToID := make(map[State]int)

	var generate func(idx int, maxVal int, current State)
	generate = func(idx int, maxVal int, current State) {
		if idx == k {
			stateToID[current] = len(states)
			states = append(states, current)
			return
		}
		for v := maxVal; v >= 0; v-- {
			current.x[idx] = v
			generate(idx+1, v, current)
		}
	}

	var startState State
	generate(0, h-1, startState)

	numStates := len(states)
	outTrans := make([]int, numStates)
	inTrans := make([][]int, numStates)
	for i := range inTrans {
		inTrans[i] = make([]int, k)
	}

	for id, s := range states {
		outS := State{}
		valid := true
		for i := 0; i < k; i++ {
			outS.x[i] = s.x[i] + 1
			if outS.x[i] >= h {
				valid = false
			}
		}
		if valid {
			outTrans[id] = stateToID[outS]
		} else {
			outTrans[id] = -1
		}

		for m := 0; m < k; m++ {
			inS := State{}
			validIn := true
			var newGaps [4]int
			for i := 0; i < k; i++ {
				if i == m {
					newGaps[i] = 0
				} else {
					g := s.x[i] + 1
					if g >= h {
						validIn = false
					}
					newGaps[i] = g
				}
			}
			if validIn {
				for i := 0; i < k; i++ {
					for j := i + 1; j < k; j++ {
						if newGaps[i] < newGaps[j] {
							newGaps[i], newGaps[j] = newGaps[j], newGaps[i]
						}
					}
				}
				for i := 0; i < k; i++ {
					inS.x[i] = newGaps[i]
				}
				inTrans[id][m] = stateToID[inS]
			} else {
				inTrans[id][m] = -1
			}
		}
	}

	type Trans struct {
		dest int
		mult int64
	}
	transitions := make([][]Trans, numStates)
	for id := 0; id < numStates; id++ {
		var dests []int
		var mults []int64
		add := func(dest int, count int64) {
			for i, d := range dests {
				if d == dest {
					mults[i] += count
					return
				}
			}
			dests = append(dests, dest)
			mults = append(mults, count)
		}
		if outTrans[id] != -1 && 4-k > 0 {
			add(outTrans[id], int64(4-k))
		}
		for m := 0; m < k; m++ {
			if inTrans[id][m] != -1 {
				add(inTrans[id][m], 1)
			}
		}
		for i, d := range dests {
			transitions[id] = append(transitions[id], Trans{d, mults[i]})
		}
	}

	dp := make([]int64, numStates)
	var initS State
	dp[stateToID[initS]] = 1

	for step := 0; step < n; step++ {
		nextDp := make([]int64, numStates)
		for id, ways := range dp {
			if ways == 0 {
				continue
			}
			for _, t := range transitions[id] {
				nextDp[t.dest] = (nextDp[t.dest] + ways*t.mult) % MOD
			}
		}
		dp = nextDp
	}

	var total int64 = 0
	for _, ways := range dp {
		total = (total + ways) % MOD
	}
	return total
}

func solveCase(n, h int) int {
	n1 := solveK(n, h, 1)
	n2 := solveK(n, h, 2)
	n3 := solveK(n, h, 3)
	n4 := solveK(n, h, 4)

	ans := (4*n1 - 6*n2 + 4*n3 - n4) % MOD
	if ans < 0 {
		ans += MOD
	}
	return int(ans)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	h := rng.Intn(minInt(n, 5)) + 1
	input := fmt.Sprintf("%d %d\n", n, h)
	expected := fmt.Sprintf("%d", solveCase(n, h))
	return input, expected
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func runCase(bin, input, expected string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
