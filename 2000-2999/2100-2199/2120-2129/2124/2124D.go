package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		a := make([]int, n)
		freq := make([]int, n+2)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			freq[a[i]]++
		}

		if k == 1 {
			fmt.Fprintln(out, "YES")
			continue
		}

		// Find the k-th smallest value.
		kth := 0
		acc := 0
		for v := 1; v <= n; v++ {
			acc += freq[v]
			if acc >= k {
				kth = v
				break
			}
		}

		// Remove all elements greater than kth (they are always deletable).
		filtered := make([]int, 0, n)
		for _, v := range a {
			if v <= kth {
				filtered = append(filtered, v)
			}
		}

		l, r := 0, len(filtered)-1
		length := len(filtered)
		ok := true
		for l < r {
			if filtered[l] == filtered[r] {
				l++
				r--
				continue
			}
			if length < k {
				ok = false
				break
			}
			if filtered[l] == kth {
				l++
				length--
			} else if filtered[r] == kth {
				r--
				length--
			} else {
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
