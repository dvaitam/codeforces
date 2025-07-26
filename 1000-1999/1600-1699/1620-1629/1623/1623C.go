package main

import (
	"bufio"
	"fmt"
	"os"
)

func canAchieve(h []int64, target int64) bool {
	n := len(h)
	b := make([]int64, n)
	copy(b, h)
	for i := n - 1; i >= 2; i-- {
		if b[i] < target {
			return false
		}
		d := (b[i] - target) / 3
		if d > h[i]/3 {
			d = h[i] / 3
		}
		b[i-1] += d
		b[i-2] += 2 * d
	}
	return b[0] >= target && b[1] >= target
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		h := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &h[i])
		}
		l, r := int64(0), int64(1e9)
		for l < r {
			mid := (l + r + 1) / 2
			if canAchieve(h, mid) {
				l = mid
			} else {
				r = mid - 1
			}
		}
		fmt.Fprintln(writer, l)
	}
}
