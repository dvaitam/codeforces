package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	a := make([]int, n)
	sum := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
		sum += a[i]
	}

	half := sum / 2
	dp := make([]bool, half+1)
	dp[0] = true
	for _, x := range a {
		for s := half; s >= x; s-- {
			if dp[s-x] {
				dp[s] = true
			}
		}
	}

	if sum%2 == 1 || !dp[half] {
		fmt.Fprintln(out, 0)
		return
	}

	g := a[0]
	for i := 1; i < n; i++ {
		g = gcd(g, a[i])
	}

	idx := -1
	best := 31
	for i, v := range a {
		x := v / g
		cnt := 0
		for x%2 == 0 {
			x /= 2
			cnt++
		}
		if cnt < best {
			best = cnt
			idx = i
		}
	}

	fmt.Fprintln(out, 1)
	fmt.Fprintln(out, idx+1)
}
