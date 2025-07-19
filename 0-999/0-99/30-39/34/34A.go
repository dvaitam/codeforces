package main

import (
	"bufio"
	"fmt"
	"os"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	for {
		_, err := fmt.Fscan(reader, &n)
		if err != nil {
			break
		}
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		ansi, ansj := 0, 1
		minDiff := abs(a[1] - a[0])
		for i := 1; i < n; i++ {
			d := abs(a[i] - a[i-1])
			if d < minDiff {
				minDiff = d
				ansi = i - 1
				ansj = i
			}
		}
		// Output 1-based indices
		fmt.Println(ansi+1, ansj+1)
	}
}
