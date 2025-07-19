package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// Node represents a vertex with its original index and weight
type Node struct {
	num int   // original index (1-based)
	wei int64 // weight
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for T > 0 {
		T--
		var n, m int
		fmt.Fscan(reader, &n, &m)
		nodes := make([]Node, n)
		var sum int64
		for i := 0; i < n; i++ {
			var w int64
			fmt.Fscan(reader, &w)
			nodes[i] = Node{num: i + 1, wei: w}
			sum += w
		}
		// sort by weight, tie-breaker by index
		sort.Slice(nodes, func(i, j int) bool {
			if nodes[i].wei == nodes[j].wei {
				return nodes[i].num < nodes[j].num
			}
			return nodes[i].wei < nodes[j].wei
		})

		// impossible if fewer edges than vertices or exactly two vertices
		if m < n || n == 2 {
			writer.WriteString("-1\n")
			continue
		}
		// branch: not enough edges for full star, use cycle + extras
		if n >= 3 && m < 2*n-3 {
			ans := sum * 2
			tEdges := m - n
			if tEdges > 0 {
				ans += int64(tEdges) * (nodes[0].wei + nodes[1].wei)
			}
			fmt.Fprintf(writer, "%d\n", ans)
			// base cycle 1-2-...-n-1
			for i := 1; i < n; i++ {
				fmt.Fprintf(writer, "%d %d\n", i, i+1)
			}
			fmt.Fprintf(writer, "%d %d\n", n, 1)
			// extra edges between two smallest weights
			for i := 0; i < tEdges; i++ {
				fmt.Fprintf(writer, "%d %d\n", nodes[0].num, nodes[1].num)
			}
		} else {
			// full star from two smallest to all others
			var ans int64
			for j := 2; j < n; j++ {
				ans += 2*nodes[j].wei + nodes[0].wei + nodes[1].wei
			}
			baseEdges := 2 * (n - 2)
			tEdges := m - baseEdges
			ans += int64(tEdges) * (nodes[0].wei + nodes[1].wei)
			fmt.Fprintf(writer, "%d\n", ans)
			// connect each other vertex to the two smallest
			for j := 2; j < n; j++ {
				fmt.Fprintf(writer, "%d %d\n", nodes[0].num, nodes[j].num)
				fmt.Fprintf(writer, "%d %d\n", nodes[1].num, nodes[j].num)
			}
			// extra edges between two smallest weights
			for i := 0; i < tEdges; i++ {
				fmt.Fprintf(writer, "%d %d\n", nodes[0].num, nodes[1].num)
			}
		}
	}
}
