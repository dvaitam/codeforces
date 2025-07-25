package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// This solution builds a subset of the vertices such that
// no divisibility chain of length three remains. The algorithm
// processes values in increasing order and greedily keeps a
// vertex only if together with previously kept divisors it does
// not create such a chain. The resulting graph is bipartite and
// the number of removed vertices is reported.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		sort.Ints(a)

		kept := make(map[int]bool)
		depth := make(map[int]int)
		keepCnt := 0

		for _, val := range a {
			maxD := 0
			// check all divisors of val that were kept before
			for d := 1; d*d <= val; d++ {
				if val%d == 0 {
					if kept[d] && depth[d] > maxD {
						maxD = depth[d]
					}
					other := val / d
					if other != d && kept[other] && depth[other] > maxD {
						maxD = depth[other]
					}
				}
			}
			if maxD+1 <= 2 { // keep the vertex only if chain length stays <=2
				kept[val] = true
				depth[val] = maxD + 1
				keepCnt++
			}
		}

		fmt.Fprintln(out, n-keepCnt)
	}
}
