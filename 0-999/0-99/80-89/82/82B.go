package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	// Special case n == 2: only one union, split arbitrarily
	if n == 2 {
		var k int
		fmt.Fscan(in, &k)
		a := make([]int, k)
		for i := 0; i < k; i++ {
			fmt.Fscan(in, &a[i])
		}
		// output first element as one set, rest as the other
		fmt.Printf("1 %d\n", a[0])
		fmt.Printf("%d", k-1)
		for i := 1; i < k; i++ {
			fmt.Printf(" %d", a[i])
		}
		fmt.Println()
		return
	}
	// General case: build co-occurrence counts
	var b [201][201]int
	var seen [201]bool
	var used [201]bool
	m := n * (n - 1) / 2
	for t := 0; t < m; t++ {
		var k int
		fmt.Fscan(in, &k)
		a := make([]int, k)
		for i := 0; i < k; i++ {
			fmt.Fscan(in, &a[i])
			seen[a[i]] = true
		}
		// count pairs
		for i := 0; i < k; i++ {
			for j := i + 1; j < k; j++ {
				x, y := a[i], a[j]
				b[x][y]++
				b[y][x]++
			}
		}
	}
	// Recover sets: numbers co-occurring in all unions (n-1 times)
	for x := 1; x <= 200; x++ {
		if seen[x] && !used[x] {
			// start new set with x
			members := []int{x}
			used[x] = true
			// find all y where x and y always co-occur
			for y := 1; y <= 200; y++ {
				if !used[y] && b[x][y] == n-1 {
					members = append(members, y)
					used[y] = true
				}
			}
			// print this set
			fmt.Printf("%d", len(members))
			for _, v := range members {
				fmt.Printf(" %d", v)
			}
			fmt.Println()
		}
	}
}
