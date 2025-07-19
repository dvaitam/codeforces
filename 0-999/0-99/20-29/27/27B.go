package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	// h[i]: number of remaining matches for player i
	// a[i]: score offset, starts at 50, +1 for win, -1 for loss
	h := make([]int, n+1)
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		h[i] = n - 1
		a[i] = 50
	}
	total := n * (n - 1) / 2
	// Read known match results (total-1 matches)
	for i := 1; i < total; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		h[u]--
		h[v]--
		a[u]++
		a[v]--
	}
	var x, y int
	b := make(map[int]bool)
	// Identify two players with one missing match and record scores of completed ones
	for i := 1; i <= n; i++ {
		if h[i] != 0 {
			if x == 0 {
				x = i
			} else {
				y = i
			}
		} else {
			b[a[i]] = true
		}
	}
	// Determine result: check if y wins (x loses) yields unique scores
	if !b[a[x]-1] && !b[a[y]+1] {
		fmt.Printf("%d %d\n", y, x)
	} else {
		fmt.Printf("%d %d\n", x, y)
	}
}
