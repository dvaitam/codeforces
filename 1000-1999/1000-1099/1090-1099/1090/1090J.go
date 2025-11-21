package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var s, t string
	if _, err := fmt.Fscan(in, &s); err != nil {
		return
	}
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	if len(s) == 0 || len(t) == 0 {
		fmt.Println(0)
		return
	}

	tBytes := []byte(t)
	limit := buildLimits(tBytes)
	pi := prefixFunction(tBytes)
	up, maxUp := buildLifting(pi, limit)

	result := countDistinct([]byte(s), tBytes, pi, up, maxUp)

	fmt.Println(result)
}

func buildLimits(t []byte) []int {
	m := len(t)
	limit := make([]int, m+1)
	if m == 0 {
		return limit
	}
	z := make([]int, m)
	l, r := 0, 0
	for i := 1; i < m; i++ {
		if i <= r {
			if z[i-l] < r-i+1 {
				z[i] = z[i-l]
			} else {
				z[i] = r - i + 1
			}
		}
		for i+z[i] < m && t[z[i]] == t[i+z[i]] {
			z[i]++
		}
		if i+z[i]-1 > r {
			l = i
			r = i + z[i] - 1
		}
	}
	for d := 1; d < m; d++ {
		val := z[d]
		if m-d < val {
			val = m - d
		}
		limit[d] = val
	}
	return limit
}

func prefixFunction(t []byte) []int {
	m := len(t)
	pi := make([]int, m)
	for i := 1; i < m; i++ {
		j := pi[i-1]
		for j > 0 && t[i] != t[j] {
			j = pi[j-1]
		}
		if t[i] == t[j] {
			j++
		}
		pi[i] = j
	}
	return pi
}

func buildLifting(pi []int, limit []int) ([][]int, []int) {
	m := len(pi)
	log := 1
	for (1 << log) <= m {
		log++
	}
	up := make([][]int, log)
	for i := range up {
		up[i] = make([]int, m+1)
	}
	for L := 1; L <= m; L++ {
		up[0][L] = pi[L-1]
	}
	for k := 1; k < log; k++ {
		for v := 0; v <= m; v++ {
			up[k][v] = up[k-1][up[k-1][v]]
		}
	}
	maxUp := make([]int, m+1)
	for L := 1; L <= m; L++ {
		parent := up[0][L]
		maxUp[L] = max(limit[L], maxUp[parent])
	}
	return up, maxUp
}

func climb(up [][]int, node, threshold int) int {
	if node <= threshold {
		return node
	}
	for k := len(up) - 1; k >= 0; k-- {
		ancestor := up[k][node]
		if ancestor > threshold {
			node = ancestor
		}
	}
	return up[0][node]
}

func countDistinct(s, t []byte, pi []int, up [][]int, maxUp []int) int64 {
	m := len(t)
	cur := 0
	var ans int64
	for idx := 0; idx < len(s); idx++ {
		ch := s[idx]
		for cur > 0 && ch != t[cur] {
			cur = pi[cur-1]
		}
		if ch == t[cur] {
			cur++
		}
		ancestor := climb(up, cur, idx)
		best := maxUp[ancestor]
		ans += int64(m - best)
		if cur == m {
			cur = pi[cur-1]
		}
	}
	return ans
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
