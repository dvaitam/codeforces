package main

import (
	"bufio"
	"fmt"
	"os"
)

// This reference solution demonstrates a straightforward way to analyze the
// game.  Each card is treated as a vertex of a directed graph.  From a
// Monocarp card there are edges to every Bicarp card that can beat it, and vice
// versa.  States with no outgoing edges are losing for the player to move.
// Propagating this information via a queue allows determining the outcome for
// every starting move.
//
// This implementation builds the full adjacency matrix and therefore works only
// for small inputs, but it follows the classic algorithm for solving games on
// graphs with possible draws.

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		ax := make([]int, n)
		for i := range ax {
			fmt.Fscan(in, &ax[i])
		}
		ay := make([]int, n)
		for i := range ay {
			fmt.Fscan(in, &ay[i])
		}
		var m int
		fmt.Fscan(in, &m)
		bx := make([]int, m)
		for i := range bx {
			fmt.Fscan(in, &bx[i])
		}
		by := make([]int, m)
		for i := range by {
			fmt.Fscan(in, &by[i])
		}

		total := n + m
		edges := make([][]int, total)
		rev := make([][]int, total)
		outdeg := make([]int, total)

		// edges from Monocarp card i -> Bicarp card j
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				if bx[j] > ay[i] {
					edges[i] = append(edges[i], n+j)
					rev[n+j] = append(rev[n+j], i)
				}
			}
			outdeg[i] = len(edges[i])
		}
		for j := 0; j < m; j++ {
			for i := 0; i < n; i++ {
				if ax[i] > by[j] {
					edges[n+j] = append(edges[n+j], i)
					rev[i] = append(rev[i], n+j)
				}
			}
			outdeg[n+j] = len(edges[n+j])
		}

		state := make([]int, total) // 0 unknown,1 win,2 lose
		queue := make([]int, 0)
		for i := 0; i < total; i++ {
			if outdeg[i] == 0 {
				state[i] = 2 // lose
				queue = append(queue, i)
			}
		}
		for len(queue) > 0 {
			v := queue[0]
			queue = queue[1:]
			for _, u := range rev[v] {
				if state[u] != 0 {
					continue
				}
				if state[v] == 2 {
					state[u] = 1 // can move to losing position -> win
					queue = append(queue, u)
				} else {
					outdeg[u]--
					if outdeg[u] == 0 {
						state[u] = 2
						queue = append(queue, u)
					}
				}
			}
		}

		win, draw, lose := 0, 0, 0
		for i := 0; i < n; i++ {
			if state[i] == 1 {
				win++
			} else if state[i] == 2 {
				lose++
			} else {
				draw++
			}
		}
		fmt.Fprintln(out, win, draw, lose)
	}
}
