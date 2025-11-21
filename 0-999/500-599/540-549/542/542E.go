package main

import (
	"bufio"
	"fmt"
	"os"
)

type FastScanner struct {
	r *bufio.Reader
}

func NewFastScanner() *FastScanner {
	return &FastScanner{r: bufio.NewReader(os.Stdin)}
}

func (fs *FastScanner) NextInt() int {
	sign := 1
	val := 0
	c, _ := fs.r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		c, _ = fs.r.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, _ = fs.r.ReadByte()
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int(c-'0')
		c, _ = fs.r.ReadByte()
	}
	return sign * val
}

func main() {
	in := NewFastScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	n := in.NextInt()
	m := in.NextInt()
	adj := make([][]int, n)
	for i := 0; i < m; i++ {
		u := in.NextInt() - 1
		v := in.NextInt() - 1
		adj[u] = append(adj[u], v)
	}

	mod := 1000000007
	sizeDP := 1 << (n - 1)
	dp := make([]int, sizeDP)
	dp[0] = 1
	for mask := 1; mask < sizeDP; mask++ {
		for i := 1; i < n; i++ {
			if mask&(1<<(i-1)) == 0 {
				continue
			}
			prevMask := mask ^ (1 << (i - 1))
			for _, u := range adj[i] {
				if i == 0 {
					continue
				}
				if u == 0 {
					dp[mask] = (dp[mask] + dp[prevMask]) % mod
				}
			}
		}
	}
	fmt.Fprintln(out, dp[sizeDP-1])
}
