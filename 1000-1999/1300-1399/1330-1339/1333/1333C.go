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

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	prefixMap := make(map[int64]int)
	prefixSum := int64(0)
	prefixMap[0] = 0
	left := 0
	var res int64

	for i := 1; i <= n; i++ {
		prefixSum += a[i-1]
		if idx, ok := prefixMap[prefixSum]; ok && idx >= left {
			left = idx + 1
		}
		prefixMap[prefixSum] = i
		res += int64(i - left)
	}

	fmt.Fprintln(writer, res)
}
