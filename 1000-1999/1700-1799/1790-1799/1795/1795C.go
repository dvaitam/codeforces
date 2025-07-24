package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &b[i])
		}

		pref := make([]int64, n)
		for i := 0; i < n; i++ {
			if i == 0 {
				pref[i] = b[i]
			} else {
				pref[i] = pref[i-1] + b[i]
			}
		}

		diff := make([]int64, n+1)
		ans := make([]int64, n)

		for i := 0; i < n; i++ {
			x := a[i]
			if i > 0 {
				x += pref[i-1]
			}
			// pos is the first index with prefix sum > x
			pos := sort.Search(n, func(j int) bool { return pref[j] > x })
			diff[i]++
			diff[pos]--
			if pos < n {
				var prev int64
				if pos > 0 {
					prev = pref[pos-1]
				}
				ans[pos] += x - prev
			}
		}

		curr := int64(0)
		for i := 0; i < n; i++ {
			curr += diff[i]
			ans[i] += curr * b[i]
		}

		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(writer, " ")
			}
			fmt.Fprint(writer, ans[i])
		}
		fmt.Fprintln(writer)
	}
}
