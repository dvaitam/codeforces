package main

import (
	"bufio"
	"fmt"
	"os"
)

type Parent struct {
	prev   int
	player byte
}

type Step struct {
	player byte
	card   int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	players := []byte{'q', 'k', 'j'}

	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)

		perms := make([][]int, 3)
		for p := 0; p < 3; p++ {
			perms[p] = make([]int, n+1)
			for i := 1; i <= n; i++ {
				fmt.Fscan(in, &perms[p][i])
			}
		}

		reachable := make([]bool, n+1)
		parent := make([]Parent, n+1)
		bestVal := []int{-1, -1, -1}
		bestCard := []int{0, 0, 0}

		reachable[1] = true
		for p := 0; p < 3; p++ {
			bestVal[p] = perms[p][1]
			bestCard[p] = 1
		}

		for i := 2; i <= n; i++ {
			for p := 0; p < 3; p++ {
				if bestVal[p] > perms[p][i] {
					reachable[i] = true
					parent[i] = Parent{prev: bestCard[p], player: players[p]}
					break
				}
			}
			if reachable[i] {
				for p := 0; p < 3; p++ {
					if perms[p][i] > bestVal[p] {
						bestVal[p] = perms[p][i]
						bestCard[p] = i
					}
				}
			}
		}

		if !reachable[n] {
			fmt.Fprintln(out, "NO")
			continue
		}

		steps := make([]Step, 0)
		cur := n
		for cur > 1 {
			par := parent[cur]
			steps = append(steps, Step{player: par.player, card: cur})
			cur = par.prev
		}

		fmt.Fprintln(out, "YES")
		fmt.Fprintln(out, len(steps))
		for i := len(steps) - 1; i >= 0; i-- {
			fmt.Fprintf(out, "%c %d\n", steps[i].player, steps[i].card)
		}
	}
}
