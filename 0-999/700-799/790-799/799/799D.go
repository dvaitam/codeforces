package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func solve(targetA, targetB, h, w int64, factors []int64) int {
	if h >= targetA && w >= targetB {
		return 0
	}
	sort.Slice(factors, func(i, j int) bool { return factors[i] > factors[j] })
	limit := 40
	if len(factors) < limit {
		limit = len(factors)
	}
	factors = factors[:limit]

	type pair struct{ w, h int64 }
	state := map[int64]int64{h: w}
	for i, f := range factors {
		next := make(map[int64]int64, len(state)*2)
		for width, height := range state {
			// skip current extension
			if height > next[width] {
				next[width] = height
			}
			// apply to width
			nw := width * f
			if nw > targetA {
				nw = targetA
			}
			if height > next[nw] {
				next[nw] = height
			}
			// apply to height
			nh := height * f
			if nh > targetB {
				nh = targetB
			}
			if nh > next[width] {
				next[width] = nh
			}
		}
		// prune dominated states
		arr := make([]pair, 0, len(next))
		for wv, hv := range next {
			arr = append(arr, pair{wv, hv})
		}
		sort.Slice(arr, func(i, j int) bool { return arr[i].w < arr[j].w })
		state = make(map[int64]int64)
		maxH := int64(0)
		for _, p := range arr {
			if p.h > maxH {
				maxH = p.h
				state[p.w] = p.h
			}
		}
		for wv, hv := range state {
			if wv >= targetA && hv >= targetB {
				return i + 1
			}
		}
	}
	return -1
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var a, b, h, w int64
	var n int
	if _, err := fmt.Fscan(reader, &a, &b, &h, &w, &n); err != nil {
		return
	}
	factors := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &factors[i])
	}
	res1 := solve(a, b, h, w, factors)
	res2 := solve(b, a, h, w, factors)
	res := -1
	if res1 != -1 && res2 != -1 {
		if res1 < res2 {
			res = res1
		} else {
			res = res2
		}
	} else if res1 != -1 {
		res = res1
	} else if res2 != -1 {
		res = res2
	}
	fmt.Println(res)
}
