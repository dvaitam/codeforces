package main

import (
	"bufio"
	"fmt"
	"os"
)

const MaxB = 1000

var dist [MaxB + 1]int

func init() {
	const INF = int(1e9)
	for i := range dist {
		dist[i] = INF
	}
	dist[1] = 0
	q := []int{1}
	for head := 0; head < len(q); head++ {
		v := q[head]
		for x := 1; x <= v; x++ {
			nxt := v + v/x
			if nxt <= MaxB && dist[nxt] > dist[v]+1 {
				dist[nxt] = dist[v] + 1
				q = append(q, nxt)
			}
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

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
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}
		c := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &c[i])
		}
		w := make([]int, n)
		for i := 0; i < n; i++ {
			w[i] = dist[b[i]]
		}
		cap := min(k, 12*n)
		dp := make([]int, cap+1)
		for i := 0; i < n; i++ {
			cost := w[i]
			val := c[i]
			if cost > cap {
				continue
			}
			for j := cap; j >= cost; j-- {
				if dp[j-cost]+val > dp[j] {
					dp[j] = dp[j-cost] + val
				}
			}
		}
		fmt.Fprintln(out, dp[min(k, cap)])
	}
}
