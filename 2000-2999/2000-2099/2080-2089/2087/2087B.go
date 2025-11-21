package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)

		m := 2 * n
		ratings := make([]int, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(in, &ratings[i])
		}

		sort.Ints(ratings)
		possible := true

		for i := 0; i < m; i += 2 {
			j := i + 1
			diff := ratings[j] - ratings[i]

			if i > 0 {
				if diff > ratings[i]-ratings[i-1] {
					possible = false
					break
				}
			}

			if j < m-1 {
				if diff > ratings[j+1]-ratings[j] {
					possible = false
					break
				}
			}
		}

		if possible {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
