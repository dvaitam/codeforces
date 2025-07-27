package main

import (
	"bufio"
	"fmt"
	"os"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}

	dist := make([]int, n+1)
	prev := make([]int, n+1)
	preX := make([]int, n+1)
	for i := 0; i <= n; i++ {
		dist[i] = -1
	}
	queue := []int{0}
	dist[0] = 0
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for x := max(0, k-(n-cur)); x <= min(k, cur); x++ {
			nxt := cur + k - 2*x
			if dist[nxt] == -1 {
				dist[nxt] = dist[cur] + 1
				prev[nxt] = cur
				preX[nxt] = x
				queue = append(queue, nxt)
			}
		}
	}

	if dist[n] == -1 || dist[n] > 500 {
		fmt.Fprintln(out, -1)
		return
	}

	// Reconstruct operations
	ops := []int{}
	cur := n
	for cur > 0 {
		ops = append(ops, preX[cur])
		cur = prev[cur]
	}
	for i, j := 0, len(ops)-1; i < j; i, j = i+1, j-1 {
		ops[i], ops[j] = ops[j], ops[i]
	}

	ones := []int{}
	zeros := make([]int, n)
	for i := 1; i <= n; i++ {
		zeros[i-1] = i
	}

	ans := 0
	for _, x := range ops {
		y := k - x
		if len(ones) < x || len(zeros) < y {
			// should not happen
			return
		}
		query := make([]int, 0, k)
		query = append(query, ones[:x]...)
		query = append(query, zeros[:y]...)
		fmt.Fprint(out, "?")
		for _, v := range query {
			fmt.Fprintf(out, " %d", v)
		}
		fmt.Fprintln(out)
		out.Flush()

		var resp int
		if _, err := fmt.Fscan(in, &resp); err != nil {
			return
		}
		ans ^= resp

		movedOnes := append([]int{}, ones[:x]...)
		movedZeros := append([]int{}, zeros[:y]...)
		ones = ones[x:]
		zeros = zeros[y:]
		ones = append(ones, movedZeros...)
		zeros = append(zeros, movedOnes...)
	}

	fmt.Fprintf(out, "! %d\n", ans)
}
