package main

import (
	"bufio"
	"fmt"
	"os"
)

// The sequence of played cards must be strictly increasing and uses every
// number from 0 to n*m-1 exactly once, so the play order is forced to be
// 0, 1, 2, ..., n*m-1. At step t (1-indexed) cow p_pos plays, where
// pos = (t-1) mod n + 1. Therefore cow in position pos must own exactly the
// numbers {pos-1 + k*n | 0 <= k < m}. We only need to match each residue
// class modulo n to a unique cow that owns precisely that arithmetic
// progression.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}

	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(in, &n, &m)

		position := make([]int, n) // position[r] = cow index assigned to residue r
		ok := true

		for cow := 1; cow <= n; cow++ {
			cards := make([]int, m)
			for i := 0; i < m; i++ {
				fmt.Fscan(in, &cards[i])
			}

			if !ok {
				// Still need to consume input, but can skip checks.
				continue
			}

			res := cards[0] % n
			seen := make([]bool, m)
			valid := true
			for _, v := range cards {
				if v%n != res {
					valid = false
					break
				}
				q := v / n
				if q < 0 || q >= m || seen[q] {
					valid = false
					break
				}
				seen[q] = true
			}
			if valid {
				for _, s := range seen {
					if !s {
						valid = false
						break
					}
				}
			}

			if valid {
				if position[res] != 0 {
					ok = false // two cows claim same residue
				} else {
					position[res] = cow
				}
			} else {
				ok = false
			}
		}

		if ok {
			for i, v := range position {
				if v == 0 {
					ok = false
					break
				}
				if i > 0 {
					fmt.Fprint(out, " ")
				}
				fmt.Fprint(out, v)
			}
		}

		if !ok {
			fmt.Fprint(out, -1)
		}
		if T > 0 {
			fmt.Fprint(out, "\n")
		}
	}
}
