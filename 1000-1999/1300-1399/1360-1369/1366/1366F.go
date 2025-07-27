package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const MOD int64 = 1_000_000_007
const INF_NEG int64 = -1 << 60

type Edge struct {
	u, v int
	w    int64
}

type Line struct {
	k int64
	b int64
}

func isBad(a, b, c Line) bool {
	// (b.b - a.b)/(a.k - b.k) >= (c.b - b.b)/(b.k - c.k)
	leftNum := (b.b - a.b) * (b.k - c.k)
	rightNum := (c.b - b.b) * (a.k - b.k)
	return leftNum >= rightNum
}

func floorDiv(a, b int64) int64 {
	// denominator b>0 expected
	if b < 0 {
		a, b = -a, -b
	}
	if a >= 0 {
		return a / b
	}
	// a negative
	return -((-a + b - 1) / b)
}

func mod(x int64) int64 {
	x %= MOD
	if x < 0 {
		x += MOD
	}
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, q int
	if _, err := fmt.Fscan(in, &n, &m, &q); err != nil {
		return
	}
	edges := make([]Edge, m)
	maxW := make([]int64, n+1)
	for i := 0; i < m; i++ {
		var u, v int
		var w int64
		fmt.Fscan(in, &u, &v, &w)
		edges[i] = Edge{u, v, w}
		if w > maxW[u] {
			maxW[u] = w
		}
		if w > maxW[v] {
			maxW[v] = w
		}
	}

	if q <= m {
		dp := make([][]int64, q+1)
		for i := range dp {
			dp[i] = make([]int64, n+1)
			for j := 1; j <= n; j++ {
				dp[i][j] = INF_NEG
			}
		}
		dp[0][1] = 0
		ans := int64(0)
		for t := 1; t <= q; t++ {
			for j := 1; j <= n; j++ {
				dp[t][j] = INF_NEG
			}
			for _, e := range edges {
				if dp[t-1][e.u] != INF_NEG {
					if val := dp[t-1][e.u] + e.w; val > dp[t][e.v] {
						dp[t][e.v] = val
					}
				}
				if dp[t-1][e.v] != INF_NEG {
					if val := dp[t-1][e.v] + e.w; val > dp[t][e.u] {
						dp[t][e.u] = val
					}
				}
			}
			best := INF_NEG
			for i := 1; i <= n; i++ {
				if dp[t][i] > best {
					best = dp[t][i]
				}
			}
			ans = (ans + mod(best)) % MOD
		}
		fmt.Fprintln(out, ans)
		return
	}

	// q > m
	dp := make([][]int64, m+1)
	for i := range dp {
		dp[i] = make([]int64, n+1)
		for j := 1; j <= n; j++ {
			dp[i][j] = INF_NEG
		}
	}
	dp[0][1] = 0
	partial := int64(0)
	for t := 1; t <= m; t++ {
		for j := 1; j <= n; j++ {
			dp[t][j] = INF_NEG
		}
		for _, e := range edges {
			if dp[t-1][e.u] != INF_NEG {
				if val := dp[t-1][e.u] + e.w; val > dp[t][e.v] {
					dp[t][e.v] = val
				}
			}
			if dp[t-1][e.v] != INF_NEG {
				if val := dp[t-1][e.v] + e.w; val > dp[t][e.u] {
					dp[t][e.u] = val
				}
			}
		}
		best := INF_NEG
		for i := 1; i <= n; i++ {
			if dp[t][i] > best {
				best = dp[t][i]
			}
		}
		partial = (partial + mod(best)) % MOD
	}

	lineMap := make(map[int64]int64)
	for v := 1; v <= n; v++ {
		slope := maxW[v]
		if slope == 0 {
			continue
		}
		best := INF_NEG
		for t := 0; t <= m; t++ {
			if dp[t][v] == INF_NEG {
				continue
			}
			val := dp[t][v] - int64(t)*slope
			if val > best {
				best = val
			}
		}
		if cur, ok := lineMap[slope]; !ok || best > cur {
			lineMap[slope] = best
		}
	}

	lines := make([]Line, 0, len(lineMap))
	for k, b := range lineMap {
		lines = append(lines, Line{k: int64(k), b: b})
	}
	sort.Slice(lines, func(i, j int) bool {
		if lines[i].k == lines[j].k {
			return lines[i].b > lines[j].b
		}
		return lines[i].k < lines[j].k
	})
	hull := make([]Line, 0)
	for _, ln := range lines {
		if len(hull) > 0 && hull[len(hull)-1].k == ln.k {
			if hull[len(hull)-1].b >= ln.b {
				continue
			} else {
				hull[len(hull)-1] = ln
				continue
			}
		}
		for len(hull) >= 2 && isBad(hull[len(hull)-2], hull[len(hull)-1], ln) {
			hull = hull[:len(hull)-1]
		}
		hull = append(hull, ln)
	}

	pos := int64(m + 1)
	ans := partial
	for i := 0; i < len(hull) && pos <= int64(q); i++ {
		var right int64 = int64(q)
		if i+1 < len(hull) {
			num := hull[i].b - hull[i+1].b
			den := hull[i+1].k - hull[i].k
			cross := floorDiv(num, den)
			if cross < right {
				right = cross
			}
		}
		if right < pos {
			continue
		}
		if right > int64(q) {
			right = int64(q)
		}
		count := right - pos + 1
		sumX := (pos + right) * count / 2
		term := (mod(hull[i].k) * mod(sumX)) % MOD
		term = (term + mod(hull[i].b)*mod(count)) % MOD
		ans = (ans + term) % MOD
		pos = right + 1
	}

	fmt.Fprintln(out, ans)
}
