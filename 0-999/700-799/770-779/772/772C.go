package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// extended gcd: returns g = gcd(a,b) and x,y such that ax + by = g
func gcdex(a, b int) (g, x, y int) {
	if a == 0 {
		return b, 0, 1
	}
	g1, x1, y1 := gcdex(b%a, a)
	x = y1 - (b/a)*x1
	y = x1
	return g1, x, y
}

// modular inverse of a mod m, m > 0, a and m should be coprime
func inv(a, m int) int {
	_, x, _ := gcdex(a, m)
	x %= m
	if x < 0 {
		x += m
	}
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	bad := make([]bool, m)
	for i := 0; i < n; i++ {
		var k int
		fmt.Fscan(in, &k)
		if k >= 0 && k < m {
			bad[k] = true
		}
	}

	mp := make(map[int][]int)
	for j := 0; j < m; j++ {
		if !bad[j] {
			d := gcd(m, j)
			mp[d] = append(mp[d], j)
		}
	}
	// collect and sort gcd groups
	gcds := make([]int, 0, len(mp))
	for d := range mp {
		gcds = append(gcds, d)
	}
	sort.Ints(gcds)
	cnt := len(gcds)

	ans := make([]int, cnt)
	best := make([]int, cnt)
	for i := range best {
		best[i] = -1
	}
	// DP from largest to smallest
	for i := cnt - 1; i >= 0; i-- {
		for j := i + 1; j < cnt; j++ {
			if gcds[j]%gcds[i] == 0 && ans[j] > ans[i] {
				ans[i] = ans[j]
				best[i] = j
			}
		}
		ans[i] += len(mp[gcds[i]])
	}
	// find starting group
	bestI := 0
	for i := 1; i < cnt; i++ {
		if ans[i] > ans[bestI] {
			bestI = i
		}
	}
	fmt.Fprintln(out, ans[bestI])

	prevK, prevG := 1, 1
	// reconstruct sequence
	for idx := bestI; idx != -1; idx = best[idx] {
		g := gcds[idx]
		for _, k := range mp[g] {
			// multiplier to go from prevK to k mod m
			modSeg := m / prevG
			a := k / prevG
			b := prevK / prevG
			x := (a * inv(b, modSeg)) % modSeg
			fmt.Fprintf(out, "%d ", x)
			prevK = k
			prevG = g
		}
	}
	fmt.Fprintln(out)
}
