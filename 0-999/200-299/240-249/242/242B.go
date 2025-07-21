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

	L := make([]int64, n)
	R := make([]int64, n)
	const INF = int64(1e18)
	minL, maxR := INF, -INF

	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &L[i], &R[i])
		if L[i] < minL {
			minL = L[i]
		}
		if R[i] > maxR {
			maxR = R[i]
		}
	}

	result := -1
	for i := 0; i < n; i++ {
		if L[i] == minL && R[i] == maxR {
			result = i + 1
			break
		}
	}
	fmt.Fprintln(writer, result)
}
