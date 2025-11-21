package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		a := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &a[i])
		}

		prevLess := make([]int, n+1)
		stack := make([]int, 0)
		for i := 1; i <= n; i++ {
			for len(stack) > 0 && a[stack[len(stack)-1]] >= a[i] {
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				prevLess[i] = 0
			} else {
				prevLess[i] = stack[len(stack)-1]
			}
			stack = append(stack, i)
		}

		prevGreater := make([]int, n+1)
		stack = stack[:0]
		for i := 1; i <= n; i++ {
			for len(stack) > 0 && a[stack[len(stack)-1]] <= a[i] {
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				prevGreater[i] = 0
			} else {
				prevGreater[i] = stack[len(stack)-1]
			}
			stack = append(stack, i)
		}

		nextGreater := make([]int, n+2)
		stack = stack[:0]
		for i := n; i >= 1; i-- {
			for len(stack) > 0 && a[stack[len(stack)-1]] <= a[i] {
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				nextGreater[i] = n + 1
			} else {
				nextGreater[i] = stack[len(stack)-1]
			}
			stack = append(stack, i)
		}

		positions := make([][]int, n+1)
		weights := make([][]int64, n+1)
		for i := 1; i <= n; i++ {
			v := a[i]
			positions[v] = append(positions[v], i)
		}
		for v := 1; v <= n; v++ {
			weights[v] = make([]int64, len(positions[v]))
		}
		counters := make([]int, n+1)
		for i := 1; i <= n; i++ {
			v := a[i]
			idx := counters[v]
			weights[v][idx] = int64(i - prevLess[i])
			counters[v]++
		}
		prefix := make([][]int64, n+1)
		for v := 1; v <= n; v++ {
			pre := make([]int64, len(weights[v])+1)
			for i := 0; i < len(weights[v]); i++ {
				pre[i+1] = pre[i] + weights[v][i]
			}
			prefix[v] = pre
		}

		var ans int64
		for j := 1; j <= n; j++ {
			comp := k - a[j]
			if comp < 1 || comp > n {
				continue
			}
			pos := positions[comp]
			if len(pos) == 0 {
				continue
			}
			L := sort.Search(len(pos), func(i int) bool { return pos[i] > prevGreater[j] })
			R := sort.Search(len(pos), func(i int) bool { return pos[i] >= j }) - 1
			if L <= R {
				sumWeights := prefix[comp][R+1] - prefix[comp][L]
				ans += sumWeights * int64(nextGreater[j]-j)
			}
		}
		fmt.Fprintln(out, ans)
	}
}
