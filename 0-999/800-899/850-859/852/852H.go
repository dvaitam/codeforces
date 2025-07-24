package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

// Point represents a 2D point with integer coordinates
type Point struct {
	x, y int64
}

// cross returns cross product (b-a)x(c-a)
func cross(a, b, c Point) int64 {
	return (b.x-a.x)*(c.y-a.y) - (b.y-a.y)*(c.x-a.x)
}

// area2 returns twice the signed area of triangle abc
func area2(a, b, c Point) float64 {
	return float64(cross(a, b, c)) / 2.0
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, K int
	if _, err := fmt.Fscan(in, &n, &K); err != nil {
		return
	}
	pts := make([]Point, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &pts[i].x, &pts[i].y)
	}

	// precompute edge visibility: edgeOK[i][j] is true if no point lies to the left of directed edge i->j
	edgeOK := make([][]bool, n)
	for i := range edgeOK {
		edgeOK[i] = make([]bool, n)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i == j {
				continue
			}
			ok := true
			for k := 0; k < n; k++ {
				if k == i || k == j {
					continue
				}
				if cross(pts[i], pts[j], pts[k]) > 0 {
					ok = false
					break
				}
			}
			edgeOK[i][j] = ok
		}
	}

	best := 0.0

	for r := 0; r < n; r++ {
		// build list of other points
		idx := make([]int, 0, n-1)
		for i := 0; i < n; i++ {
			if i == r {
				continue
			}
			idx = append(idx, i)
		}
		// sort by polar angle around pts[r]
		sort.Slice(idx, func(a, b int) bool {
			ax := float64(pts[idx[a]].x - pts[r].x)
			ay := float64(pts[idx[a]].y - pts[r].y)
			bx := float64(pts[idx[b]].x - pts[r].x)
			by := float64(pts[idx[b]].y - pts[r].y)
			return math.Atan2(ay, ax) < math.Atan2(by, bx)
		})
		m := len(idx)
		if m < K-1 { // not enough points besides r
			continue
		}

		// dp[i][j][t] -> map[start]area
		dp := make([][][]map[int]float64, m)
		for i := range dp {
			dp[i] = make([][]map[int]float64, m)
			for j := range dp[i] {
				dp[i][j] = make([]map[int]float64, K+1)
			}
		}

		// initialization for triangles
		for a := 0; a < m; a++ {
			ia := idx[a]
			if !edgeOK[r][ia] {
				continue
			}
			for b := a + 1; b < m; b++ {
				ib := idx[b]
				if !edgeOK[ia][ib] {
					continue
				}
				if cross(pts[r], pts[ia], pts[ib]) <= 0 {
					continue
				}
				area := area2(pts[r], pts[ia], pts[ib])
				if dp[a][b][3] == nil {
					dp[a][b][3] = make(map[int]float64)
				}
				dp[a][b][3][a] = area
			}
		}

		for i := 0; i < m; i++ {
			for j := i + 1; j < m; j++ {
				for t := 3; t < K; t++ {
					mp := dp[i][j][t]
					if mp == nil {
						continue
					}
					for start, val := range mp {
						pi := idx[i]
						pj := idx[j]
						for k := j + 1; k < m; k++ {
							pk := idx[k]
							if cross(pts[pi], pts[pj], pts[pk]) <= 0 {
								continue
							}
							if !edgeOK[pj][pk] {
								continue
							}
							area := val + area2(pts[r], pts[pj], pts[pk])
							if dp[j][k][t+1] == nil {
								dp[j][k][t+1] = make(map[int]float64)
							}
							if area > dp[j][k][t+1][start] {
								dp[j][k][t+1][start] = area
							}
						}
					}
				}
			}
		}

		// finalize polygons
		for i := 0; i < m; i++ {
			for j := i + 1; j < m; j++ {
				mp := dp[i][j][K]
				if mp == nil {
					continue
				}
				pi := idx[i]
				pj := idx[j]
				if !edgeOK[pj][r] {
					continue
				}
				if cross(pts[pi], pts[pj], pts[r]) <= 0 {
					continue
				}
				for start, val := range mp {
					s := idx[start]
					if cross(pts[pj], pts[r], pts[s]) <= 0 {
						continue
					}
					if val > best {
						best = val
					}
				}
			}
		}
	}

	fmt.Printf("%.2f\n", best)
}
