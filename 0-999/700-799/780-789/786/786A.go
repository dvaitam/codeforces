package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	k := make([]int, 2)
	moves := make([][]int, 2)
	for p := 0; p < 2; p++ {
		fmt.Fscan(in, &k[p])
		moves[p] = make([]int, k[p])
		for i := 0; i < k[p]; i++ {
			fmt.Fscan(in, &moves[p][i])
		}
	}

	const (
		UNKNOWN = 0
		WIN     = 1
		LOSE    = -1
	)

	state := make([][]int, 2)
	deg := make([][]int, 2)
	for p := 0; p < 2; p++ {
		state[p] = make([]int, n)
		deg[p] = make([]int, n)
		for i := 0; i < n; i++ {
			deg[p][i] = k[p]
		}
	}

	queue := make([][2]int, 0)
	for p := 0; p < 2; p++ {
		state[p][0] = LOSE
		queue = append(queue, [2]int{p, 0})
	}

	for head := 0; head < len(queue); head++ {
		p := queue[head][0]
		pos := queue[head][1]
		cur := state[p][pos]
		opp := 1 - p
		for _, mv := range moves[opp] {
			prev := pos - mv
			prev = ((prev % n) + n) % n
			if state[opp][prev] != UNKNOWN {
				continue
			}
			if cur == LOSE {
				state[opp][prev] = WIN
				queue = append(queue, [2]int{opp, prev})
			} else if cur == WIN {
				deg[opp][prev]--
				if deg[opp][prev] == 0 {
					state[opp][prev] = LOSE
					queue = append(queue, [2]int{opp, prev})
				}
			}
		}
	}

	for p := 0; p < 2; p++ {
		for i := 1; i < n; i++ {
			if i > 1 {
				fmt.Fprint(out, " ")
			}
			switch state[p][i] {
			case WIN:
				fmt.Fprint(out, "Win")
			case LOSE:
				fmt.Fprint(out, "Lose")
			default:
				fmt.Fprint(out, "Loop")
			}
		}
		fmt.Fprintln(out)
	}
}
