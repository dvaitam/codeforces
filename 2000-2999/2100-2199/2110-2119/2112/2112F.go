package main

import (
	"bufio"
	"fmt"
	"os"
)

const inf int64 = 1 << 60

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	// d1[u][v] – weight of direct edge u -> v (0 on diagonal), INF if none.
	d1 := make([][]int64, n)
	dist := make([][]int64, n)
	for i := 0; i < n; i++ {
		d1[i] = make([]int64, n)
		dist[i] = make([]int64, n)
		for j := 0; j < n; j++ {
			if i == j {
				d1[i][j] = 0
				dist[i][j] = 0
			} else {
				d1[i][j] = inf
				dist[i][j] = inf
			}
		}
	}

	for i := 0; i < m; i++ {
		var x, y int
		var z int64
		fmt.Fscan(in, &x, &y, &z)
		x--
		y--
		if z < d1[y][x] {
			d1[y][x] = z
			dist[y][x] = z
		}
	}

	// Floyd–Warshall for all-pairs shortest paths.
	for k := 0; k < n; k++ {
		for i := 0; i < n; i++ {
			if dist[i][k] == inf {
				continue
			}
			ik := dist[i][k]
			rowK := dist[k]
			for j := 0; j < n; j++ {
				if rowK[j] == inf {
					continue
				}
				if ik+rowK[j] < dist[i][j] {
					dist[i][j] = ik + rowK[j]
				}
			}
		}
	}

	// Precompute candidate lists where a multi-edge path can beat a direct edge.
	cand := make([][]int, n)
	for i := 0; i < n; i++ {
		for p := 0; p < n; p++ {
			if dist[p][i] < inf && dist[p][i] < d1[p][i] {
				cand[i] = append(cand[i], p)
			}
		}
	}

	var q int
	fmt.Fscan(in, &q)
	a := make([]int64, n)
	best1 := make([]int64, n)
	best1Idx := make([]int, n)
	best2 := make([]int64, n)

	for ; q > 0; q-- {
		var k int64
		fmt.Fscan(in, &k)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		// For each target i, find two smallest values of a[j] + d1[j][i].
		for i := 0; i < n; i++ {
			best1[i], best2[i] = inf, inf
			best1Idx[i] = -1
		}

		for j := 0; j < n; j++ {
			aj := a[j]
			for i := 0; i < n; i++ {
				if d1[j][i] == inf {
					continue
				}
				val := aj + d1[j][i]
				if val < best1[i] {
					best2[i] = best1[i]
					best1[i] = val
					best1Idx[i] = j
				} else if val < best2[i] {
					best2[i] = val
				}
			}
		}

		ans := make([]byte, n)
		for i := 0; i < n; i++ {
			ans[i] = '0'
			for _, p := range cand[i] {
				bw := best1[i]
				if best1Idx[i] == p {
					bw = best2[i]
				}
				if bw == inf {
					continue
				}
				need := a[p] + dist[p][i] - bw + 1
				if need <= 0 || need <= k {
					ans[i] = '1'
					break
				}
			}
		}

		fmt.Fprintln(out, string(ans))
	}
}
