package main

import (
	"bufio"
	"fmt"
	"os"
)

const inf int64 = 1 << 60

var (
	n int
	b int64
	c []int64
	d []int64
	g [][]int
)

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

// dfs returns dp0, dp1 slices and subtree size
func dfs(u int) ([]int64, []int64, int) {
	dp0 := []int64{0, c[u]}
	dp1 := []int64{inf, c[u] - d[u]}
	size := 1
	for _, v := range g[u] {
		child0, child1, sz := dfs(v)
		new0 := make([]int64, size+sz+1)
		new1 := make([]int64, size+sz+1)
		for i := range new0 {
			new0[i] = inf
			new1[i] = inf
		}
		for i := 0; i <= size; i++ {
			for j := 0; j <= sz; j++ {
				if dp0[i] != inf && child0[j] != inf {
					if val := dp0[i] + child0[j]; val < new0[i+j] {
						new0[i+j] = val
					}
				}
				if dp1[i] != inf {
					useChild := child0[j]
					if child1[j] < useChild {
						useChild = child1[j]
					}
					if useChild != inf {
						if val := dp1[i] + useChild; val < new1[i+j] {
							new1[i+j] = val
						}
					}
				}
			}
		}
		dp0 = new0
		dp1 = new1
		size += sz
	}
	return dp0, dp1, size
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	if _, err := fmt.Fscan(reader, &n, &b); err != nil {
		return
	}
	c = make([]int64, n+1)
	d = make([]int64, n+1)
	g = make([][]int, n+1)
	for i := 1; i <= n; i++ {
		var ci, di int64
		fmt.Fscan(reader, &ci, &di)
		c[i] = ci
		d[i] = di
		if i > 1 {
			var x int
			fmt.Fscan(reader, &x)
			g[x] = append(g[x], i)
		}
	}
	dp0, dp1, _ := dfs(1)
	ans := 0
	for k := 0; k < len(dp0); k++ {
		best := dp0[k]
		if dp1[k] < best {
			best = dp1[k]
		}
		if best <= b && k > ans {
			ans = k
		}
	}
	fmt.Fprintln(writer, ans)
}
