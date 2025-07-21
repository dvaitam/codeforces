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
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	// initial best days is limited by first and last tile
	best := a[0]
	if a[n-1] < best {
		best = a[n-1]
	}
	// check adjacent pairs for two-tile gap
	for i := 0; i+1 < n; i++ {
		// the walk becomes impossible the day after both tiles i and i+1 are destroyed
		// but they both are destroyed together when day > max(a[i], a[i+1])
		// thus last possible day is max(a[i], a[i+1])
		mx := a[i]
		if a[i+1] > mx {
			mx = a[i+1]
		}
		if mx < best {
			best = mx
		}
	}
	fmt.Println(best)
}
