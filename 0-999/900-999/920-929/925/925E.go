package main

import (
	"bufio"
	"fmt"
	"os"
)

type Node struct {
	parent     int
	t          int
	sub        int
	on         bool
	displeased bool
}

var (
	n, m  int
	nodes []Node
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fscan(reader, &n, &m)
	nodes = make([]Node, n+1)
	for i := 2; i <= n; i++ {
		fmt.Fscan(reader, &nodes[i].parent)
	}
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &nodes[i].t)
	}
	events := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &events[i])
	}

	displeasedCount := 0
	// initial all on=false, sub=0, displeased=false

	writer := bufio.NewWriter(os.Stdout)
	for _, q := range events {
		x := q
		if x > 0 {
			// employee x leaves
			if nodes[x].displeased {
				nodes[x].displeased = false
				displeasedCount--
			}
			nodes[x].on = true
			// update ancestors
			v := nodes[x].parent
			for v != 0 {
				nodes[v].sub++
				if !nodes[v].on {
					if !nodes[v].displeased && nodes[v].sub > nodes[v].t {
						nodes[v].displeased = true
						displeasedCount++
					}
				}
				v = nodes[v].parent
			}
		} else {
			x = -x
			// employee x returns
			nodes[x].on = false
			if nodes[x].sub > nodes[x].t {
				if !nodes[x].displeased {
					nodes[x].displeased = true
					displeasedCount++
				}
			} else {
				if nodes[x].displeased {
					nodes[x].displeased = false
					displeasedCount--
				}
			}
			v := nodes[x].parent
			for v != 0 {
				nodes[v].sub--
				if !nodes[v].on {
					if nodes[v].displeased && nodes[v].sub <= nodes[v].t {
						nodes[v].displeased = false
						displeasedCount--
					} else if !nodes[v].displeased && nodes[v].sub > nodes[v].t {
						nodes[v].displeased = true
						displeasedCount++
					}
				}
				v = nodes[v].parent
			}
		}
		fmt.Fprintf(writer, "%d ", displeasedCount)
	}
	writer.WriteByte('\n')
	writer.Flush()
}
