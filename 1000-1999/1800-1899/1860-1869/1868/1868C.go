package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353
const MAXH int = 60

var pow2 [MAXH + 2]int64
var pow2mod [MAXH + 2]int64
var fullNodes [MAXH + 1]int64
var pairsFull [MAXH + 1][]int64

func initPrecalc() {
	pow2[0] = 1
	pow2mod[0] = 1
	for i := 1; i <= MAXH+1; i++ {
		pow2[i] = pow2[i-1] << 1
		pow2mod[i] = pow2mod[i-1] * 2 % MOD
	}
	for h := 0; h <= MAXH; h++ {
		fullNodes[h] = pow2[h+1] - 1
	}
	pairsFull[0] = []int64{0, 1}
	for h := 1; h <= MAXH; h++ {
		pairs := make([]int64, 2*h+2)
		pairs[1] = 1
		// root-other pairs
		for d := 1; d <= h; d++ {
			pairs[d+1] = (pairs[d+1] + pow2mod[d]) % MOD
		}
		prev := pairsFull[h-1]
		for l, v := range prev {
			pairs[l] = (pairs[l] + 2*v) % MOD
		}
		for dl := 0; dl < h; dl++ {
			for dr := 0; dr < h; dr++ {
				l := dl + dr + 3
				pairs[l] = (pairs[l] + pow2mod[dl]*pow2mod[dr]) % MOD
			}
		}
		pairsFull[h] = pairs
	}
}

type Info struct {
	nodes []int64
	pairs []int64
}

func prefixCount(x, n int64) int64 {
	l, r := x, x
	cnt := int64(0)
	for l <= n {
		if r > n {
			cnt += n - l + 1
		} else {
			cnt += r - l + 1
		}
		l *= 2
		r = r*2 + 1
	}
	return cnt
}

func solve(x, n int64) Info {
	size := prefixCount(x, n)
	if size == 0 {
		return Info{nil, []int64{0}}
	}
	h := 0
	for h+1 <= MAXH && fullNodes[h+1] <= size {
		h++
	}
	if size == fullNodes[h] {
		nodes := make([]int64, h+1)
		for d := 0; d <= h; d++ {
			nodes[d] = pow2mod[d]
		}
		pairs := make([]int64, len(pairsFull[h]))
		copy(pairs, pairsFull[h])
		return Info{nodes, pairs}
	}
	var left, right Info
	if x*2 <= n {
		left = solve(x*2, n)
	}
	if x*2+1 <= n {
		right = solve(x*2+1, n)
	}
	maxDepth := 1
	if len(left.nodes)+1 > maxDepth {
		maxDepth = len(left.nodes) + 1
	}
	if len(right.nodes)+1 > maxDepth {
		maxDepth = len(right.nodes) + 1
	}
	nodes := make([]int64, maxDepth)
	nodes[0] = 1
	for i, v := range left.nodes {
		if i+1 >= len(nodes) {
			tmp := make([]int64, i+2)
			copy(tmp, nodes)
			nodes = tmp
		}
		nodes[i+1] = (nodes[i+1] + v) % MOD
	}
	for i, v := range right.nodes {
		if i+1 >= len(nodes) {
			tmp := make([]int64, i+2)
			copy(tmp, nodes)
			nodes = tmp
		}
		nodes[i+1] = (nodes[i+1] + v) % MOD
	}
	maxLen := 1
	if len(left.pairs) > maxLen {
		maxLen = len(left.pairs)
	}
	if len(right.pairs) > maxLen {
		maxLen = len(right.pairs)
	}
	if len(left.nodes)+len(right.nodes)+3 > maxLen {
		maxLen = len(left.nodes) + len(right.nodes) + 3
	}
	pairs := make([]int64, maxLen)
	pairs[1] = 1
	for i, v := range left.nodes {
		idx := i + 2
		if idx >= len(pairs) {
			tmp := make([]int64, idx+1)
			copy(tmp, pairs)
			pairs = tmp
		}
		pairs[idx] = (pairs[idx] + v) % MOD
	}
	for i, v := range right.nodes {
		idx := i + 2
		if idx >= len(pairs) {
			tmp := make([]int64, idx+1)
			copy(tmp, pairs)
			pairs = tmp
		}
		pairs[idx] = (pairs[idx] + v) % MOD
	}
	for i, v := range left.pairs {
		if i >= len(pairs) {
			tmp := make([]int64, i+1)
			copy(tmp, pairs)
			pairs = tmp
		}
		pairs[i] = (pairs[i] + v) % MOD
	}
	for i, v := range right.pairs {
		if i >= len(pairs) {
			tmp := make([]int64, i+1)
			copy(tmp, pairs)
			pairs = tmp
		}
		pairs[i] = (pairs[i] + v) % MOD
	}
	for i, a := range left.nodes {
		for j, b := range right.nodes {
			idx := i + j + 3
			if idx >= len(pairs) {
				tmp := make([]int64, idx+1)
				copy(tmp, pairs)
				pairs = tmp
			}
			pairs[idx] = (pairs[idx] + a*b%MOD) % MOD
		}
	}
	return Info{nodes, pairs}
}

func powMod(a, b int64) int64 {
	res := int64(1)
	a %= MOD
	for b > 0 {
		if b&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		b >>= 1
	}
	return res
}

func precomputeSums(m, maxK int) []int64 {
	sums := make([]int64, maxK+1)
	powArr := make([]int64, m)
	for i := range powArr {
		powArr[i] = 1
	}
	for k := 1; k <= maxK; k++ {
		var s int64
		for i := 0; i < m; i++ {
			powArr[i] = powArr[i] * int64(i) % MOD
			s += powArr[i]
		}
		sums[k] = s % MOD
	}
	return sums
}

func main() {
	initPrecalc()
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int64
		var m int
		fmt.Fscan(in, &n, &m)
		info := solve(1, n)
		maxLen := len(info.pairs) - 1
		sums := precomputeSums(m, maxLen)
		ans := int64(0)
		for k := 1; k <= maxLen; k++ {
			ck := info.pairs[k] % MOD
			if ck == 0 {
				continue
			}
			F := (powMod(int64(m), int64(k+1)) - sums[k]) % MOD
			if F < 0 {
				F += MOD
			}
			ans = (ans + ck*F%MOD*powMod(int64(m), n-int64(k))%MOD) % MOD
		}
		fmt.Fprintln(out, ans%MOD)
	}
}
