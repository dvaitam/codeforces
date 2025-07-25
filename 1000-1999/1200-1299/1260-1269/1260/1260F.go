package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const MOD int64 = 1_000_000_007

type EventMap map[int]int64

var (
	n        int
	L, R     []int
	invLen   []int64
	g        [][]int
	maxColor int
	tot      int64
	totW     []int64
	prefix   []int64
	ans      int64
)

func mod(x int64) int64 {
	x %= MOD
	if x < 0 {
		x += MOD
	}
	return x
}

func modPow(a, e int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		e >>= 1
	}
	return res
}

func computeContribution(m EventMap) int64 {
	if len(m) == 0 {
		return 0
	}
	type pair struct {
		pos   int
		delta int64
	}
	arr := make([]pair, 0, len(m))
	for k, v := range m {
		arr = append(arr, pair{pos: k, delta: v})
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i].pos < arr[j].pos })

	curr := int64(0)
	prev := 1
	var w1, w2 int64
	for _, e := range arr {
		pos := e.pos
		if pos > maxColor+1 {
			pos = maxColor + 1
		}
		if pos > prev {
			sumTot := mod(prefix[pos-1] - prefix[prev-1])
			length := int64(pos - prev)
			w1 = mod(w1 + curr*sumTot%MOD)
			w2 = mod(w2 + mod(curr*curr%MOD*length%MOD))
		}
		curr = mod(curr + e.delta)
		prev = e.pos
		if prev > maxColor+1 {
			prev = maxColor + 1
		}
	}
	if prev <= maxColor {
		sumTot := mod(prefix[maxColor] - prefix[prev-1])
		length := int64(maxColor - prev + 1)
		w1 = mod(w1 + curr*sumTot%MOD)
		w2 = mod(w2 + mod(curr*curr%MOD*length%MOD))
	}
	res := mod(w1 - w2)
	return res
}

func dfs(v, p int) EventMap {
	em := make(EventMap)
	em[L[v]] = mod(em[L[v]] + invLen[v])
	em[R[v]+1] = mod(em[R[v]+1] - invLen[v])
	for _, to := range g[v] {
		if to == p {
			continue
		}
		child := dfs(to, v)
		contrib := computeContribution(child)
		ans = mod(ans + mod(tot*contrib%MOD))
		// merge child into em using small-to-large
		if len(em) < len(child) {
			em, child = child, em
		}
		for k, v := range child {
			em[k] = mod(em[k] + v)
		}
	}
	return em
}

func main() {
	in := bufio.NewReader(os.Stdin)
	fmt.Fscan(in, &n)
	L = make([]int, n+1)
	R = make([]int, n+1)
	invLen = make([]int64, n+1)
	maxColor = 0
	tot = 1
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &L[i], &R[i])
		if R[i] > maxColor {
			maxColor = R[i]
		}
		length := int64(R[i] - L[i] + 1)
		tot = tot * length % MOD
		invLen[i] = modPow(length, MOD-2)
	}
	g = make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}
	totWDiff := make([]int64, maxColor+3)
	for i := 1; i <= n; i++ {
		totWDiff[L[i]] = mod(totWDiff[L[i]] + invLen[i])
		totWDiff[R[i]+1] = mod(totWDiff[R[i]+1] - invLen[i])
	}
	totW = make([]int64, maxColor+2)
	prefix = make([]int64, maxColor+2)
	cur := int64(0)
	for i := 1; i <= maxColor; i++ {
		cur = mod(cur + totWDiff[i])
		totW[i] = cur
		prefix[i] = mod(prefix[i-1] + totW[i])
	}
	dfs(1, 0)
	fmt.Println(ans % MOD)
}
