package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

const LOG = 20

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

func ceilLog2(x int) int {
	if x <= 1 {
		return 0
	}
	return bits.Len(uint(x - 1))
}

func lca(u, v int, depth []int, up [][LOG]int) int {
	if depth[u] < depth[v] {
		u, v = v, u
	}
	diff := depth[u] - depth[v]
	for j := 0; diff > 0; j++ {
		if diff&1 == 1 {
			u = up[u][j]
		}
		diff >>= 1
	}
	if u == v {
		return u
	}
	for j := LOG - 1; j >= 0; j-- {
		if up[u][j] != up[v][j] {
			u = up[u][j]
			v = up[v][j]
		}
	}
	return up[u][0]
}

func dist(u, v int, depth []int, up [][LOG]int) int {
	l := lca(u, v, depth, up)
	return depth[u] + depth[v] - 2*depth[l]
}

func main() {
	in := NewFastScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	t := in.NextInt()
	for ; t > 0; t-- {
		n := in.NextInt()
		parents := make([]int, n+1)
		for i := 2; i <= n; i++ {
			parents[i] = in.NextInt()
		}

		depth := make([]int, n+1)
		up := make([][LOG]int, n+1)
		deg := make([]int, n+1)

		res := make([]int, n)
		res[0] = 1
		depth[1] = 0
		for j := 0; j < LOG; j++ {
			up[1][j] = 0
		}

		a, b := 1, 1
		diam := 0
		invalid := false

		for i := 2; i <= n; i++ {
			p := parents[i]
			deg[p]++
			deg[i]++
			if deg[p] > 3 || deg[i] > 3 {
				invalid = true
			}

			depth[i] = depth[p] + 1
			up[i][0] = p
			for j := 1; j < LOG; j++ {
				up[i][j] = up[up[i][j-1]][j-1]
			}

			da := dist(i, a, depth, up)
			if da > diam {
				diam = da
				b = i
			}
			db := dist(i, b, depth, up)
			if db > diam {
				diam = db
				a = i
			}

			if invalid {
				res[i-1] = -1
				continue
			}

			radius := (diam + 1) / 2
			dHeight := radius + 1
			dNodes := ceilLog2(i + 1)
			ans := dHeight
			if dNodes > ans {
				ans = dNodes
			}
			res[i-1] = ans
		}

		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, res[i])
		}
		fmt.Fprintln(out)
	}
}
