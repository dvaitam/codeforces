package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const MOD uint64 = 998244353
const INV2 uint64 = (MOD + 1) / 2

type Interval struct {
	L uint64
	R uint64
}

func decompose(l, r uint64, mp map[uint64][]uint64) {
	for l <= r {
		maxLen := l & -l
		remain := r - l + 1
		for maxLen > remain {
			maxLen >>= 1
		}
		mp[maxLen] = append(mp[maxLen], l)
		l += maxLen
	}
}

func sumInterval(l, r uint64) uint64 {
	cnt := (r - l + 1) % MOD
	firstLast := (l%MOD + r%MOD) % MOD
	return cnt * firstLast % MOD * INV2 % MOD
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var nA int
	if _, err := fmt.Fscan(in, &nA); err != nil {
		return
	}
	segA := make(map[uint64][]uint64)
	for i := 0; i < nA; i++ {
		var l, r uint64
		fmt.Fscan(in, &l, &r)
		decompose(l, r, segA)
	}
	var nB int
	fmt.Fscan(in, &nB)
	segB := make(map[uint64][]uint64)
	for i := 0; i < nB; i++ {
		var l, r uint64
		fmt.Fscan(in, &l, &r)
		decompose(l, r, segB)
	}

	var intervals []Interval
	for length, arrA := range segA {
		arrB, ok := segB[length]
		if !ok {
			continue
		}
		for _, a := range arrA {
			for _, b := range arrB {
				start := a ^ b
				intervals = append(intervals, Interval{start, start + length - 1})
			}
		}
	}

	if len(intervals) == 0 {
		fmt.Println(0)
		return
	}

	sort.Slice(intervals, func(i, j int) bool {
		if intervals[i].L == intervals[j].L {
			return intervals[i].R < intervals[j].R
		}
		return intervals[i].L < intervals[j].L
	})

	curL := intervals[0].L
	curR := intervals[0].R
	ans := uint64(0)
	for _, seg := range intervals[1:] {
		if seg.L <= curR+1 {
			if seg.R > curR {
				curR = seg.R
			}
		} else {
			ans = (ans + sumInterval(curL, curR)) % MOD
			curL = seg.L
			curR = seg.R
		}
	}
	ans = (ans + sumInterval(curL, curR)) % MOD

	fmt.Println(ans % MOD)
}
