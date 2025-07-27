package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	var s string
	fmt.Fscan(reader, &s)

	patterns := []string{"abc", "acb", "bac", "bca", "cab", "cba"}
	pref := make([][]int, len(patterns))
	for i := range pref {
		pref[i] = make([]int, n+1)
	}

	for p, pat := range patterns {
		for i := 0; i < n; i++ {
			mismatch := 0
			if s[i] != pat[i%3] {
				mismatch = 1
			}
			pref[p][i+1] = pref[p][i] + mismatch
		}
	}

	for ; m > 0; m-- {
		var l, r int
		fmt.Fscan(reader, &l, &r)
		best := n
		for p := range patterns {
			cost := pref[p][r] - pref[p][l-1]
			if cost < best {
				best = cost
			}
		}
		fmt.Fprintln(writer, best)
	}
}
