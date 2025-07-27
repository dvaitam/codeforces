package main

import (
	"bufio"
	"fmt"
	"os"
)

// Solution for Codeforces problem 1594D - The Number of Imposters.
// We model each statement as a constraint between two players. If a player i says
// that player j is an "imposter", then their roles must differ; otherwise if the
// statement says "crewmate", their roles must be the same. We then attempt to
// assign each player a role of 0 (crewmate) or 1 (imposter) so that all
// constraints are satisfied. The graph constraints form a parity graph, so each
// connected component can be 2-colored. If a component is contradictory, the
// test case has no solution. Otherwise, for each component we choose the better
// orientation (either as-is or flipped) to maximize the number of imposters.

type edge struct {
	to  int
	typ int // 0 = same, 1 = different
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		g := make([][]edge, n+1)
		for i := 0; i < m; i++ {
			var a, b int
			var s string
			fmt.Fscan(in, &a, &b, &s)
			t := 0
			if s == "imposter" {
				t = 1
			}
			g[a] = append(g[a], edge{b, t})
			g[b] = append(g[b], edge{a, t})
		}

		color := make([]int, n+1)
		for i := range color {
			color[i] = -1
		}

		ans := 0
		ok := true
		for i := 1; i <= n && ok; i++ {
			if color[i] != -1 {
				continue
			}
			queue := []int{i}
			color[i] = 0
			cnt := [2]int{1, 0}
			for len(queue) > 0 && ok {
				v := queue[0]
				queue = queue[1:]
				for _, e := range g[v] {
					need := color[v] ^ e.typ
					if color[e.to] == -1 {
						color[e.to] = need
						cnt[need]++
						queue = append(queue, e.to)
					} else if color[e.to] != need {
						ok = false
						break
					}
				}
			}
			if !ok {
				break
			}
			if cnt[0] > cnt[1] {
				ans += cnt[0]
			} else {
				ans += cnt[1]
			}
		}
		if !ok {
			fmt.Fprintln(out, -1)
		} else {
			fmt.Fprintln(out, ans)
		}
	}
}
