package main

import (
	"bufio"
	"fmt"
	"os"
)

func check(arr []int, half, k int) bool {
	cnt := make(map[int]int)
	for _, v := range arr {
		r := v % k
		if r < 0 {
			r += k
		}
		cnt[r]++
	}
	for _, c := range cnt {
		if c >= half {
			return true
		}
	}
	return false
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		freq := make(map[int]int)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			freq[a[i]]++
		}
		half := n / 2
		maxC := 0
		for _, c := range freq {
			if c > maxC {
				maxC = c
			}
		}
		if maxC >= half {
			fmt.Fprintln(out, -1)
			continue
		}
		ans := 1
		seen := make(map[int]bool)
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				diff := a[i] - a[j]
				if diff < 0 {
					diff = -diff
				}
				if diff == 0 {
					continue
				}
				for d := 1; d*d <= diff; d++ {
					if diff%d == 0 {
						if !seen[d] && d > ans {
							seen[d] = true
							if check(a, half, d) {
								ans = d
							}
						}
						d2 := diff / d
						if d2 != d && !seen[d2] && d2 > ans {
							seen[d2] = true
							if check(a, half, d2) {
								ans = d2
							}
						}
					}
				}
			}
		}
		fmt.Fprintln(out, ans)
	}
}
