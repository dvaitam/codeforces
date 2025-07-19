package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	s := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &s[i])
	}
	m := len(s[0])

	size := 1 << m
	c := make([]uint64, size)
	// build initial masks for pairs
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			mask := 0
			for w := 0; w < m; w++ {
				if s[i][w] == s[j][w] {
					mask |= 1 << w
				}
			}
			bit := (uint64(1) << uint(i)) | (uint64(1) << uint(j))
			c[mask] |= bit
		}
	}
	// propagate supersets: for each mask, include from masks with one extra bit
	for mask := size - 1; mask >= 0; mask-- {
		for b := 0; b < m; b++ {
			if (mask & (1 << b)) == 0 {
				c[mask] |= c[mask|(1<<b)]
			}
		}
		if mask == 0 {
			break
		}
	}
	// binomial coefficients
	d := make([][]float64, m+1)
	for i := 0; i <= m; i++ {
		d[i] = make([]float64, i+1)
		d[i][0] = 1
		d[i][i] = 1
		for j := 1; j < i; j++ {
			d[i][j] = d[i-1][j] + d[i-1][j-1]
		}
	}
	// compute expected value
	var ans float64
	for mask := 0; mask < size; mask++ {
		t := bits.OnesCount(uint(mask))
		cnt := bits.OnesCount64(c[mask])
		ans += float64(cnt) / float64(n) * (1.0 / d[m][t])
	}
	fmt.Printf("%.15f\n", ans)
}
