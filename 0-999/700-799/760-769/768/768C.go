package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k, x int
	if _, err := fmt.Fscan(reader, &n, &k, &x); err != nil {
		return
	}
	const maxVal = 1024
	cnt := make([]int, maxVal)
	for i := 0; i < n; i++ {
		var v int
		fmt.Fscan(reader, &v)
		cnt[v]++
	}

	for iter := 0; iter < k; iter++ {
		next := make([]int, maxVal)
		odd := 0
		for v := 0; v < maxVal; v++ {
			c := cnt[v]
			if c == 0 {
				continue
			}
			var toXor int
			if odd == 0 {
				toXor = (c + 1) / 2
			} else {
				toXor = c / 2
			}
			next[v^x] += toXor
			next[v] += c - toXor
			if c%2 == 1 {
				odd ^= 1
			}
		}
		cnt = next
	}

	minVal, maxValIdx := -1, -1
	for i := 0; i < maxVal; i++ {
		if cnt[i] > 0 {
			if minVal == -1 {
				minVal = i
			}
			maxValIdx = i
		}
	}
	fmt.Fprintln(writer, maxValIdx, minVal)
}
