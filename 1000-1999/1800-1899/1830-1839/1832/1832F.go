package main

import (
	"bufio"
	"fmt"
	"os"
)

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k int
	var x, m int64
	if _, err := fmt.Fscan(reader, &n, &k, &x, &m); err != nil {
		return
	}
	l := make([]int64, n)
	r := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &l[i], &r[i])
	}

	// This implementation only computes zombies entering without using generators.
	// It does not attempt to optimize generator placement.
	var result int64
	for i := 0; i < n; i++ {
		manual := r[i] - l[i]
		zombies := x - manual
		if zombies < 0 {
			zombies = 0
		}
		result += zombies
	}
	fmt.Fprintln(writer, result)
}
