package main

import (
	"bufio"
	"fmt"
	"os"
)

const N = 5005

var (
	g   [N][]int
	ord []int
	d   [N]int
	s   = 73
	c   int
)

func check(w *bufio.Writer, r *bufio.Reader, x int) bool {
	fmt.Fprintf(w, "? %d\n", x)
	w.Flush()
	var y int
	fmt.Fscan(r, &y)
	return y != 0
}

func dfs(v, p int, w *bufio.Writer, r *bufio.Reader) {
	ord = append(ord, v)
	for _, u := range g[v] {
		if u != p && d[u] >= s {
			c = u
		}
	}
	for _, u := range g[v] {
		if u != p && d[u] >= s {
			if u == c || check(w, r, u) {
				dfs(u, v, w, r)
				return
			}
		}
	}
}

func cfs(v, p int) {
	for _, u := range g[v] {
		if u != p {
			cfs(u, v)
			if d[u]+1 > d[v] {
				d[v] = d[u] + 1
			}
		}
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		for v := 1; v <= n; v++ {
			d[v] = 0
			g[v] = g[v][:0]
		}
		for i := 1; i < n; i++ {
			var x, y int
			fmt.Fscan(reader, &x, &y)
			g[x] = append(g[x], y)
			g[y] = append(g[y], x)
		}
		q := 0
		for i := 2; i <= n; i++ {
			if len(g[i]) == 1 {
				q = i
			}
		}
		fl := false
		for i := 0; i < s; i++ {
			if check(writer, reader, q) {
				fl = true
			}
		}
		if fl {
			fmt.Fprintf(writer, "! %d\n", q)
			continue
		}
		ord = ord[:0]
		cfs(1, -1)
		dfs(1, -1, writer, reader)
		l := 0
		r := len(ord) - 1
		for l != r {
			mid := (l + r + 1) >> 1
			if check(writer, reader, ord[mid]) {
				l = mid
			} else {
				r = mid - 2
				l--
				if l < 0 {
					l = 0
				}
				if r < 0 {
					r = 0
				}
			}
		}
		fmt.Fprintf(writer, "! %d\n", ord[l])
	}
}
