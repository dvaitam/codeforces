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

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	a := make([]int, n)
	maxA := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
		if a[i] > maxA {
			maxA = a[i]
		}
	}

	limit := 200000
	if maxA > limit {
		limit = maxA
	}
	dist := make([][]int, limit+1)

	for _, x := range a {
		steps := 0
		for {
			dist[x] = append(dist[x], steps)
			if x == 0 {
				break
			}
			x /= 2
			steps++
		}
	}

	ans := int(^uint(0) >> 1) // max int
	for i := 0; i <= limit; i++ {
		if len(dist[i]) < k {
			continue
		}
		sort.Ints(dist[i])
		sum := 0
		for j := 0; j < k; j++ {
			sum += dist[i][j]
		}
		if sum < ans {
			ans = sum
		}
	}

	fmt.Fprintln(out, ans)
}
