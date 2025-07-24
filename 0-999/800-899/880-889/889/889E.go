package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	candidates := make(map[int64]struct{})
	if a[0] > 0 {
		candidates[a[0]-1] = struct{}{}
	}
	limit := int64(100000)
	if a[0]-1 <= limit {
		for x := int64(0); x < a[0]; x++ {
			candidates[x] = struct{}{}
		}
	} else {
		for _, v := range a[1:] {
			if v == 0 {
				continue
			}
			for k := int64(1); ; k++ {
				val := k*v - 1
				if val >= a[0] || val < 0 {
					break
				}
				candidates[val] = struct{}{}
				if k*v-1 > limit {
					break
				}
			}
		}
		for i := int64(0); i < limit && a[0]-1-i >= 0; i++ {
			candidates[a[0]-1-i] = struct{}{}
		}
	}

	var best int64
	for x := range candidates {
		if x < 0 || x >= a[0] {
			continue
		}
		cur := x
		var sum int64
		for i := 0; i < n; i++ {
			cur = cur % a[i]
			sum += cur
		}
		if sum > best {
			best = sum
		}
	}
	fmt.Fprintln(out, best)
}
