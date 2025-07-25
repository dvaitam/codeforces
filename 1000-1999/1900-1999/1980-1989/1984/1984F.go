package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const MOD int64 = 998244353

type anchor struct {
	idx int
	typ int // 0 = const, 1 = T - val
	val int64
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

// check verifies if there exists an array satisfying the constraints
func check(n int, m int64, s []byte, b []int64) bool {
	anchors := make([]anchor, 0, 2*n+2)
	anchors = append(anchors, anchor{idx: 0, typ: 0, val: 0})
	for i := 1; i <= n; i++ {
		ch := s[i-1]
		if ch == 'P' {
			anchors = append(anchors, anchor{idx: i, typ: 0, val: b[i-1]})
		} else if ch == 'S' {
			anchors = append(anchors, anchor{idx: i - 1, typ: 1, val: b[i-1]})
		}
	}

	sort.Slice(anchors, func(i, j int) bool {
		if anchors[i].idx == anchors[j].idx {
			return anchors[i].typ < anchors[j].typ
		}
		return anchors[i].idx < anchors[j].idx
	})

	const INF int64 = 1 << 60
	L, R := int64(-INF), int64(INF)

	for i := 0; i < len(anchors)-1; i++ {
		a := anchors[i]
		b2 := anchors[i+1]
		dist := b2.idx - a.idx
		if dist < 0 {
			dist = -dist
		}
		limit := int64(dist) * m
		switch {
		case a.typ == 0 && b2.typ == 0:
			diff := abs(b2.val - a.val)
			if diff > limit {
				return false
			}
		case a.typ == 1 && b2.typ == 1:
			diff := abs(b2.val - a.val)
			if diff > limit {
				return false
			}
		case a.typ == 0 && b2.typ == 1:
			lower := a.val + b2.val - limit
			upper := a.val + b2.val + limit
			if lower > L {
				L = lower
			}
			if upper < R {
				R = upper
			}
			if L > R {
				return false
			}
		case a.typ == 1 && b2.typ == 0:
			lower := b2.val + a.val - limit
			upper := b2.val + a.val + limit
			if lower > L {
				L = lower
			}
			if upper < R {
				R = upper
			}
			if L > R {
				return false
			}
		}
	}
	return L <= R
}

func solveCase(n int, m int64, str string, b []int64) int64 {
	s := []byte(str)
	var qs []int
	for i, ch := range s {
		if ch == '?' {
			qs = append(qs, i)
		}
	}
	var ans int64
	var dfs func(int)
	dfs = func(idx int) {
		if idx == len(qs) {
			if check(n, m, s, b) {
				ans++
				if ans >= MOD {
					ans %= MOD
				}
			}
			return
		}
		pos := qs[idx]
		s[pos] = 'P'
		dfs(idx + 1)
		s[pos] = 'S'
		dfs(idx + 1)
		s[pos] = '?'
	}
	dfs(0)
	return ans % MOD
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		var m int64
		fmt.Fscan(in, &n, &m)
		var str string
		fmt.Fscan(in, &str)
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}
		res := solveCase(n, m, str, b)
		fmt.Fprintln(out, res%MOD)
	}
}
