package main

import (
	"bufio"
	"fmt"
	"os"
)

type segTree struct {
	n   int
	val []int64
}

func newSegTree(a []int64) *segTree {
	n := len(a)
	size := 1
	for size < n {
		size <<= 1
	}
	val := make([]int64, 2*size)
	const inf int64 = 1 << 60
	for i := range val {
		val[i] = inf
	}
	for i := 0; i < n; i++ {
		val[size+i] = a[i]
	}
	for i := size - 1; i > 0; i-- {
		if val[i<<1] < val[i<<1|1] {
			val[i] = val[i<<1]
		} else {
			val[i] = val[i<<1|1]
		}
	}
	return &segTree{n: n, val: val}
}

func (st *segTree) rangeMin(l, r int) int64 {
	if l > r {
		return 1 << 60
	}
	l += len(st.val) / 2
	r += len(st.val) / 2
	res := int64(1 << 60)
	for l <= r {
		if l&1 == 1 {
			if st.val[l] < res {
				res = st.val[l]
			}
			l++
		}
		if r&1 == 0 {
			if st.val[r] < res {
				res = st.val[r]
			}
			r--
		}
		l >>= 1
		r >>= 1
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var N int
	var D, Ts, Tf, Tw int64
	if _, err := fmt.Fscan(in, &N, &D, &Ts, &Tf, &Tw); err != nil {
		return
	}

	P := make([]int64, N)
	for i := 0; i < N; i++ {
		fmt.Fscan(in, &P[i])
	}

	// Helper prefix sums of P for quick range queries.
	prefP := make([]int64, N+1)
	for i := 0; i < N; i++ {
		prefP[i+1] = prefP[i] + P[i]
	}

	// If flying is faster, we fly whenever possible after crossing the initial visited prefix.
	if Ts > Tf {
		best := int64(1 << 62)
		c := 2*Tw + Ts - Tf
		// feasibility check per k.
		maxK := 0
		if c < 0 {
			// pick the furthest feasible k.
			for k := 0; k < N; k++ {
				need := int64(k) * D
				sum := prefP[k+1]
				ok := false
				if k == N-1 {
					ok = sum >= need
				} else {
					ok = sum > need
				}
				if ok {
					maxK = k
				}
			}
		}
		candidates := []int{0}
		if c < 0 {
			candidates = append(candidates, maxK)
		}
		for _, k := range candidates {
			need := int64(k) * D
			sum := prefP[k+1]
			if k == N-1 {
				if sum < need {
					continue
				}
			} else {
				if sum <= need {
					continue
				}
			}
			time := int64(2*k)*Tw + int64(k)*Ts + int64(N-1-k)*Tf
			if time < best {
				best = time
			}
		}
		fmt.Fprintln(out, best)
		return
	}

	// Build prefix of (P[i+1] - D) to help find when stamina drops below D.
	// prefix[i] corresponds to island i (0-based).
	prefix := make([]int64, N)
	for i := 0; i < N-1; i++ {
		prefix[i+1] = prefix[i] + (P[i+1] - D)
	}

	// Segment tree over prefix[0..N-2] (positions where a move starts).
	segArr := prefix[:N-1]
	st := newSegTree(segArr)

	// Finds the smallest x >= pos and x <= N-2 such that prefix[x] < threshold.
	// Returns -1 if stamina never drops below D starting from pos.
	firstFail := func(pos int, stamina int64) int {
		if pos >= N-1 {
			return -1
		}
		threshold := prefix[pos] + D - stamina
		if st.rangeMin(pos, N-2) >= threshold {
			return -1
		}
		l, r := pos, N-2
		for l < r {
			m := (l + r) >> 1
			if st.rangeMin(pos, m) < threshold {
				r = m
			} else {
				l = m + 1
			}
		}
		return l
	}

	// dp[i] = minimal time to reach island N starting from island i with stamina P[i].
	dp := make([]int64, N)
	dp[N-1] = 0
	for i := N - 2; i >= 0; i-- {
		fail := firstFail(i, P[i])
		if fail == -1 {
			dp[i] = int64(N-1-i) * Ts
		} else {
			swimLen := int64(fail - i)
			dp[i] = swimLen*Ts + Tf + dp[fail+1]
		}
	}

	best := int64(1 << 62)

	for k := 0; k < N; k++ {
		// k is number of edges in visited prefix (walked to island k and back), position after swims is island k.
		sum := prefP[k+1]
		need := int64(k) * D
		if k == N-1 {
			if sum < need {
				continue
			}
		} else {
			if sum <= need {
				continue
			}
		}
		stamina := sum - need
		timePrefix := int64(2*k)*Tw + int64(k)*Ts
		if k == N-1 {
			if timePrefix < best {
				best = timePrefix
			}
			continue
		}
		fail := firstFail(k, stamina)
		var rem int64
		if fail == -1 {
			rem = int64(N-1-k) * Ts
		} else {
			rem = int64(fail-k)*Ts + Tf + dp[fail+1]
		}
		total := timePrefix + rem
		if total < best {
			best = total
		}
	}

	fmt.Fprintln(out, best)
}
