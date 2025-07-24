package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxBits = 21
const maxMask = 1 << maxBits

func insert(first, second []int, mask, idx int) {
	if idx > first[mask] {
		second[mask] = first[mask]
		first[mask] = idx
	} else if idx > second[mask] && idx != first[mask] {
		second[mask] = idx
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	first := make([]int, maxMask)
	second := make([]int, maxMask)
	for i := 0; i < maxMask; i++ {
		first[i] = -1
		second[i] = -1
	}

	for i, v := range a {
		insert(first, second, v, i)
	}

	for bit := 0; bit < maxBits; bit++ {
		for mask := 0; mask < maxMask; mask++ {
			if mask&(1<<bit) == 0 {
				other := mask | (1 << bit)
				if first[other] != -1 {
					insert(first, second, mask, first[other])
				}
				if second[other] != -1 {
					insert(first, second, mask, second[other])
				}
			}
		}
	}

	ans := 0
	for i := 0; i < n; i++ {
		best := 0
		for b := maxBits - 1; b >= 0; b-- {
			if a[i]&(1<<b) != 0 {
				continue
			}
			cand := best | (1 << b)
			if second[cand] > i {
				best = cand
			}
		}
		val := a[i] | best
		if val > ans {
			ans = val
		}
	}

	fmt.Fprintln(out, ans)
}
