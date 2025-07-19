package main

import (
	"bufio"
	"fmt"
	"os"
)

// gen constructs a vector of length n of positive integers such that
// the sum of their squares is a perfect square.
func gen(n int) []int {
	switch n {
	case 1:
		return []int{1}
	case 2:
		return []int{3, 4}
	}
	ret := make([]int, n)
	for i := 0; i < n; i++ {
		ret[i] = 1
	}
	diff := n - 1
	if diff%2 == 0 {
		ret[n-2] = 2
		diff += 3
	}
	// diff is now odd
	ret[n-1] = diff / 2
	return ret
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	row := gen(n)
	col := gen(m)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			fmt.Fprintf(writer, "%d", row[i]*col[j])
			if j+1 < m {
				writer.WriteByte(' ')
			}
		}
		writer.WriteByte('\n')
	}
}
