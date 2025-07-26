package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Lamp struct {
	a int
	b int
}

// DSU find for scheduling slots
func find(parent []int, x int) int {
	if x <= 0 {
		return 0
	}
	if parent[x] != x {
		parent[x] = find(parent, parent[x])
	}
	return parent[x]
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		lamps := make([]Lamp, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &lamps[i].a, &lamps[i].b)
		}
		sort.Slice(lamps, func(i, j int) bool {
			return lamps[i].b > lamps[j].b
		})
		parent := make([]int, n+1)
		for i := 0; i <= n; i++ {
			parent[i] = i
		}
		ans := int64(0)
		for _, l := range lamps {
			slot := find(parent, l.a)
			if slot > 0 {
				ans += int64(l.b)
				parent[slot] = find(parent, slot-1)
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
