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
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(in, &n, &k)

		if perm := buildPermutation(n, k); perm == nil {
			fmt.Fprintln(out, 0)
		} else {
			for i, val := range perm {
				if i > 0 {
					fmt.Fprint(out, " ")
				}
				fmt.Fprint(out, val)
			}
			fmt.Fprintln(out)
		}
	}
}

func buildPermutation(n, k int) []int {
	total := n * (n - 1) / 2
	target := total - k

	dp := make([][]bool, n+1)
	parent := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		dp[i] = make([]bool, total+1)
		parent[i] = make([]int, total+1)
		for j := range parent[i] {
			parent[i][j] = -1
		}
	}

	dp[0][0] = true
	for used := 0; used < n; used++ {
		for sum := 0; sum <= total; sum++ {
			if !dp[used][sum] {
				continue
			}
			for seg := 1; used+seg <= n; seg++ {
				add := seg * (seg - 1) / 2
				nextSum := sum + add
				if nextSum > total || dp[used+seg][nextSum] {
					continue
				}
				dp[used+seg][nextSum] = true
				parent[used+seg][nextSum] = seg
			}
		}
	}

	if !dp[n][target] {
		return nil
	}

	segments := []int{}
	for pos, sum := n, target; pos > 0; {
		seg := parent[pos][sum]
		if seg == -1 {
			return nil
		}
		segments = append(segments, seg)
		pos -= seg
		sum -= seg * (seg - 1) / 2
	}
	for i, j := 0, len(segments)-1; i < j; i, j = i+1, j-1 {
		segments[i], segments[j] = segments[j], segments[i]
	}

	perm := make([]int, 0, n)
	cur := n
	for _, seg := range segments {
		start := cur - seg + 1
		for val := start; val <= cur; val++ {
			perm = append(perm, val)
		}
		cur -= seg
	}
	return perm
}
