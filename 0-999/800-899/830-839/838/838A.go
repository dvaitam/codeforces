package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	grid := make([][]byte, n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		grid[i] = []byte(s)
	}

	prefix := make([][]int, n+1)
	for i := range prefix {
		prefix[i] = make([]int, m+1)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			val := 0
			if grid[i][j] == '1' {
				val = 1
			}
			prefix[i+1][j+1] = prefix[i+1][j] + prefix[i][j+1] - prefix[i][j] + val
		}
	}

	maxNM := n
	if m > maxNM {
		maxNM = m
	}
	ans := n * m
	for k := 2; k <= maxNM; k++ {
		total := 0
		for r := 0; r < n; r += k {
			for c := 0; c < m; c += k {
				r2 := r + k
				if r2 > n {
					r2 = n
				}
				c2 := c + k
				if c2 > m {
					c2 = m
				}
				ones := prefix[r2][c2] - prefix[r][c2] - prefix[r2][c] + prefix[r][c]
				size := k * k
				if size-ones < ones {
					total += size - ones
				} else {
					total += ones
				}
			}
		}
		if total < ans {
			ans = total
		}
	}

	fmt.Println(ans)
}
