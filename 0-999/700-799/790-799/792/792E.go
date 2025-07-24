package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	candMap := make(map[int64]struct{})
	for _, x := range a {
		for j := int64(1); j*j <= x; j++ {
			candMap[j] = struct{}{}
			if j-1 > 0 {
				candMap[j-1] = struct{}{}
			}
			q := x / j
			candMap[q] = struct{}{}
			if q-1 > 0 {
				candMap[q-1] = struct{}{}
			}
		}
	}
	cand := make([]int64, 0, len(candMap))
	for k := range candMap {
		if k > 0 {
			cand = append(cand, k)
		}
	}
	sort.Slice(cand, func(i, j int) bool { return cand[i] < cand[j] })
	best := int64(1<<63 - 1)
	for _, k := range cand {
		total := int64(0)
		feasible := true
		for _, x := range a {
			t := (x + k) / (k + 1)
			maxSet := x / k
			if t > maxSet {
				feasible = false
				break
			}
			total += t
			if total >= best {
				break
			}
		}
		if feasible && total < best {
			best = total
		}
	}
	fmt.Println(best)
}
