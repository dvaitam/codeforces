package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

type Edge struct {
	u, v int
}

const (
	maxN = 32768
	maxM = 65536
)

func gadgetValue(l int) int64 {
	// Value contributed by a chain gadget of length l (excluding global node 1).
	// Derived from simulation: l^2 + 3l - 2 (includes the edge to node 1).
	return int64(l*l + 3*l - 2)
}

func findCoins(rem int64, edgesAvail, nodesAvail int) (c11, c8, c2 int, ok bool) {
	if rem == 0 {
		return 0, 0, 0, true
	}
	max11 := int(min(rem/11, int64(edgesAvail/4)))
	if nodesAvail/3 < max11 {
		max11 = nodesAvail / 3
	}
	for c11 = max11; c11 >= 0; c11-- {
		rem1 := rem - int64(11*c11)
		if rem1 < 0 {
			continue
		}
		if rem1%2 == 1 {
			continue // 8 and 2 coins are even, so remainder must be even
		}
		edges1 := edgesAvail - 4*c11
		nodes1 := nodesAvail - 3*c11
		if edges1 < 0 || nodes1 < 0 {
			continue
		}
		max8 := int(min(rem1/8, int64(edges1/3)))
		if nodes1/2 < max8 {
			max8 = nodes1 / 2
		}
		for c8 = max8; c8 >= 0; c8-- {
			rem2 := rem1 - int64(8*c8)
			if rem2 < 0 {
				continue
			}
			if rem2%2 != 0 {
				continue
			}
			edges2 := edges1 - 3*c8
			nodes2 := nodes1 - 2*c8
			c2 = int(rem2 / 2)
			if edges2 >= c2 && nodes2 >= c2 {
				return c11, c8, c2, true
			}
		}
	}
	return 0, 0, 0, false
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var K int64
	if _, err := fmt.Fscan(in, &K); err != nil {
		return
	}

	impossible := map[int64]bool{1: true, 3: true, 5: true, 7: true, 9: true}
	if impossible[K] {
		fmt.Println("-1 -1")
		return
	}

	edges := make([]Edge, 0)
	nextNode := 2
	nodesUsed := 1
	edgesUsed := 0
	remaining := K

	addGadget := func(l int) {
		// Nodes: g1 ... gL where g1 = nextNode
		g1 := nextNode
		g2 := g1 + 1
		gL := g1 + l - 1
		edges = append(edges, Edge{1, g1})      // anchor to global node 1
		edges = append(edges, Edge{g1, g2})     // (g1, g2)
		edges = append(edges, Edge{g1, gL})     // (g1, gL)
		for x := g1 + 2; x <= gL; x++ {          // g3 .. gL
			edges = append(edges, Edge{g2, x})  // (g2, gx)
		}
		for x := g1 + 2; x <= gL; x++ {          // (gx, gx-1)
			edges = append(edges, Edge{x, x - 1})
		}
		nextNode += l
		nodesUsed += l
		edgesUsed += 2*l - 1
		remaining -= gadgetValue(l)
	}

	// Build gadgets until the remainder can be handled by coin gadgets.
	for {
		edgesAvail := maxM - edgesUsed
		nodesAvail := maxN - nodesUsed
		if remaining < 0 || edgesAvail < 0 || nodesAvail <= 0 {
			fmt.Println("-1 -1")
			return
		}

		// Try to finish with coins.
		if remaining == 0 {
			break
		}
		if impossible[remaining] {
			// Need another gadget to move away from unreachable small values.
		} else {
			c11, c8, c2, ok := findCoins(remaining, edgesAvail, nodesAvail)
			if ok {
				// Add coin gadgets.
				for i := 0; i < c11; i++ {
					a := nextNode
					b := nextNode + 1
					c := nextNode + 2
					edges = append(edges, Edge{1, a}, Edge{1, b}, Edge{a, b}, Edge{a, c})
					nextNode += 3
					nodesUsed += 3
					edgesUsed += 4
				}
				for i := 0; i < c8; i++ {
					a := nextNode
					b := nextNode + 1
					edges = append(edges, Edge{1, a}, Edge{1, b}, Edge{a, b})
					nextNode += 2
					nodesUsed += 2
					edgesUsed += 3
				}
				for i := 0; i < c2; i++ {
					edges = append(edges, Edge{1, nextNode})
					nextNode++
					nodesUsed++
					edgesUsed++
				}
				remaining = 0
				break
			}
		}

		// Need another chain gadget.
		if nodesAvail < 4 || edgesAvail < 5 {
			fmt.Println("-1 -1")
			return
		}
		target := int(math.Sqrt(float64(remaining))) + 1
		limitByEdges := (edgesAvail + 1) / 2
		limitByNodes := nodesAvail - 1
		l := target
		if l > limitByEdges {
			l = limitByEdges
		}
		if l > limitByNodes {
			l = limitByNodes
		}
		if l < 3 {
			fmt.Println("-1 -1")
			return
		}
		// Ensure gadget value does not exceed remaining.
		for l >= 3 && gadgetValue(l) > remaining {
			l--
		}
		if l < 3 {
			fmt.Println("-1 -1")
			return
		}
		addGadget(l)
	}

	if remaining != 0 || len(edges) > maxM || nextNode-1 > maxN {
		fmt.Println("-1 -1")
		return
	}

	out := bufio.NewWriter(os.Stdout)
	fmt.Fprintf(out, "%d %d\n", nextNode-1, len(edges))
	for _, e := range edges {
		fmt.Fprintf(out, "%d %d\n", e.u, e.v)
	}
	out.Flush()
}
