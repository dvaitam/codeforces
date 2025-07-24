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

	var p int
	if _, err := fmt.Fscan(in, &p); err != nil {
		return
	}

	const INF int = 1 << 30
	dp := make([]int, p+1)
	for i := 1; i <= p; i++ {
		dp[i] = INF
	}
	dp[0] = 0

	for s := 2; ; s++ {
		t := s * (s - 1) / 2
		if t > p {
			break
		}
		for j := t; j <= p; j++ {
			if v := dp[j-t] + s; v < dp[j] {
				dp[j] = v
			}
		}
	}

	n := dp[p]
	unidir := n*(n-1)/2 - p
	fmt.Fprintln(out, n, unidir)
}
