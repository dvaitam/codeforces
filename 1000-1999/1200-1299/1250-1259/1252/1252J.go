package main

import (
	"bufio"
	"fmt"
	"os"
)

type pair struct {
	idx int
	val int64
}

func applySegment(dpPrev []int64, length int, slope2 int64, G2 int) []int64 {
	K := len(dpPrev) - 1
	inf := int64(-1 << 60)
	dpNext := make([]int64, K+1)
	for i := range dpNext {
		dpNext[i] = inf
	}
	deq := [2][]pair{}
	for k := 0; k <= K; k++ {
		for p := 0; p < 2; p++ {
			for len(deq[p]) > 0 && deq[p][0].idx < k-length {
				deq[p] = deq[p][1:]
			}
		}
		parity := k & 1
		if dpPrev[k] > inf/2 {
			val := dpPrev[k] - slope2*int64(k)
			dq := deq[parity]
			for len(dq) > 0 && dq[len(dq)-1].val <= val {
				dq = dq[:len(dq)-1]
			}
			dq = append(dq, pair{idx: k, val: val})
			deq[parity] = dq
		}
		for p := 0; p < 2; p++ {
			q := (k & 1) ^ p
			dq := deq[q]
			if len(dq) == 0 {
				continue
			}
			candidate := slope2*int64(k) + dq[0].val + int64(length)*int64(G2) - int64(G2)*int64((length&1)^p)
			if candidate > dpNext[k] {
				dpNext[k] = candidate
			}
		}
	}
	return dpNext
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var N, K int
	var G1, G2, G3 int
	if _, err := fmt.Fscan(in, &N, &K, &G1, &G2, &G3); err != nil {
		return
	}
	var s string
	if _, err := fmt.Fscan(in, &s); err != nil {
		return
	}

	segments := make([]int, 0)
	last := -1
	for i := 0; i <= N; i++ {
		if i == N || s[i] == '#' {
			segments = append(segments, i-last-1)
			last = i
		}
	}
	m := len(segments) - 1 // number of rocks

	inf := int64(-1 << 60)
	Kmax := K
	dp0 := make([]int64, Kmax+1)
	dp1 := make([]int64, Kmax+1)
	for i := 1; i <= Kmax; i++ {
		dp0[i] = inf
		dp1[i] = inf
	}
	dp1[0] = inf

	slope2 := int64(2*G1 - G2)
	for i := 0; i <= m; i++ {
		new0 := make([]int64, Kmax+1)
		new1 := make([]int64, Kmax+1)
		for x := 0; x <= Kmax; x++ {
			new0[x] = inf
			new1[x] = inf
		}
		for prev := 0; prev <= 1; prev++ {
			var dpPrev []int64
			if prev == 0 {
				dpPrev = dp0
			} else {
				dpPrev = dp1
			}
			length := segments[i] - prev
			if length < 0 {
				continue
			}
			temp := applySegment(dpPrev, length, slope2, G2)
			for k := 0; k <= Kmax; k++ {
				if temp[k] > new0[k] {
					new0[k] = temp[k]
				}
			}
			if i < m && length >= 1 && segments[i+1] > 0 {
				temp2 := applySegment(dpPrev, length-1, slope2, G2)
				for k := 0; k <= Kmax; k++ {
					if temp2[k] <= inf/2 {
						continue
					}
					v := temp2[k] + int64(2*G3)
					if v > new1[k] {
						new1[k] = v
					}
				}
			}
		}
		dp0, dp1 = new0, new1
	}
	ans := inf
	for k := 0; k <= Kmax; k++ {
		if dp0[k] > ans {
			ans = dp0[k]
		}
		if dp1[k] > ans {
			ans = dp1[k]
		}
	}
	if ans < 0 {
		ans = 0
	}
	fmt.Println(ans / 2)
}
