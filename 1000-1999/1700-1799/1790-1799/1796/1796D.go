package main

import (
	"bufio"
	"fmt"
	"os"
)

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, k int
		var x int64
		fmt.Fscan(reader, &n, &k, &x)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}

		b := make([]int64, n+1)
		prefix := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			b[i] = a[i-1] - x
			prefix[i] = prefix[i-1] + b[i]
		}

		res := int64(0)
		if x >= 0 {
			minPref := make([]int64, n+1)
			minPref[0] = 0
			for i := 1; i <= n; i++ {
				if prefix[i] < minPref[i-1] {
					minPref[i] = prefix[i]
				} else {
					minPref[i] = minPref[i-1]
				}
			}
			for r := 1; r <= n; r++ {
				lim := k
				if r < lim {
					lim = r
				}
				for l := 1; l <= lim; l++ {
					val := prefix[r] - prefix[r-l] + int64(2*l)*x
					if val > res {
						res = val
					}
				}
				if r > k {
					val := prefix[r] + int64(2*k)*x - minPref[r-k]
					if val > res {
						res = val
					}
				}
			}
			if res < 0 {
				res = 0
			}
			fmt.Fprintln(writer, res)
			continue
		}

		// x < 0
		cur := int64(0)
		for i := 1; i <= n; i++ {
			cur += b[i]
			if cur < 0 {
				cur = 0
			}
			if cur > res {
				res = cur
			}
		}
		start := n - k + 1
		if start < 1 {
			start = 1
		}
		for length := start; length <= n; length++ {
			add := int64(2) * x * int64(length+k-n)
			window := prefix[length] - prefix[0]
			val := window + add
			if val > res {
				res = val
			}
			for r := length + 1; r <= n; r++ {
				window += b[r] - b[r-length]
				val = window + add
				if val > res {
					res = val
				}
			}
		}
		if res < 0 {
			res = 0
		}
		fmt.Fprintln(writer, res)
	}
}
