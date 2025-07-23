package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

type Mat struct{ a00, a01, a10, a11 int64 }

func mul(a, b Mat, mod int64) Mat {
	return Mat{
		(a.a00*b.a00 + a.a01*b.a10) % mod,
		(a.a00*b.a01 + a.a01*b.a11) % mod,
		(a.a10*b.a00 + a.a11*b.a10) % mod,
		(a.a10*b.a01 + a.a11*b.a11) % mod,
	}
}

func mulVec(m Mat, v0, v1, mod int64) (int64, int64) {
	x0 := (m.a00*v0 + m.a01*v1) % mod
	x1 := (m.a10*v0 + m.a11*v1) % mod
	return x0, x1
}

func precompute(base []int64, mod int64, K int64) [][]Mat {
	N := len(base)
	maxP := bits.Len64(uint64(K)) + 1
	step := make([][]Mat, maxP)
	step[0] = make([]Mat, N)
	for r := 0; r < N; r++ {
		prev := r - 1
		if prev < 0 {
			prev += N
		}
		step[0][r] = Mat{base[r] % mod, base[prev] % mod, 1 % mod, 0}
	}
	for p := 0; p < maxP-1; p++ {
		step[p+1] = make([]Mat, N)
		shift := 1 << uint(p)
		for r := 0; r < N; r++ {
			r2 := (r + shift) % N
			step[p+1][r] = mul(step[p][r2], step[p][r], mod)
		}
	}
	return step
}

func apply(step [][]Mat, startR int, L int64, mod int64) (Mat, int) {
	res := Mat{1, 0, 0, 1}
	r := startR
	p := 0
	for L > 0 {
		if L&1 == 1 {
			res = mul(step[p][r], res, mod)
			r = (r + (1 << uint(p))) % len(step[0])
		}
		L >>= 1
		p++
	}
	return res, r
}

func solve(K, P int64, base []int64, mods map[int64]int64) int64 {
	if K == 0 {
		return 0 % P
	}
	if K == 1 {
		return 1 % P
	}
	step := precompute(base, P, K)
	N := len(base)
	// important indices: mods and mods+1
	importantMap := map[int64]struct{}{}
	for j := range mods {
		importantMap[j] = struct{}{}
		importantMap[j+1] = struct{}{}
	}
	importantMap[K] = struct{}{}
	// gather and sort
	events := make([]int64, 0, len(importantMap))
	for x := range importantMap {
		if x >= 1 {
			events = append(events, x)
		}
	}
	sort.Slice(events, func(i, j int) bool { return events[i] < events[j] })
	// remove duplicates
	uniq := events[:0]
	var last int64 = -1
	for _, x := range events {
		if x != last {
			uniq = append(uniq, x)
			last = x
		}
	}
	events = uniq
	getS := func(idx int64) int64 {
		if v, ok := mods[idx]; ok {
			return v % P
		}
		return base[int(idx%int64(N))] % P
	}
	v0, v1 := int64(1%P), int64(0)
	cur := int64(1)
	r := int(cur % int64(N))
	ei := 0
	for cur <= K-1 {
		x := events[ei]
		if x > K-1 {
			x = K
		}
		if cur <= x-1 {
			length := (x - 1) - cur + 1
			m, r2 := apply(step, r, length, P)
			v0, v1 = mulVec(m, v0, v1, P)
			r = r2
			cur = x
		}
		if cur > K-1 {
			break
		}
		// handle single index cur (== events[ei])
		sx := getS(cur)
		sx1 := getS(cur - 1)
		m := Mat{sx % P, sx1 % P, 1 % P, 0}
		v0, v1 = mulVec(m, v0, v1, P)
		r = (r + 1) % N
		cur++
		ei++
	}
	return v0 % P
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var K, P int64
	if _, err := fmt.Fscan(in, &K, &P); err != nil {
		return
	}
	var N int
	fmt.Fscan(in, &N)
	base := make([]int64, N)
	for i := 0; i < N; i++ {
		fmt.Fscan(in, &base[i])
		base[i] %= P
	}
	var M int
	fmt.Fscan(in, &M)
	mods := make(map[int64]int64, M)
	for i := 0; i < M; i++ {
		var j, v int64
		fmt.Fscan(in, &j, &v)
		mods[j] = v % P
	}
	ans := solve(K, P, base, mods)
	out := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(out, ans)
	out.Flush()
}
