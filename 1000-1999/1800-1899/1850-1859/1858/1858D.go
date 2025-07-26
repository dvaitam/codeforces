package main

import (
	"bufio"
	"fmt"
	"os"
)

func maxLen(s string, k int, target byte) []int {
	n := len(s)
	res := make([]int, k+1)
	for c := 0; c <= k; c++ {
		left, cnt, best := 0, 0, 0
		for right := 0; right < n; right++ {
			if s[right] != target {
				cnt++
			}
			for cnt > c {
				if s[left] != target {
					cnt--
				}
				left++
			}
			if cur := right - left + 1; cur > best {
				best = cur
			}
		}
		res[c] = best
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		var s string
		fmt.Fscan(reader, &s)

		len0 := maxLen(s, k, '0')
		len1 := maxLen(s, k, '1')

		for a := 1; a <= n; a++ {
			best := 0
			for c0 := 0; c0 <= k; c0++ {
				c1 := k - c0
				l0 := len0[c0]
				l1 := len1[c1]
				if l0+l1 > n {
					continue
				}
				val := a*l0 + l1
				if val > best {
					best = val
				}
			}
			if a > 1 {
				fmt.Fprint(writer, " ")
			}
			fmt.Fprint(writer, best)
		}
		fmt.Fprintln(writer)
	}
}
