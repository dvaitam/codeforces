package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: provide a complete implementation for problem E.
// This placeholder attempts a simple heuristic based on moving one
// element as early as possible. It may not pass all tests.

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	pref := make([]int64, n+1)
	for i := 0; i < n; i++ {
		pref[i+1] = pref[i] + a[i]
	}

	best := pref[n] // doing nothing
	for j := 0; j < n; j++ {
		// earliest step we can start using a[j]
		s := (j + 2) / 2
		cand := pref[s] + a[j]*int64(n-s)
		if cand > best {
			best = cand
		}
	}

	fmt.Fprintln(out, best)
}
