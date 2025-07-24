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

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		if _, err := fmt.Fscan(reader, &n); err != nil {
			return
		}
		freq := make([]int, n+1)
		maxFreq := 0
		for i := 0; i < n; i++ {
			var c int
			fmt.Fscan(reader, &c)
			freq[c]++
			if freq[c] > maxFreq {
				maxFreq = freq[c]
			}
		}

		counts := make([]int, 0)
		for _, v := range freq {
			if v > 0 {
				counts = append(counts, v)
			}
		}

		valid := func(s int) bool {
			for _, x := range counts {
				k := (x + s - 1) / s
				if k*(s-1) > x {
					return false
				}
			}
			return true
		}

		lo := 1
		hi := maxFreq + 1
		for lo < hi {
			mid := (lo + hi + 1) / 2
			if valid(mid) {
				lo = mid
			} else {
				hi = mid - 1
			}
		}
		s := lo
		m := 0
		for _, x := range counts {
			m += (x + s - 1) / s
		}
		fmt.Fprintln(writer, m)
	}
}
