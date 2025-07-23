package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	t := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &t[i])
	}

	isGood := func(i int) bool {
		if i%2 == 0 {
			return t[i] < t[i+1]
		}
		return t[i] > t[i+1]
	}

	bad := []int{}
	for i := 0; i < n-1; i++ {
		if !isGood(i) {
			bad = append(bad, i)
		}
	}
	if len(bad) > 4 {
		fmt.Println(0)
		return
	}

	candidateSet := make(map[int]struct{})
	for _, b := range bad {
		candidateSet[b] = struct{}{}
		candidateSet[b+1] = struct{}{}
	}
	candidates := []int{}
	for idx := range candidateSet {
		if idx >= 0 && idx < n {
			candidates = append(candidates, idx)
		}
	}

	seen := make(map[[2]int]struct{})
	ans := 0

	check := func(i, j int) bool {
		if i == j {
			return false
		}
		t[i], t[j] = t[j], t[i]
		defer func() { t[i], t[j] = t[j], t[i] }()

		indices := []int{}
		indices = append(indices, bad...)
		indices = append(indices, i-1, i, j-1, j)
		for _, idx := range indices {
			if idx >= 0 && idx < n-1 {
				if idx%2 == 0 {
					if t[idx] >= t[idx+1] {
						return false
					}
				} else {
					if t[idx] <= t[idx+1] {
						return false
					}
				}
			}
		}
		return true
	}

	for _, i := range candidates {
		for j := 0; j < n; j++ {
			if i == j {
				continue
			}
			a, b := i, j
			if a > b {
				a, b = b, a
			}
			key := [2]int{a, b}
			if _, ok := seen[key]; ok {
				continue
			}
			if check(i, j) {
				ans++
			}
			seen[key] = struct{}{}
		}
	}
	fmt.Println(ans)
}
