package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the problem described in problemA.txt.
// It simulates the motion of n scooters on a narrow road where
// scooter i moves i meters per second but cannot overtake the
// scooter in front and must stay at least 1 meter behind.
// After one second we output the distance from each scooter to
// point B.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	dist := make([]int, n)
	if n > 0 {
		dist[0] = a[0] - 1
	}
	for i := 1; i < n; i++ {
		cand := a[i] - (i + 1)
		if cand < dist[i-1]+1 {
			cand = dist[i-1] + 1
		}
		dist[i] = cand
	}

	for i := 0; i < n; i++ {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, dist[i])
	}
	fmt.Fprintln(out)
}
