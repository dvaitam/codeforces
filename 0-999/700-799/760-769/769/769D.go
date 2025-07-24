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

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}

	const maxVal = 10000
	freq := make([]int64, maxVal+1)
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(in, &x)
		if x >= 0 && x <= maxVal {
			freq[x]++
		}
	}

	if k == 0 {
		var ans int64
		for _, cnt := range freq {
			if cnt > 1 {
				ans += cnt * (cnt - 1) / 2
			}
		}
		fmt.Fprintln(out, ans)
		return
	}

	masks := []int{}
	totalBits := 14
	for mask := 0; mask < (1 << totalBits); mask++ {
		if popcount(mask) == k {
			masks = append(masks, mask)
		}
	}

	var ans int64
	for x := 0; x <= maxVal; x++ {
		if freq[x] == 0 {
			continue
		}
		for _, mask := range masks {
			y := x ^ mask
			if y > maxVal {
				continue
			}
			if y > x {
				ans += freq[x] * freq[y]
			}
		}
	}

	fmt.Fprintln(out, ans)
}

func popcount(x int) int {
	x = x - ((x >> 1) & 0x5555)
	x = (x & 0x3333) + ((x >> 2) & 0x3333)
	x = (x + (x >> 4)) & 0x0F0F
	x = x + (x >> 8)
	return x & 0x3F
}
