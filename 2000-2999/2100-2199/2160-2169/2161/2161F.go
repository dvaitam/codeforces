package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	mod  int64 = 1000000007
	maxN       = 5000
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
	c, err := fs.r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		if err != nil {
			return 0
		}
		c, err = fs.r.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, err = fs.r.ReadByte()
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int(c-'0')
		c, err = fs.r.ReadByte()
		if err != nil {
			break
		}
	}
	return sign * val
}

func main() {
	scanner := NewFastScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	pow2 := make([]int64, maxN+1)
	pow2[0] = 1
	for i := 1; i <= maxN; i++ {
		pow2[i] = (pow2[i-1] * 2) % mod
	}

	t := scanner.NextInt()
	for ; t > 0; t-- {
		n := scanner.NextInt()
		adj := make([][]int, n)
		for i := 0; i < n-1; i++ {
			u := scanner.NextInt() - 1
			v := scanner.NextInt() - 1
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}

		if n <= 1 {
			fmt.Fprintln(out, 0)
			continue
		}

		parent := make([]int, n)
		for i := range parent {
			parent[i] = -1
		}
		order := make([]int, 0, n)
		stack := []int{0}
		for len(stack) > 0 {
			v := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			order = append(order, v)
			for _, to := range adj[v] {
				if to == parent[v] {
					continue
				}
				parent[to] = v
				stack = append(stack, to)
			}
		}

		sub := make([]int, n)
		for i := len(order) - 1; i >= 0; i-- {
			v := order[i]
			size := 1
			for _, to := range adj[v] {
				if to == parent[v] {
					continue
				}
				size += sub[to]
			}
			sub[v] = size
		}

		ans := int64(0)
		for v := 1; v < n; v++ {
			size := sub[v]
			other := n - size
			a := pow2[size] - 1
			if a < 0 {
				a += mod
			}
			b := pow2[other] - 1
			if b < 0 {
				b += mod
			}
			// Edge between v and parent[v] appears in MST of subsets touching both components.
			ans += (a * b) % mod
			if ans >= mod {
				ans -= mod
			}
		}

		fmt.Fprintln(out, ans)
	}
}
