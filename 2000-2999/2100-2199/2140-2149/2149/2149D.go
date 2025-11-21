package main

import (
	"bufio"
	"fmt"
	"os"
)

func calc(pos []int) int64 {
	if len(pos) <= 1 {
		return 0
	}
	adjusted := make([]int64, len(pos))
	for i, p := range pos {
		adjusted[i] = int64(p - i)
	}
	median := adjusted[len(pos)/2]
	var res int64
	for _, v := range adjusted {
		diff := v - median
		if diff < 0 {
			diff = -diff
		}
		res += diff
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		var s string
		fmt.Fscan(in, &n)
		fmt.Fscan(in, &s)

		posA := make([]int, 0, n)
		posB := make([]int, 0, n)
		for i := 0; i < n; i++ {
			if s[i] == 'a' {
				posA = append(posA, i+1)
			} else {
				posB = append(posB, i+1)
			}
		}

		ans := calc(posA)
		if v := calc(posB); v < ans {
			ans = v
		}
		fmt.Fprintln(out, ans)
	}
}
