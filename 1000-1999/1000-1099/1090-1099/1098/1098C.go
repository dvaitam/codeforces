package main

import (
	"bufio"
	"fmt"
	"os"
)

// getSum computes sum of sequence: starting with m = n, decrement by k each time and k *= i
func getSum(n, i, k int64) int64 {
	var sum int64
	m := n
	for m > 0 {
		sum += m
		m -= k
		k *= i
	}
	return sum
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, s int64
	if _, err := fmt.Fscan(reader, &n, &s); err != nil {
		return
	}
	// try each possible block size i
	for i := int64(1); i < n; i++ {
		total := getSum(n, i, 1)
		if i == 1 && total < s {
			fmt.Fprintln(writer, "No")
			return
		}
		if total > s {
			continue
		}
		// found feasible i
		fmt.Fprintln(writer, "Yes")
		// case i == 1: identity permutation
		if i == 1 {
			for j := int64(1); j <= n; j++ {
				if j > 1 {
					writer.WriteByte(' ')
				}
				fmt.Fprint(writer, j)
			}
			fmt.Fprintln(writer)
			return
		}
		// general case
		k := i
		m := n - 1
		// we assume the largest element n is placed as first block
		s -= n
		start := int64(1)
		finish := int64(1)
		// build remaining elements
		for m > 0 {
			// find smallest e in [1..k] such that getSum(m, i, e) <= s
			l, r := int64(1), k
			for l < r {
				e := (l + r) / 2
				if getSum(m, i, e) > s {
					l = e + 1
				} else {
					r = e
				}
			}
			k = l
			// print k elements: blocks of size i with increasing labels starting from start
			for j := int64(0); j < k; j++ {
				if start+(j/i) > n {
					// safety check
					fmt.Fprint(writer, n)
				} else {
					fmt.Fprint(writer, start+(j/i))
				}
				writer.WriteByte(' ')
			}
			// update for next iteration
			start = finish + 1
			finish += k
			s -= m
			m -= k
			k *= i
		}
		fmt.Fprintln(writer)
		return
	}
	fmt.Fprintln(writer, "No")
}
