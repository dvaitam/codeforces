package main

import (
	"bufio"
	"fmt"
	"os"
)

// This solution implements a straightforward probability DP that follows
// the tournament described in the problem statement. It directly simulates
// matches between sets of players. The approach is not optimised for the
// maximum constraints but demonstrates the intended algorithm.

type Node struct {
	probs map[int]float64
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	weights := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &weights[i])
	}

	// determine k so that n - k/2 is a power of two
	power := 1
	for power*2 < n {
		power *= 2
	}
	N := power
	k := 2 * (n - N)
	if k < 0 {
		k = 0
	}

	nodes := make([]Node, 0)
	// first round for first k participants
	for i := 1; i <= k; i += 2 {
		a := i
		b := i + 1
		pa := float64(weights[a]) / float64(weights[a]+weights[b])
		pb := 1 - pa
		m := make(map[int]float64)
		m[a] = pa
		m[b] = pb
		nodes = append(nodes, Node{m})
	}
	// remaining participants go directly
	for i := k + 1; i <= n; i++ {
		m := make(map[int]float64)
		m[i] = 1
		nodes = append(nodes, Node{m})
	}

	// simulate tournament
	for len(nodes) > 1 {
		next := make([]Node, 0)
		for i := 0; i < len(nodes); i += 2 {
			left := nodes[i]
			right := nodes[i+1]
			nd := make(map[int]float64)
			for p1, prob1 := range left.probs {
				for p2, prob2 := range right.probs {
					w1 := float64(weights[p1])
					w2 := float64(weights[p2])
					win1 := prob1 * prob2 * w1 / (w1 + w2)
					win2 := prob1 * prob2 * w2 / (w1 + w2)
					nd[p1] += win1
					nd[p2] += win2
				}
			}
			next = append(next, Node{nd})
		}
		nodes = next
	}

	finalDist := nodes[0].probs
	for i := 1; i <= n; i++ {
		if p, ok := finalDist[i]; ok {
			fmt.Fprintf(out, "%.10f ", p)
		} else {
			fmt.Fprintf(out, "0 ")
		}
	}
}
