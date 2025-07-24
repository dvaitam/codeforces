package main

import (
	"bufio"
	"fmt"
	"os"
)

func absInt64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	var k int64
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
	}

	// Precompute powers of k within reasonable bounds
	limit := int64(1e15)
	powers := make([]int64, 0)
	p := int64(1)
	for {
		powers = append(powers, p)
		if k == 1 {
			break
		}
		if k == -1 {
			if p == 1 {
				p = -1
				continue
			}
			break
		}
		if absInt64(p) > limit/absInt64(k) {
			break
		}
		p *= k
	}

	freq := make(map[int64]int64)
	freq[0] = 1
	var sum int64
	var result int64
	for _, v := range arr {
		sum += v
		for _, power := range powers {
			result += freq[sum-power]
		}
		freq[sum]++
	}

	fmt.Fprintln(writer, result)
}
