package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const INF int64 = 1 << 60

// computeCost returns minimal cost to move from start to target (inclusive)
// along the DAG defined by arrays a, b and nxt.
// n is the number of squares. target can be n+1 representing T.
func computeCost(n int, a, b []int64, nxt []int, start, target int) int64 {
	if start == target {
		return 0
	}
	dp := make([]int64, target+2)
	for i := range dp {
		dp[i] = INF
	}
	dp[start] = 0
	for i := start; i < target && i <= n; i++ {
		if dp[i] == INF {
			continue
		}
		// move to i+1
		if i < n && i+1 <= target {
			if v := dp[i] + a[i]; v < dp[i+1] {
				dp[i+1] = v
			}
		}
		// jump to next greater
		ng := nxt[i]
		if ng == 0 {
			ng = n + 1
		}
		if ng <= target {
			if v := dp[i] + b[i]; v < dp[ng] {
				dp[ng] = v
			}
		}
	}
	return dp[target]
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, q int
	if _, err := fmt.Fscan(reader, &n, &q); err != nil {
		return
	}
	p := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &p[i])
	}
	a := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		var x int64
		fmt.Fscan(reader, &x)
		a[i] = x
	}
	b := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		var x int64
		fmt.Fscan(reader, &x)
		b[i] = x
	}

	// compute next greater index to the right for each position
	nxt := make([]int, n+1)
	stack := []int{}
	for i := n; i >= 1; i-- {
		for len(stack) > 0 && p[stack[len(stack)-1]] <= p[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) > 0 {
			nxt[i] = stack[len(stack)-1]
		} else {
			nxt[i] = 0
		}
		stack = append(stack, i)
	}

	S := make(map[int]bool)
	for ; q > 0; q-- {
		var x int
		fmt.Fscan(reader, &x)
		if S[x] {
			delete(S, x)
		} else {
			S[x] = true
		}
		keys := make([]int, 0, len(S))
		for k := range S {
			keys = append(keys, k)
		}
		sort.Ints(keys)
		prev := 1
		var ans int64 = 0
		for _, v := range keys {
			ans += computeCost(n, a, b, nxt, prev, v)
			prev = v
		}
		ans += computeCost(n, a, b, nxt, prev, n+1)
		fmt.Fprintln(writer, ans)
	}
}
