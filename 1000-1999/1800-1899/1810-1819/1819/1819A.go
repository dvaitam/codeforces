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

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		freq := make([]int, n+3)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			if a[i] <= n+1 {
				freq[a[i]]++
			}
		}
		m := 0
		for freq[m] > 0 {
			m++
		}
		// find first and last position of m+1
		L, R := -1, -1
		for i, v := range a {
			if v == m+1 {
				if L == -1 {
					L = i
				}
				R = i
			}
		}
		if L == -1 {
			// m+1 not present
			ok := false
			for _, v := range a {
				if v > m {
					ok = true
					break
				}
			}
			if !ok {
				for x := 0; x < m; x++ {
					if freq[x] > 1 {
						ok = true
						break
					}
				}
			}
			if ok {
				fmt.Fprintln(out, "Yes")
			} else {
				fmt.Fprintln(out, "No")
			}
		} else {
			left := make([]int, m)
			right := make([]int, m)
			for i := 0; i < m; i++ {
				left[i] = -1
				right[i] = -1
			}
			for i, v := range a {
				if v < m {
					if left[v] == -1 {
						left[v] = i
					}
					right[v] = i
				}
			}
			ok := true
			for x := 0; x < m; x++ {
				if left[x] >= L && right[x] <= R {
					ok = false
					break
				}
			}
			if ok {
				fmt.Fprintln(out, "Yes")
			} else {
				fmt.Fprintln(out, "No")
			}
		}
	}
}
