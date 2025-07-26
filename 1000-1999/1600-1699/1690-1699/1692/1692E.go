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
		var n, s int
		fmt.Fscan(in, &n, &s)
		a := make([]int, n)
		total := 0
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			total += a[i]
		}
		if total < s {
			fmt.Fprintln(out, -1)
			continue
		}
		if total == s {
			fmt.Fprintln(out, 0)
			continue
		}

		best := 0
		left := 0
		sum := 0
		for right := 0; right < n; right++ {
			sum += a[right]
			for left <= right && sum > s {
				sum -= a[left]
				left++
			}
			if sum == s {
				length := right - left + 1
				if length > best {
					best = length
				}
			}
		}
		fmt.Fprintln(out, n-best)
	}
}
