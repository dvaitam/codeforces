package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	r := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &r[i])
	}

	sufMaxIdx := make([]int, n)
	maxIdx := n - 1
	sufMaxIdx[n-1] = maxIdx
	for i := n - 2; i >= 0; i-- {
		if r[i] >= r[maxIdx] {
			maxIdx = i
		}
		sufMaxIdx[i] = maxIdx
	}

	minIdx := 0
	for j := 1; j < n-1; j++ {
		if r[minIdx] < r[j] && r[j] < r[sufMaxIdx[j+1]] {
			fmt.Println(minIdx+1, j+1, sufMaxIdx[j+1]+1)
			return
		}
		if r[j] < r[minIdx] {
			minIdx = j
		}
	}

	fmt.Println("-1 -1 -1")
}
