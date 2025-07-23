package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	row1 := make([]int, n-1)
	row2 := make([]int, n-1)
	for i := 0; i < n-1; i++ {
		fmt.Fscan(reader, &row1[i])
	}
	for i := 0; i < n-1; i++ {
		fmt.Fscan(reader, &row2[i])
	}
	b := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &b[i])
	}
	prefix1 := make([]int, n)
	for i := 1; i < n; i++ {
		prefix1[i] = prefix1[i-1] + row1[i-1]
	}
	suffix2 := make([]int, n+1)
	for i := n - 1; i >= 1; i-- {
		suffix2[i] = suffix2[i+1] + row2[i-1]
	}
	const inf int = int(1e9)
	minTotal := inf
	for j1 := 1; j1 <= n; j1++ {
		timeThere := suffix2[j1] + b[j1-1] + prefix1[j1-1]
		for j2 := 1; j2 <= n; j2++ {
			if j1 == j2 {
				continue
			}
			timeBack := prefix1[j2-1] + b[j2-1] + suffix2[j2]
			total := timeThere + timeBack
			if total < minTotal {
				minTotal = total
			}
		}
	}
	fmt.Println(minTotal)
}
