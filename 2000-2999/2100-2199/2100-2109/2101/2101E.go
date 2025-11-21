package main

import (
	"bufio"
	"fmt"
	"os"
)

type FastScanner struct {
	r *bufio.Reader
}

func NewFastScanner(r *bufio.Reader) *FastScanner {
	return &FastScanner{r: r}
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

func (fs *FastScanner) NextString() string {
	c, _ := fs.r.ReadByte()
	for c <= ' ' {
		c, _ = fs.r.ReadByte()
	}
	buf := make([]byte, 0, 16)
	for c > ' ' {
		buf = append(buf, c)
		c, _ = fs.r.ReadByte()
	}
	return string(buf)
}

func main() {
	in := NewFastScanner(bufio.NewReader(os.Stdin))
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	t := in.NextInt()
	for ; t > 0; t-- {
		n := in.NextInt()
		s := in.NextString()
		adj := make([][]int, n)
		for i := 0; i < n-1; i++ {
			u := in.NextInt() - 1
			v := in.NextInt() - 1
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}

		ones := make([]int, 0)
		for i := 0; i < n; i++ {
			if s[i] == '1' {
				ones = append(ones, i)
			}
		}
		k := len(ones)
		if k == 0 {
			for i := 0; i < n; i++ {
				fmt.Fprint(out, "-1 ")
			}
			fmt.Fprintln(out)
			continue
		}
		if k == 1 {
			for i := 0; i < n; i++ {
				if s[i] == '1' {
					fmt.Fprint(out, "1 ")
				} else {
					fmt.Fprint(out, "-1 ")
				}
			}
			fmt.Fprintln(out)
			continue
		}

		// BFS to find farthest
		maxDist := int64(0)
		var source int
		for idx := 0; idx < k; idx++ {
			u := ones[idx]
			dist := make([]int64, n)
			for i := 0; i < n; i++ {
				dist[i] = -1
			}
			q := make([]int, 0, n)
			dist[u] = 0
			q = append(q, u)
			for qi := 0; qi < len(q); qi++ {
				v := q[qi]
				for _, to := range adj[v] {
					if dist[to] == -1 {
						dist[to] = dist[v] + 1
						q = append(q, to)
					}
				}
			}
			for _, node := range ones {
				if dist[node] > maxDist {
					maxDist = dist[node]
					source = u
				}
			}
		}

		// BFS from source to get order by distance
		dist := make([]int64, n)
		for i := range dist {
			dist[i] = -1
		}
		order := make([]int, 0, n)
		q := make([]int, 0, n)
		dist[source] = 0
		q = append(q, source)
		order = append(order, source)
		for qi := 0; qi < len(q); qi++ {
			v := q[qi]
			for _, to := range adj[v] {
				if dist[to] == -1 {
					dist[to] = dist[v] + 1
					q = append(q, to)
					order = append(order, to)
				}
			}
		}

		res := make([]int64, n)
		for i := range res {
			res[i] = -1
		}
		res[source] = 1

		for _, v := range order[1:] {
			best := res[v] // either already set or -1
			for _, to := range adj[v] {
				if dist[to] == dist[v]+1 {
					if res[to] == -1 {
						continue
					}
					val := res[to] + 1
					if val > best {
						best = val
					}
				}
			}
			if best < 1 {
				best = 1
			}
			res[v] = best
		}

		for i := 0; i < n; i++ {
			fmt.Fprint(out, res[i], " ")
		}
		fmt.Fprintln(out)
	}
}
