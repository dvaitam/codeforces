package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const inf = int(1e9)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)

	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		var s string
		fmt.Fscan(in, &s)

		limits := make([]int, n)
		for i := range limits {
			limits[i] = n
		}

		cur := inf
		for v := n; v >= 1; v-- {
			if s[v-1] == 'p' {
				cur = v
			}
			if cur != inf && limits[v-1] > cur {
				limits[v-1] = cur
			}
		}

		for i := 1; i <= n; i++ {
			if s[i-1] != 's' {
				continue
			}
			startVal := n - i + 2
			if startVal < 1 {
				startVal = 1
			}
			for v := startVal; v <= n; v++ {
				if limits[v-1] > i-1 {
					limits[v-1] = i - 1
				}
			}
		}

		sort.Ints(limits)
		ok := true
		for i := 0; i < n; i++ {
			if limits[i] < i+1 {
				ok = false
				break
			}
		}

		if ok {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
